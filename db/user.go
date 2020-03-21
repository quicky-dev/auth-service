package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/quicky-dev/auth-service/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

var users = util.MongoClient.Database("Auth").Collection("Users")

type User struct {
	ID                         primitive.ObjectID `bson: "_id, omitempty"`
	Username                   string             `bson:"thisname"`
	Email                      string             `bson:"email"`
	Password                   string             `bson:"password"`
	Verified                   bool               `bson:"verified"`
	VerificationCode           string             `bson:"verificationCode"`
	VerificationExpirationDate time.Time          `bson:"verificationExpirationDate"`
	LastSignIn                 time.Time          `bson:"lastSignIn"`
	CreatedAt                  time.Time          `bson:"createdAt"`
}

func (this *User) Save() (string, error) {
	insertion, err := users.InsertOne(context.TODO(), this)

	if err != nil {
		log.Printf(err.Error())
		log.Printf("Couldn't create the verified this: %s", this.Username)
	}

	thisID := util.GetObjectIdFromInsertion(insertion.InsertedID)

	if thisID == "" {
		return "", fmt.Errorf("Couldn't obtain the ObjectID for this: %s", this.Username)
	}

	return thisID, nil
}

func (this *User) VerifyEmail(objectID string, verificationCode string) (bool, error) {
	objectIDHex, err := util.GetObjectIDFromString(objectID)

	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": objectIDHex, "verificationCode": verificationCode,
		"verified": false, "verificationExpirationDate": bson.D{{"$gt", time.Now()}}}
	update := bson.D{{"$set", bson.M{"verified": true, "verificationExpirationDate": time.Time{}}}}
	log.Println(filter)

	err = users.FindOneAndUpdate(context.TODO(), filter, update).Decode(&this)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (this *User) Login() error {
	filter := bson.M{"username": this.Username}
	foundUser := new(User)

	err := users.FindOne(context.TODO(), filter).Decode(&foundUser)
	if err != nil {
		return err
	}

	success := util.ComparePasswordWithHash(this.Password, foundUser.Password)

	if success {
		this = foundUser
		return nil
	}

	return errors.New("Could not successfully login.")
}
