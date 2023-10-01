package v2

import (
	"context"
	"fmt"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

type DBFeed struct {
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
) *DBFeed {
	return &DBFeed{
		Id:        database.TransformUUIDToString(uuid),
		Url:       url,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type IFeedRepository interface {
	GetAll() []*DBFeed
	GetOne(uuid string) *DBFeed
	Create(url string) *DBFeed
	Delete(uuid string) bool
}

type SFeedRepository struct{}

// Return all Feeds in the DB
func (rep *SFeedRepository) GetAll() []*DBFeed {
	var feeds []*DBFeed

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
func (rep *SFeedRepository) GetOne(uuid string) *DBFeed {
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
func (rep *SFeedRepository) Create(url string) *DBFeed {
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

type MocksFeedRepository struct {
	mock.Mock
}

func (m *MocksFeedRepository) Create(url string) *DBFeed {
	args := m.Called(url)
	return args.Get(0).(*DBFeed)
}

func (m *MocksFeedRepository) Delete(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}

func (m *MocksFeedRepository) GetOne(uuid string) *DBFeed {
	args := m.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBFeed)
}

func (m *MocksFeedRepository) GetAll() []*DBFeed {
	args := m.Called()
	return args.Get(0).([]*DBFeed)
}
