package repository

import (
	"context"
	"fmt"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

type IArticleRepository interface {
	Create(title string, published string, link string, source string, author string, match []string, tags ...string) *database.Article
	Find(key string, val string) *database.Article
}

type sArticleRep struct {
}

func newArticleRepository() IArticleRepository {
	return &sArticleRep{}
}

// Create a new Article in the DB
func (r sArticleRep) Create(title string, published string, link string, source string, author string, match []string, tags ...string) *database.Article {
	article := database.NewArticle(title, published, link, source, author, match, tags...)

	_, err := database.ArticleCollection.InsertOne(context.TODO(), article)

	if err != nil {
		return nil
	}

	return article
}

// Find an Artcile by a key
func (r sArticleRep) Find(key string, val string) *database.Article {
	var article *database.Article

	fmt.Println("Searching for Article with key: ", key, " value: ", val)

	err := database.ArticleCollection.FindOne(context.TODO(), bson.D{{Key: key, Value: val}}).
		Decode(&article)

	if err != nil {
		fmt.Println("Error finding:", err)
		return nil
	}

	return article
}
