package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/quicky-dev/auth-service/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var users = util.MongoClient.Database("Auth").Collection("Users")

// TODO file an issue on the mongo-driver org, because this apparently creates
// an index for usernames :thonk:
func setupIndexes() error {
	_, err := users.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bson.M{"username": ""},
			Options: options.Index().SetUnique(true)})

	if err != nil {
		return err
	}

	return nil
}

func initializeModel() bool {
	err := setupIndexes()
	log.Println(err)
	return true
}

var Initialized = initializeModel()

type User struct {
	ID                         primitive.ObjectID `bson:"_id,omitempty"`
	Username                   string             `bson:"username"`
	Email                      string             `bson:"email"`
	Password                   string             `bson:"password"`
	Verified                   bool               `bson:"verified"`
	VerificationCode           string             `bson:"verificationCode"`
	RefreshToken               string             `bson:"refreshToken"`
	VerificationExpirationDate time.Time          `bson:"verificationExpirationDate"`
	LastSignIn                 time.Time          `bson:"lastSignIn"`
	CreatedAt                  time.Time          `bson:"createdAt"`
}

func prepareUserForInsertion(user *User) error {
	hashedPassword, err := util.HashAndSaltPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	generatedString, err := util.GenerateRandomString(64)
	if err != nil {
		return err
	}

	user.VerificationCode = generatedString
	user.VerificationExpirationDate = time.Now().AddDate(0, 0, 2)
	user.Verified = false

	return nil
}

func (this *User) Save() (string, error) {
	if err := prepareUserForInsertion(this); err != nil {
		return "", err
	}

	insertion, err := users.InsertOne(context.TODO(), this)

	if err != nil {
		log.Println(err)
		err = util.ConvertMongoErrorToAPIError(err)
		return "", err
	}

	thisID := util.GetObjectIdFromInsertion(insertion.InsertedID)

	if thisID == "" {
		return "", fmt.Errorf("Couldn't obtain the ObjectID for user: %s", this.Username)
	}

	return thisID, nil
}

// Verify the email for a user given the
func (this *User) VerifyEmail(objectID string, verificationCode string) (bool, error) {
	objectIDHex, err := util.GetObjectIDFromString(objectID)

	if err != nil {
		return false, err
	}

	filter := bson.M{
		"_id":                        objectIDHex,
		"verificationCode":           verificationCode,
		"verified":                   false,
		"verificationExpirationDate": bson.D{{"$gt", time.Now()}}}

	update := bson.D{{
		"$set",
		bson.M{"verified": true, "verificationExpirationDate": time.Time{}}}}

	err = users.FindOneAndUpdate(context.TODO(), filter, update).Decode(&this)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Login the user by verifying that the user exists and that the hashed & salted
// password matches the hashed and salted password sent from the login request.
func (this *User) Login() error {
	foundUser := new(User)
	filter := bson.M{"username": this.Username}

	if err := users.FindOne(context.TODO(), filter).Decode(&foundUser); err != nil {
		err = util.ConvertMongoErrorToAPIError(err)
		return err
	}

	if foundUser.Verified == false {
		return errors.New("The user is currently not verified.")
	}

	if success := util.ComparePasswordWithHash(this.Password, foundUser.Password); success != true {
		return errors.New("Could not successfully login.")
	}

	generatedString, err := util.GenerateRandomString(48)

	if err != nil {
		return errors.New("Couldn't issue a refresh token.")
	}

	updateRefreshToken := bson.D{{"$set", bson.M{"refreshToken": generatedString, "lastSignIn": time.Now()}}}

	err = users.FindOneAndUpdate(
		context.TODO(),
		filter,
		updateRefreshToken,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&this)

	if err != nil {
		return err
	}

	return nil
}

// Get a new refresh token for a user given that the current refresh token
// being used is valid.
func (this *User) ValidateAndReplaceToken(receivedRefreshToken string) error {
	foundUser := new(User)

	filter := bson.M{"username": this.Username}

	err := users.FindOne(context.TODO(), filter).Decode(&foundUser)
	if err != nil {
		return err
	}

	if time.Now().Sub(foundUser.LastSignIn) > 5*24*time.Hour {
		return errors.New("User has not signed in for 5 days. Force login.")
	}

	if foundUser.RefreshToken != receivedRefreshToken {
		return errors.New("The refresh tokens do not match. Can't issue new one.")
	}

	generatedString, err := util.GenerateRandomString(48)
	if err != nil {
		return errors.New("Couldn't issue a refresh token.")
	}

	updateRefreshToken := bson.D{{"$set", bson.M{"refreshToken": generatedString}}}
	err = users.FindOneAndUpdate(
		context.TODO(),
		filter,
		updateRefreshToken,
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&this)

	if err != nil {
		return errors.New("Couldn't update the auth token for the new user. Sign in or try again.")
	}

	return nil
}
