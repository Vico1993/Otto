package v2

import (
	"context"
	"fmt"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Feed struct {
	Id        string    `db:"id"`
	Url       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewFeed(
	uuid pgtype.UUID,
	url string,
	createdAt time.Time,
	updatedAt time.Time,
) *Feed {
	return &Feed{
		Id:        database.TransformUUIDToString(uuid),
		Url:       url,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type IFeedRepository interface {
	GetAll() []*Feed
	GetOne(uuid string) *Feed
	Create(url string) *Feed
	Delete(uuid string) bool
}

type SFeedRepository struct{}

// Return all Feeds in the DB
func (rep *SFeedRepository) GetAll() []*Feed {
	var feeds []*Feed

	q := `SELECT id, url, created_at, updated_at FROM feeds`
	rows, err := database.Connection.Query(context.Background(), q)

	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
	}

	var id pgtype.UUID
	var url string
	var createdAt time.Time
	var updatedAt time.Time
	params := []any{&id, &url, &createdAt, &updatedAt}
	_, err = pgx.ForEachRow(rows, params, func() error {
		feeds = append(
			feeds,
			NewFeed(
				id,
				url,
				createdAt,
				updatedAt,
			),
		)

		return nil
	})

	if err != nil {
		fmt.Println("Error ForEach", err.Error())
	}

	return feeds
}

// Return one feed, nil if not found
func (rep *SFeedRepository) GetOne(uuid string) *Feed {
	q := `SELECT id, url, created_at, updated_at FROM feeds where id=$1`

	var id pgtype.UUID
	var url string
	var createdAt time.Time
	var updatedAt time.Time
	err := database.Connection.QueryRow(
		context.Background(),
		q,
		uuid,
	).Scan(
		&id,
		&url,
		&createdAt,
		&updatedAt,
	)

	// if null throw an error
	if err != nil {
		return nil
	}

	return NewFeed(
		id,
		url,
		createdAt,
		updatedAt,
	)
}

// Create one feed
func (rep *SFeedRepository) Create(url string) *Feed {
	q := `INSERT INTO feeds (id, url) VALUES ($1, $2);`

	newId := uuid.New().String()
	_, err := database.Connection.Exec(
		context.Background(),
		q,
		newId,
		url,
	)
	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	return rep.GetOne(newId)
}

// Delete one feed from the db
func (rep *SFeedRepository) Delete(uuid string) bool {
	q := `DELETE FROM feeds where id=$1`
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
