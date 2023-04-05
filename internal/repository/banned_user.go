package repository

import (
	"context"
	"fmt"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

var BannedUser sBannedUserRepository

type sBannedUserRepository struct{}

// Create a new Banned User in the DB
func (r sBannedUserRepository) Create(telegramId int64, firstName string, lastName string, userName string, lang string, isBot bool) *database.BannedUser {
	entity := database.NewBannedUser(telegramId, firstName, lastName, userName, lang, isBot)

	_, err := database.BannedUserCollection.InsertOne(context.TODO(), entity)
	if err != nil {
		return nil
	}

	return entity
}

// Find an banned user by it's username
func (r sBannedUserRepository) Find(key string, val any) *database.BannedUser {
	var bannedUser *database.BannedUser

	err := database.BannedUserCollection.
		FindOne(context.TODO(), bson.D{{Key: key, Value: val}}).
		Decode(&bannedUser)
	if err != nil {
		fmt.Println("Couldn't find the user: " + err.Error())
		return nil
	}

	return bannedUser
}
