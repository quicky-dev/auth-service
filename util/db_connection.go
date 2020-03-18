package util

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func obtainDatabaseClient() *mongo.Client {
	credentials := new(options.Credential)

	if os.Getenv("PRODUCTION_MODE") == "" {
		credentials.Username = os.Getenv("AUTH_DB_USER")
		credentials.Password = os.Getenv("AUTH_DB_PASSWORD")
	}

	clientOptions := new(options.ClientOptions)
	clientOptions.Auth = credentials
	clientOptions.ApplyURI(os.Getenv("MONGODB_URI"))

	err := clientOptions.Validate()

	if err != nil {
		log.Fatal(err)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)
	log.Println("Successfully connected to mongo.")

	return client
}

var MongoClient = obtainDatabaseClient()
