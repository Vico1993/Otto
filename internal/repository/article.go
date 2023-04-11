package repository

import (
	"context"
	"fmt"

	"github.com/Vico1993/Otto/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Article sArticleRepository = sArticleRepository{
	collection: database.ArticleCollection,
}

type sArticleRepository struct {
	collection *mongo.Collection
}

// Create a new Article in the DB
func (r sArticleRepository) Create(title string, published string, link string, source string, tags ...string) *database.Article {
	article := database.NewArticle(title, published, link, source, tags...)

	_, err := r.collection.InsertOne(context.TODO(), article)

	if err != nil {
		return nil
	}

	return article
}

// Find an Artcile by a key
func (r sArticleRepository) Find(key string, val string) *database.Article {
	var article *database.Article

	err := r.collection.FindOne(context.TODO(), bson.D{{Key: key, Value: val}}).
		Decode(&article)

	if err != nil {
		fmt.Println("Error finding:", err)
		return nil
	}

	return article
}
