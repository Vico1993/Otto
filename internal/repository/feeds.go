package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Vico1993/Otto/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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
	GetByChatId(uuid string) []*DBFeed
	LinkChatAndFeed(feedId string, chatId string) bool
	UnLinkChatAndFeed(feedId string, chatId string) bool
	DisableFeed(feedId string) bool
}

type SFeedRepository struct {
	conn *pgxpool.Pool
}

// Return all Feeds in the DB
func (rep *SFeedRepository) GetAll() []*DBFeed {
	q := `SELECT id, url, disabled, created_at, updated_at FROM feeds`

	return rep.query(q)
}

// Return only active feeds
func (rep *SFeedRepository) GetAllActive() []*DBFeed {
	q := `SELECT id, url, disabled, created_at, updated_at FROM feeds WHERE disabled=FALSE`

	return rep.query(q)
}

// Return one feed, nil if not found
func (rep *SFeedRepository) GetOne(uuid string) *DBFeed {
	q := `SELECT id, url, disabled, created_at, updated_at FROM feeds where id=$1`

	feeds := rep.query(q, uuid)

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

// Return all Feeds in the DB
func (rep *SFeedRepository) GetByChatId(uuid string) []*DBFeed {
	q := `SELECT
		f.id,
		f.url,
		f.disabled,
		f.created_at,
		f.updated_at
		FROM feeds as f
		INNER JOIN chat_feed as cf
			ON cf.feed_id = f.id
		WHERE cf.chat_id=$1`

	return rep.query(q, uuid)
}

// Subscribed chat to a feed
func (rep *SFeedRepository) LinkChatAndFeed(feedId string, chatId string) bool {
	q := `INSERT INTO chat_feed (chat_id, feed_id) VALUES ($1, $2);`
	_, err := rep.execute(q, chatId, feedId)

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
	res, err := rep.execute(q, chatId, feedId)

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

	res, err := rep.execute(q, feedId)

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
func (rep *SFeedRepository) query(q string, param ...any) []*DBFeed {
	var feeds []*DBFeed

	rows, err := rep.conn.Query(context.Background(), q, param...)
	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
		return feeds
	}
	defer rows.Close()

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

// Wrap logic for executing queries
func (rep *SFeedRepository) execute(q string, param ...any) (pgconn.CommandTag, error) {
	res, err := rep.conn.Exec(
		context.Background(),
		q,
		param...,
	)

	return res, err
}
