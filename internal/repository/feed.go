package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

var Feed sFeedRepository

type sFeedRepository struct{}

// Create a new Feed in the DB
func (r sFeedRepository) Create(url string, chatId string) *database.Feed {
	entity := database.NewFeed(url, chatId)

	_, err := database.FeedCollection.InsertOne(context.TODO(), entity)
	if err != nil {
		return nil
	}

	return entity
}

// Find feed by ChatId
func (r sFeedRepository) FindByChatId(chatId string) []database.Feed {
	var feeds []database.Feed

	cursor, err := database.FeedCollection.Find(context.TODO(), bson.D{{
		Key:   "chatid",
		Value: chatId,
	}})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	if err = cursor.All(context.TODO(), &feeds); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return feeds
}

// Update LastTimeChecked with the current time for a key url / chatId
func (r sFeedRepository) SetLastTimeCheck(url string, chatId string) bool {
	_, err := database.FeedCollection.UpdateOne(
		context.TODO(),
		bson.D{{
			Key:   "chatid",
			Value: chatId,
		}, {
			Key:   "url",
			Value: url,
		}},
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "lasttimeparsed",
				Value: time.Now(),
			}},
		}},
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
