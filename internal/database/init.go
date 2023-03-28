package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ArticleCollection *mongo.Collection = nil

func Init() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check if everthing is working
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// TODO: Setup a correct way for dev and server. with Collection / database in the .env
	ArticleCollection = client.Database("otto").Collection("article")
}
