package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type DBArticle struct {
	Id        string    `db:"id"`
	FeedId    string    `db:"feed_id"`
	Title     string    `db:"title"`
	Source    string    `db:"source"`
	Author    string    `db:"author"`
	Link      string    `db:"link"`
	Tags      []string  `db:"tags"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewArticle(
	uuid pgtype.UUID,
	Feedid pgtype.UUID,
	title string,
	source string,
	author string,
	link string,
	tags []string,
	createdAt time.Time,
	updatedAt time.Time,
) *DBArticle {

	var trimedTags []string
	for _, tag := range tags {
		trimedTags = append(trimedTags, strings.TrimSpace(tag))
	}

	return &DBArticle{
		Id:        database.TransformUUIDToString(uuid),
		FeedId:    database.TransformUUIDToString(Feedid),
		Title:     title,
		Source:    source,
		Author:    author,
		Link:      link,
		Tags:      trimedTags,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type IArticleRepository interface {
	GetAll() []*DBArticle
	GetOne(uuid string) *DBArticle
	GetByTitle(title string) *DBArticle
	GetByFeedId(uuid string) []*DBArticle
	Create(feedId string, title string, source string, author string, link string, tags []string) *DBArticle
	Delete(uuid string) bool
	GetByChatAndTime(chatId string) []*DBArticle
}

type SArticleRepository struct {
	conn *pgxpool.Pool
}

// Return all Article in the DB
func (rep *SArticleRepository) GetAll() []*DBArticle {
	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles`
	return rep.query(q)
}

// Return one Article, nil if not found
func (rep *SArticleRepository) GetOne(uuid string) *DBArticle {
	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles where id=$1`
	articles := rep.query(q, uuid)

	if len(articles) == 0 {
		return nil
	}

	return articles[0]
}

// Return one article link to a feed
func (rep *SArticleRepository) GetByFeedId(uuid string) []*DBArticle {
	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles where feed_id=$1`
	return rep.query(q, uuid)
}

// Return one Article by it's title
func (rep *SArticleRepository) GetByTitle(title string) *DBArticle {
	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles where title=$1`
	articles := rep.query(q, title)

	if len(articles) == 0 {
		return nil
	}

	return articles[0]
}

// Create one article
func (rep *SArticleRepository) Create(feedId string, title string, source string, author string, link string, tags []string) *DBArticle {
	q := `INSERT INTO articles (id, feed_id, title, source, author, link, tags) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	newId := uuid.New().String()
	_, err := rep.execute(
		q,
		newId,
		feedId,
		title,
		source,
		author,
		link,
		pq.Array(tags),
	)

	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	return rep.GetOne(newId)
}

// Delete one article from the db
func (rep *SArticleRepository) Delete(uuid string) bool {
	q := `DELETE FROM articles where id=$1`
	res, err := rep.execute(q, uuid)

	// if null throw an error
	if err != nil {
		fmt.Println("Couldn't delete", err.Error())
		return false
	}

	isDelete := false
	if res.RowsAffected() == 1 {
		isDelete = true
	}

	return isDelete
}

// Retrieve new Article based on chat last time parsed
func (rep *SArticleRepository) GetByChatAndTime(chatId string) []*DBArticle {
	var articles []*DBArticle
	q := `
		SELECT
			a.*
		FROM CHATS as c
		INNER JOIN chat_feed as cf
			ON cf.chat_id = c.id
		INNER JOIN articles as a
			ON a.feed_id = cf.feed_id
		LEFT JOIN chat_article as aa
			ON aa.chat_id = c.id
		WHERE c.id = $1
		AND aa.id IS NULL -- Make sure is has been published notified yet
		AND ( a.created_at > c.last_time_parsed OR c.last_time_parsed IS NULL ) -- Check the last time parsed
		AND c.tags && a.tags -- Only pick articles that are parts of chat tags
	`

	rows, err := getConnection().Query(context.Background(), q, chatId)

	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
		return articles
	}

	var id pgtype.UUID
	var feedId pgtype.UUID
	var title string
	var source string
	var author string
	var link string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	params := []any{&id, &feedId, &title, &source, &author, &link, &tags, &createdAt, &updatedAt}
	_, err = pgx.ForEachRow(rows, params, func() error {

		articles = append(
			articles,
			NewArticle(
				id,
				feedId,
				title,
				source,
				author,
				link,
				tags,
				createdAt,
				updatedAt,
			),
		)

		return nil
	})

	if err != nil {
		fmt.Println("Error ForEach", err.Error())
	}

	return articles
}

// Wrap logic for retrieving feeds
func (rep *SArticleRepository) query(q string, param ...any) []*DBArticle {
	var articles []*DBArticle

	rows, err := rep.conn.Query(context.Background(), q, param...)
	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
	}

	var id pgtype.UUID
	var feedId pgtype.UUID
	var title string
	var source string
	var author string
	var link string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	params := []any{&id, &feedId, &title, &source, &author, &link, &tags, &createdAt, &updatedAt}
	_, err = pgx.ForEachRow(rows, params, func() error {

		articles = append(
			articles,
			NewArticle(
				id,
				feedId,
				title,
				source,
				author,
				link,
				tags,
				createdAt,
				updatedAt,
			),
		)

		return nil
	})

	if err != nil {
		fmt.Println("Error ForEach", err.Error())
	}

	return articles
}

// Wrap logic for executing queries
func (rep *SArticleRepository) execute(q string, param ...any) (pgconn.CommandTag, error) {
	res, err := rep.conn.Exec(
		context.Background(),
		q,
		param...,
	)

	return res, err
}
