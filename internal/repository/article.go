package repository

import (
	"context"
	"fmt"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

// Find an Artcile by it's title
func FindArticleByTitle(title string) *database.Article {
	var article *database.Article

	err := database.ArticleCollection.
		FindOne(context.TODO(), bson.D{{Key: "title", Value: title}}).
		Decode(&article)
	if err != nil {
		fmt.Println("Couldn't find the article: " + err.Error())
		return nil
	}

	return article
}

// Create a new Article in the DB
func CreateArticle(title string, published string, link string, source string, tags ...string) *database.Article {
	article := database.NewArticle(title, published, link, source, tags...)

	_, err := database.ArticleCollection.InsertOne(context.TODO(), article)
	if err != nil {
		fmt.Println("Couldn't find the article: " + err.Error())
		return nil
	}

	return article
}
