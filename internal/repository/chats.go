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

type DBChat struct {
	Id               string     `db:"id"`
	TelegramChatId   string     `db:"telegram_chat_id"`
	TelegramUserId   *string    `db:"telegram_user_id, omitempty"`
	TelegramThreadId *string    `db:"telegram_thread_id, omitempty"`
	Tags             []string   `db:"tags"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	LastTimeParsed   *time.Time `db:"last_time_parsed"`
}

func NewChat(
	uuid pgtype.UUID,
	telegramChatId string,
	telegramUserId *string,
	telegramThreadId *string,
	tags []string,
	createdAt time.Time,
	updatedAt time.Time,
	lastTimeParsed *time.Time,
) *DBChat {
	var trimedTags []string
	for _, tag := range tags {
		trimedTags = append(trimedTags, strings.TrimSpace(tag))
	}

	return &DBChat{
		Id:               database.TransformUUIDToString(uuid),
		TelegramChatId:   telegramChatId,
		TelegramUserId:   telegramUserId,
		TelegramThreadId: telegramThreadId,
		Tags:             trimedTags,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		LastTimeParsed:   lastTimeParsed,
	}
}

type IChatRepository interface {
	GetAll() []*DBChat
	GetOne(uuid string) *DBChat
	GetByTelegramChatId(chatId string) *DBChat
	GetByTelegramChatIdAndThreadId(chatId string, threadId string) *DBChat
	Create(telegramChatId string, telegramUserId string, telegramThreadId string, tags []string) *DBChat
	Delete(uuid string) bool
	UpdateTags(uuid string, tags []string) bool
	UpdateParsed(uuid string) bool
}

type SChatRepository struct {
	conn *pgxpool.Pool
}

// Return all Chats in the DB
func (rep *SChatRepository) GetAll() []*DBChat {
	q := `SELECT id, telegram_chat_id, telegram_user_id, telegram_thread_id, tags, created_at, updated_at, last_time_parsed FROM chats`
	return rep.query(q)
}

// Return one chat, nil if not found
func (rep *SChatRepository) GetOne(uuid string) *DBChat {
	q := `SELECT id, telegram_chat_id, telegram_user_id, telegram_thread_id, tags, created_at, updated_at, last_time_parsed FROM chats where id=$1`

	chats := rep.query(q, uuid)

	if len(chats) == 0 {
		return nil
	}

	return chats[0]
}

// Return one chat, nil if not found
func (rep *SChatRepository) GetByTelegramChatId(chatId string) *DBChat {
	q := `SELECT id, telegram_chat_id, telegram_user_id, telegram_thread_id, tags, created_at, updated_at, last_time_parsed FROM chats where telegram_chat_id=$1`

	chats := rep.query(q, chatId)

	if len(chats) == 0 {
		return nil
	}

	return chats[0]
}

// Return one chat by it's id and thread id, nil if not found
func (rep *SChatRepository) GetByTelegramChatIdAndThreadId(chatId string, threadId string) *DBChat {
	q := `SELECT id, telegram_chat_id, telegram_user_id, telegram_thread_id, tags, created_at, updated_at, last_time_parsed FROM chats where telegram_chat_id=$1 and telegram_thread_id=$2`

	chats := rep.query(q, chatId, threadId)
	if len(chats) == 0 {
		return nil
	}

	return chats[0]
}

// Create one chat
func (rep *SChatRepository) Create(
	telegramChatId string,
	telegramUserId string,
	telegramThreadId string,
	tags []string,
) *DBChat {
	q := `INSERT INTO chats (id, telegram_chat_id, telegram_user_id, telegram_thread_id, tags) VALUES ($1, $2, $3, $4, $5);`

	var thread *string
	if telegramThreadId == "" {
		thread = nil
	} else {
		thread = &telegramThreadId
	}

	newId := uuid.New().String()
	_, err := rep.execute(q, newId, telegramChatId, telegramUserId, thread, pq.Array(tags))

	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	return rep.GetOne(newId)
}

// Delete one chat from the db
func (rep *SChatRepository) Delete(uuid string) bool {
	q := `DELETE FROM chats where id=$1`
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

// Appends tags to chat
func (rep *SChatRepository) UpdateTags(uuid string, tags []string) bool {
	q := `UPDATE chats SET tags=$1 WHERE id=$2;`

	res, err := rep.execute(q, pq.Array(tags), uuid)

	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	if res.RowsAffected() == 1 {
		return true
	}

	return false
}

// Update the parsed time
func (rep *SChatRepository) UpdateParsed(uuid string) bool {
	q := `UPDATE chats SET last_time_parsed=NOW() WHERE id=$1;`

	res, err := rep.execute(q, uuid)
	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	if res.RowsAffected() == 1 {
		return true
	}

	return false
}

// Wrap logic for retrieving feeds
func (rep *SChatRepository) query(q string, param ...any) []*DBChat {
	var chats []*DBChat

	rows, err := rep.conn.Query(context.Background(), q, param...)
	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
		return chats
	}
	defer rows.Close()

	var id pgtype.UUID
	var telegramChatId string
	var telegramUserId *string
	var telegramThreadId *string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	var lastTimeParsed *time.Time
	params := []any{&id, &telegramChatId, &telegramUserId, &telegramThreadId, &tags, &createdAt, &updatedAt, &lastTimeParsed}
	_, err = pgx.ForEachRow(rows, params, func() error {

		chats = append(
			chats,
			NewChat(
				id,
				telegramChatId,
				telegramUserId,
				telegramThreadId,
				tags,
				createdAt,
				updatedAt,
				lastTimeParsed,
			),
		)

		return nil
	})
	if err != nil {
		fmt.Println("Error ForEach", err.Error())
	}

	return chats
}

// Wrap logic for executing queries
func (rep *SChatRepository) execute(q string, param ...any) (pgconn.CommandTag, error) {
	res, err := rep.conn.Exec(
		context.Background(),
		q,
		param...,
	)

	return res, err
}
