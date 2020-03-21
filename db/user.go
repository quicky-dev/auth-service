package db

import (
	"context"
	"fmt"
	"github.com/quicky-dev/auth-service/util"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

var unverifiedUsers = util.MongoClient.Database("Accounts").Collection("UnverifiedUsers")

type UnverifiedUser struct {
	Username         string    `bson:"username"`
	Email            string    `bson:"email"`
	Password         string    `bson:"password"`
	VerificationCode string    `bson:"verificationCode"`
	ExpirationDate   time.Time `bson:"expirationDate"`
}

func (this UnverifiedUser) toVerifiedUser() *VerifiedUser {
	verifiedUser := new(VerifiedUser)

	verifiedUser.Username = this.Username
	verifiedUser.Email = this.Email
	verifiedUser.Password = this.Password

	return verifiedUser
}

func CreateUnverifiedUser(user *UnverifiedUser) (string, error) {
	user.ExpirationDate = time.Now().AddDate(0, 0, 1)
	insertion, err := unverifiedUsers.InsertOne(context.TODO(), user)

	if err != nil {
		log.Printf(err.Error())
		log.Printf("Couldn't create the verified user: %s", user.Username)
	}

	userID := util.GetObjectIdFromInsertion(insertion.InsertedID)

	if userID == "" {
		return "", fmt.Errorf("Couldn't obtain the ObjectID for user: %s", user.Username)
	}

	return userID, nil
}

func VerifyEmailForUser(objectID string, verificationCode string) (string, error) {
	objectIDHex, err := util.GetObjectIDFromString(objectID)

	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": objectIDHex, "verificationCode": verificationCode}
	user := new(UnverifiedUser)

	err = unverifiedUsers.FindOneAndDelete(context.TODO(), filter).Decode(user)

	if err != nil {
		return "", err
	}

	userID, err := CreateVerifiedUser(user.toVerifiedUser())

	if err != nil {
		return "", err
	}

	return userID, nil
}
