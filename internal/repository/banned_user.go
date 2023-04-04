package repository

import (
	"context"
	"fmt"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

// Find an banned user by it's username
func FindBannedUserByTelegramId(id int64) *database.BannedUser {
	var bannedUser *database.BannedUser

	err := database.BannedUserCollection.
		FindOne(context.TODO(), bson.D{{Key: "telegramid", Value: id}}).
		Decode(&bannedUser)
	if err != nil {
		fmt.Println("Couldn't find the user: " + err.Error())
		return nil
	}

	return bannedUser
}

// Create a new Banned User in the DB
func CreateBannedUser(telegramId int64, firstName string, lastName string, userName string, lang string, isBot bool) *database.BannedUser {
	article := database.NewBannedUser(telegramId, firstName, lastName, userName, lang, isBot)

	_, err := database.BannedUserCollection.InsertOne(context.TODO(), article)
	if err != nil {
		fmt.Println("Couldn't insert the user: " + err.Error())
		return nil
	}

	return article
}
