package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(){

	mongoURI := os.Getenv("MONGODB_URI")

	if mongoURI == "" {
		log.Fatal("MONGODB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal("MongoDB connection failed: ", err)
	}

	//Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed: ", err)
	}

	log.Println("MongoDB connected successfully!")
	Client = client

}

func OpenCollection(collectionName string) *mongo.Collection {

	databaseName := os.Getenv("DATABASE_NAME")

	if databaseName == "" {
		log.Fatal("DATABASE_NAME not set")
	}

	return Client.Database(databaseName).Collection(collectionName)
}