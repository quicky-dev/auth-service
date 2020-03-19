package db

import (
	"context"
	"fmt"
	"github.com/quicky-dev/auth-service/util"
	"log"
	"time"
)

var verifiedUsers = util.MongoClient.Database("Accounts").Collection("VerifiedUsers")

type VerifiedUser struct {
	Username   string
	Email      string
	Password   string
	LastSignIn time.Time
	CreatedAt  time.Time
}

func CreateVerifiedUser(user *VerifiedUser) (string, error) {
	user.CreatedAt = time.Now()
	insertion, err := verifiedUsers.InsertOne(context.TODO(), user)

	if err != nil {
		log.Printf("Couldn't create the verified user: %s", user.Username)
		return "", err
	}

	userID := util.GetObjectIdFromInsertion(insertion.InsertedID)

	if userID == "" {
		return "", fmt.Errorf("Couldn't obtain the ObjectID for user: %s", user.Username)
	}

	log.Printf("Successfully created the verified user: %s", user.Username)
	return userID, nil
}
