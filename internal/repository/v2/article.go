package v2

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
)

type dbArticle struct {
	Id        string    `db:"id"`
	FeedId    string    `db:"feed_id"`
	Title     string    `db:"title"`
	Source    string    `db:"source"`
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
	link string,
	tags []string,
	createdAt time.Time,
	updatedAt time.Time,
) *dbArticle {

	var trimedTags []string
	for _, tag := range tags {
		trimedTags = append(trimedTags, strings.TrimSpace(tag))
	}

	return &dbArticle{
		Id:        database.TransformUUIDToString(uuid),
		FeedId:    database.TransformUUIDToString(Feedid),
		Title:     title,
		Source:    source,
		Link:      link,
		Tags:      trimedTags,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type IArticleRepository interface {
	GetAll() []*dbArticle
	GetOne(uuid string) *dbArticle
	GetByFeedId(uuid string) []*dbArticle
	Create(feedId string, title string, source string, link string, tags []string) *dbArticle
	Delete(uuid string) bool
}

type SArticleRepository struct{}

// Return all Article in the DB
func (rep *SArticleRepository) GetAll() []*dbArticle {
	var articles []*dbArticle

	q := `SELECT id, feed_id, title, source, link, tags, created_at, updated_at FROM articles`
	rows, err := database.Connection.Query(context.Background(), q)

	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
	}

	var id pgtype.UUID
	var feedId pgtype.UUID
	var title string
	var source string
	var link string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	params := []any{&id, &feedId, &title, &source, &link, &tags, &createdAt, &updatedAt}
	_, err = pgx.ForEachRow(rows, params, func() error {

		articles = append(
			articles,
			NewArticle(
				id,
				feedId,
				title,
				source,
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
func (rep *SArticleRepository) GetOne(uuid string) *dbArticle {
	q := `SELECT id, feed_id, title, source, link, tags, created_at, updated_at FROM articles where id=$1`

	var id pgtype.UUID
	var feedId pgtype.UUID
	var title string
	var source string
	var link string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	err := database.Connection.QueryRow(
		context.Background(),
		q,
		uuid,
	).Scan(
		&id,
		&feedId,
		&title,
		&source,
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
		link,
		tags,
		createdAt,
		updatedAt,
	)
}

// Return one Article, nil if not found
func (rep *SArticleRepository) GetByFeedId(uuid string) []*dbArticle {
	var articles []*dbArticle

	q := `SELECT id, feed_id, title, source, link, tags, created_at, updated_at FROM articles where feed_id=$1`
	rows, err := database.Connection.Query(context.Background(), q, uuid)

	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
	}

	var id pgtype.UUID
	var feedId pgtype.UUID
	var title string
	var source string
	var link string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	params := []any{&id, &feedId, &title, &source, &link, &tags, &createdAt, &updatedAt}
	_, err = pgx.ForEachRow(rows, params, func() error {

		articles = append(
			articles,
			NewArticle(
				id,
				feedId,
				title,
				source,
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

// Create one article
func (rep *SArticleRepository) Create(feedId string, title string, source string, link string, tags []string) *dbArticle {
	q := `INSERT INTO articles (id, feed_id, title, source, link, tags) VALUES ($1, $2, $3, $4, $5, $6);`

	newId := uuid.New().String()
	_, err := database.Connection.Exec(
		context.Background(),
		q,
		newId,
		feedId,
		title,
		source,
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
	res, err := database.Connection.Exec(context.Background(), q, uuid)

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
