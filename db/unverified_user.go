package db

import (
	"context"
	"fmt"
	"github.com/quicky-dev/auth-service/util"
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
