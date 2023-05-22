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
var BannedUserCollection *mongo.Collection = nil
var FeedCollection *mongo.Collection = nil

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

	ArticleCollection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("article")
	BannedUserCollection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("banned_user")
	FeedCollection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("feeds")
}
