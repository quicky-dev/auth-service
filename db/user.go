package db

import (
	"context"
	"fmt"
	"github.com/quicky-dev/auth-service/util"
	"log"
	"time"
)

var verifiedUsers = util.MongoClient.Database("Accounts").Collection("VerifiedUsers")
var unverifiedUsers = util.MongoClient.Database("Accounts").Collection("UnverifiedUsers")

type VerifiedUser struct {
	Username   string
	Email      string
	Password   string
	LastSignIn time.Time
	CreatedAt  time.Time
}

type UnverifiedUser struct {
	Username         string    `bson:"username"`
	Email            string    `bson:"email"`
	Password         string    `bson:"password"`
	VerificationCode string    `bson:"verificationCode"`
	ExpirationDate   time.Time `bson:"expirationDate"`
}

func AddUnverifiedUser(user *UnverifiedUser) (string, error) {
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

func AddVerifiedUser(user VerifiedUser) (string, error) {
	user.CreatedAt = time.Now()
	insertion, err := verifiedUsers.InsertOne(context.TODO(), user)

	if err != nil {
		log.Printf("Couldn't create the verified user: %s", user.Username)
		return "", err
	}

	userID := util.GetObjectIdFromInsertion(insertion)

	if userID == "" {
		return "", fmt.Errorf("Couldn't obtain the ObjectID for user: %s", user.Username)
	}

	log.Printf("Successfully created the verified user: %s", user.Username)
	return userID, nil
}
