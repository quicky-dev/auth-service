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
