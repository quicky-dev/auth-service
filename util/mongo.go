package util

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
