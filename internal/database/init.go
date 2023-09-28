package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ArticleCollection *mongo.Collection = nil
var BannedUserCollection *mongo.Collection = nil
var ChatCollection *mongo.Collection = nil

func Init() {
	_ = migrations()

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
	ChatCollection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("chats")
}

func migrations() error {
	m, err := migrate.New(
		"file://internal/database/migrations",
		os.Getenv("DB_URI"),
	)
	if err != nil {
		fmt.Println("Couldn't start the migrations process")
		return err
	}
	if err := m.Up(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done migrations")

	return nil
}
