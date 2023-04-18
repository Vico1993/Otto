package repository

import (
	"context"
	"fmt"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

var User IUserRepository = &sUserRepository{}

type sUserRepository struct{}

type IUserRepository interface {
	Create(chatId int64, telegramId int64, firstName string, lastName string, userName string, lang string, isBot bool, isBanned bool) *database.User
	FindById(id int64) *database.User
}

// Create a new user
func (s *sUserRepository) Create(chatId int64, telegramId int64, firstName string, lastName string, userName string, lang string, isBot bool, isBanned bool) *database.User {
	entity := database.NewUser(chatId, telegramId, firstName, lastName, userName, lang, isBot, isBanned)

	_, err := database.UserCollection.InsertOne(context.TODO(), entity)
	if err != nil {
		return nil
	}

	return entity
}

// Find user by it's id
func (s *sUserRepository) FindById(id int64) *database.User {
	var user *database.User

	err := database.UserCollection.
		FindOne(context.TODO(), bson.D{{Key: "telegramId", Value: id}}).
		Decode(&user)
	if err != nil {
		fmt.Println("Couldn't find the user: " + err.Error())
		return nil
	}

	return user
}
