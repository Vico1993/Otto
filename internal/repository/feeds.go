package repository

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
	Disabled  bool      `db:"disabled"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewFeed(
	uuid pgtype.UUID,
	url string,
	disabled bool,
	createdAt time.Time,
	updatedAt time.Time,
) *DBFeed {
	return &DBFeed{
		Id:        database.TransformUUIDToString(uuid),
		Url:       url,
		Disabled:  disabled,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type IFeedRepository interface {
	GetAll() []*DBFeed
	GetAllActive() []*DBFeed
	GetOne(uuid string) *DBFeed
	Create(url string) *DBFeed
	Delete(uuid string) bool
	GetByChatId(uuid string) []string
	LinkChatAndFeed(feedId string, chatId string) bool
	UnLinkChatAndFeed(feedId string, chatId string) bool
	DisableFeed(feedId string) bool
}

type SFeedRepository struct{}

// Return all Feeds in the DB
func (rep *SFeedRepository) GetAll() []*DBFeed {
	q := `SELECT id, url, disabled, created_at, updated_at FROM feeds`

	return rep.retrieveFeeds(q)
}

// Return only active feeds
func (rep *SFeedRepository) GetAllActive() []*DBFeed {
	q := `SELECT id, url, disabled, created_at, updated_at FROM feeds WHERE disabled=FALSE`

	return rep.retrieveFeeds(q)
}

// Return one feed, nil if not found
func (rep *SFeedRepository) GetOne(uuid string) *DBFeed {
	q := `SELECT id, url, disabled, created_at, updated_at FROM feeds where id=$1`

	feeds := rep.retrieveFeeds(q, uuid)

	fmt.Println(feeds)

	if len(feeds) == 0 {
		return nil
	}

	return feeds[0]
}

// Create one feed
func (rep *SFeedRepository) Create(url string) *DBFeed {
	q := `INSERT INTO feeds (id, url) VALUES ($1, $2);`

	newId := uuid.New().String()
	_, err := getConnection().Exec(
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

// Return all Feeds in the DB
func (rep *SFeedRepository) GetByChatId(uuid string) []string {
	var feeds []string

	q := `SELECT
		f.url
		FROM feeds as f
		INNER JOIN chat_feed as cf
			ON cf.feed_id = f.id
		WHERE cf.chat_id=$1`
	rows, err := getConnection().Query(context.Background(), q, uuid)

	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
	}

	var url string
	params := []any{&url}
	_, err = pgx.ForEachRow(rows, params, func() error {
		feeds = append(
			feeds,
			url,
		)

		return nil
	})

	if err != nil {
		fmt.Println("Error ForEach", err.Error())
	}

	return feeds
}

// Subscribed chat to a feed
func (rep *SFeedRepository) LinkChatAndFeed(feedId string, chatId string) bool {
	q := `INSERT INTO chat_feed (chat_id, feed_id) VALUES ($1, $2);`

	_, err := getConnection().Exec(
		context.Background(),
		q,
		chatId,
		feedId,
	)

	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)

		return false
	}

	return true
}

// Delete one feed from the db
func (rep *SFeedRepository) UnLinkChatAndFeed(feedId string, chatId string) bool {
	q := `DELETE FROM chat_feed where chat_id=$1 AND feed_id=$2`
	res, err := getConnection().Exec(context.Background(), q, chatId, feedId)

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

func (rep *SFeedRepository) DisableFeed(feedId string) bool {
	q := `UPDATE feeds SET disabled=TRUE WHERE id=$1;`

	res, err := getConnection().Exec(
		context.Background(),
		q,
		feedId,
	)
	if err != nil {
		fmt.Println("Couldn't disabled feed")
		fmt.Println(err)
		return false
	}

	if res.RowsAffected() == 1 {
		return true
	}

	return false
}

// Wrap logic for retrieving feeds
func (rep *SFeedRepository) retrieveFeeds(q string, param ...any) []*DBFeed {
	var feeds []*DBFeed

	rows, err := getConnection().Query(context.Background(), q, param...)
	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
		return feeds
	}

	var id pgtype.UUID
	var url string
	var disabled bool
	var createdAt time.Time
	var updatedAt time.Time
	params := []any{&id, &url, &disabled, &createdAt, &updatedAt}
	_, err = pgx.ForEachRow(rows, params, func() error {
		feeds = append(
			feeds,
			NewFeed(
				id,
				url,
				disabled,
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

func (m *MocksFeedRepository) GetAllActive() []*DBFeed {
	args := m.Called()
	return args.Get(0).([]*DBFeed)
}

func (m *MocksFeedRepository) GetByChatId(uuid string) []string {
	args := m.Called(uuid)
	return args.Get(0).([]string)
}

func (m *MocksFeedRepository) LinkChatAndFeed(feedId string, chatId string) bool {
	args := m.Called(feedId, chatId)
	return args.Get(0).(bool)
}

func (m *MocksFeedRepository) UnLinkChatAndFeed(feedId string, chatId string) bool {
	args := m.Called(feedId, chatId)
	return args.Get(0).(bool)
}

func (m *MocksFeedRepository) DisableFeed(feedId string) bool {
	args := m.Called(feedId)
	return args.Get(0).(bool)
}
