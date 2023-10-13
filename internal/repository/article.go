package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
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

type SArticleRepository struct{}

// Return all Article in the DB
func (rep *SArticleRepository) GetAll() []*DBArticle {
	var articles []*DBArticle

	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles`
	rows, err := getConnection().Query(context.Background(), q)

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

// Return one Article, nil if not found
func (rep *SArticleRepository) GetOne(uuid string) *DBArticle {
	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles where id=$1`

	var id pgtype.UUID
	var feedId pgtype.UUID
	var title string
	var source string
	var author string
	var link string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	err := getConnection().QueryRow(
		context.Background(),
		q,
		uuid,
	).Scan(
		&id,
		&feedId,
		&title,
		&source,
		&author,
		&link,
		&tags,
		&createdAt,
		&updatedAt,
	)

	// if null throw an error
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return NewArticle(
		id,
		feedId,
		title,
		source,
		author,
		link,
		tags,
		createdAt,
		updatedAt,
	)
}

// Return one article link to a feed
func (rep *SArticleRepository) GetByFeedId(uuid string) []*DBArticle {
	var articles []*DBArticle

	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles where feed_id=$1`
	rows, err := getConnection().Query(context.Background(), q, uuid)

	// if null throw an error
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

// Return one Article by it's title
func (rep *SArticleRepository) GetByTitle(title string) *DBArticle {
	q := `SELECT id, feed_id, title, source, author, link, tags, created_at, updated_at FROM articles where title=$1`

	var id pgtype.UUID
	var feedId pgtype.UUID
	var source string
	var link string
	var author string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	err := getConnection().QueryRow(
		context.Background(),
		q,
		title,
	).Scan(
		&id,
		&feedId,
		&title,
		&source,
		&author,
		&link,
		&tags,
		&createdAt,
		&updatedAt,
	)

	// if null throw an error
	if err != nil {
		return nil
	}

	return NewArticle(
		id,
		feedId,
		title,
		source,
		author,
		link,
		tags,
		createdAt,
		updatedAt,
	)
}

// Create one article
func (rep *SArticleRepository) Create(feedId string, title string, source string, author string, link string, tags []string) *DBArticle {
	q := `INSERT INTO articles (id, feed_id, title, source, author, link, tags) VALUES ($1, $2, $3, $4, $5, $6, $7);`

	newId := uuid.New().String()
	_, err := getConnection().Exec(
		context.Background(),
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
	res, err := getConnection().Exec(context.Background(), q, uuid)

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

type MocksArticleRepository struct {
	mock.Mock
}

func (m *MocksArticleRepository) Create(feedId string, title string, source string, author string, link string, tags []string) *DBArticle {
	args := m.Called(feedId, title, source, author, link, tags)
	return args.Get(0).(*DBArticle)
}

func (m *MocksArticleRepository) Delete(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}

func (m *MocksArticleRepository) GetByTitle(title string) *DBArticle {
	args := m.Called(title)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBArticle)
}

func (m *MocksArticleRepository) GetByFeedId(uuid string) []*DBArticle {
	args := m.Called(uuid)

	return args.Get(0).([]*DBArticle)
}

func (m *MocksArticleRepository) GetOne(uuid string) *DBArticle {
	args := m.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBArticle)
}

func (m *MocksArticleRepository) GetAll() []*DBArticle {
	args := m.Called()
	return args.Get(0).([]*DBArticle)
}

func (m *MocksArticleRepository) GetByChatAndTime(chatId string) []*DBArticle {
	args := m.Called(chatId)

	return args.Get(0).([]*DBArticle)
}
