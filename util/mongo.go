package util

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetObjectIdFromInsertion(insertedID interface{}) string {
	objectID, status := insertedID.(primitive.ObjectID)

	if status != true {
		return ""

	}

	return objectID.Hex()
}

func GetObjectIDFromString(objectID string) (primitive.ObjectID, error) {
	objectIDHex, err := primitive.ObjectIDFromHex(objectID)

	if err != nil {
		return primitive.ObjectID{}, err
	}

	return objectIDHex, nil
}

func ConvertMongoErrorToAPIError(err error) error {
	mongoErr := err.(mongo.WriteException)
	mongoErrCode := mongoErr.WriteErrors[0].Code

	switch mongoErrCode {
	case 11000:
		return errors.New("Username already exists.")
	default:
		return errors.New("Internal Server error.")
	}
}
