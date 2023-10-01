package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ArticleCollection *mongo.Collection = nil
var BannedUserCollection *mongo.Collection = nil
var ChatCollection *mongo.Collection = nil

var Connection *pgx.Conn = nil

func TransformUUIDToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

func Init() {
	_ = migrations()

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URI"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	Connection = conn

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
		fmt.Println(err)
		return err
	}
	if err := m.Up(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done migrations")

	return nil
}
