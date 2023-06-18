package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IChatRepository interface {
	GetAll() []*database.Chat
	FindByChatId(chatId string) *database.Chat
	UpdateFeedCheckForUrl(url string, articleFound int, chatId string) bool
	PushNewFeed(url string, chatId string) bool
	Create(chatid string, userid int64, tags []string, feeds []string) *database.Chat
}

type sChatRep struct{}

// Initiate the Chat Repository
func newChatRepository() IChatRepository {
	return &sChatRep{}
}

// Retrieve all Chat from the DB
func (r sChatRep) GetAll() []*database.Chat {
	var chats []*database.Chat

	cur, err := database.ChatCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	defer cur.Close(context.Background())

	// Decod each element found
	for cur.Next(context.Background()) {
		var result database.Chat
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		chats = append(chats, &result)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error while parsing")
		fmt.Println(err.Error())
		return nil
	}

	return chats
}

// Find feed by ChatId
func (r sChatRep) FindByChatId(chatId string) *database.Chat {
	var chat database.Chat

	err := database.ChatCollection.FindOne(context.TODO(), bson.D{{
		Key:   "chatid",
		Value: chatId,
	}}).Decode(&chat)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &chat
}

// Update LastTimeChecked with the current time for a key url / chatId
func (r sChatRep) UpdateFeedCheckForUrl(
	url string,
	articleFound int,
	chatId string,
) bool {
	_, err := database.ChatCollection.UpdateOne(
		context.TODO(),
		bson.D{{
			Key:   "chatid",
			Value: chatId,
		}},
		bson.D{{Key: "$set", Value: bson.M{
			"feeds.$[e].lasttimeparsed": time.Now(),
			"feeds.$[e].articleFound":   articleFound,
		}}},
		options.Update().SetArrayFilters(
			options.ArrayFilters{
				Filters: []interface{}{
					bson.D{{
						Key:   "e.url",
						Value: url,
					}},
				},
			},
		),
	)

	// if error display the error
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// Add a new feed for the chat
func (r sChatRep) PushNewFeed(url string, chatId string) bool {
	_, err := database.ChatCollection.UpdateOne(
		context.TODO(),
		bson.D{{
			Key:   "chatid",
			Value: chatId,
		}},
		bson.D{{
			Key: "$push",
			Value: bson.M{
				"feeds": database.NewFeed(url),
			},
		}},
	)

	// if error display the error
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

// Create a new Chat
func (r sChatRep) Create(
	chatid string,
	userid int64,
	tags []string,
	feedsUrl []string,
) *database.Chat {
	var listOfFeeds []database.Feed
	for _, url := range feedsUrl {
		listOfFeeds = append(listOfFeeds, *database.NewFeed(url))
	}

	chat := database.NewChat(
		chatid,
		userid,
		listOfFeeds,
		tags...,
	)

	_, err := database.ChatCollection.InsertOne(context.TODO(), chat)
	if err != nil {
		return nil
	}

	return chat
}
