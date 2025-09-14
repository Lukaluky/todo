package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	TasksColl *mongo.Collection
)

func ConnectMongo() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

    if err != nil {
		log.Fatalf("Mongo connect error : %s", err)
	}

	MongoClient = client
	dbName := os.Getenv("MONGO_DB")
	collName := os.Getenv("MONGO_TASKS_COLLECTION")

	if dbName == "" {
		dbName = "tododb"
	}
	if collName == "" {
		collName = "tasks"
	}

	TasksColl = MongoClient.Database(dbName).Collection(collName)




}