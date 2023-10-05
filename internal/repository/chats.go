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

type DBChat struct {
	Id             string     `db:"id"`
	TelegramChatId string     `db:"telegram_chat_id"`
	TelegramUserId string     `db:"telegram_user_id, omitempty"`
	Tags           []string   `db:"tags"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	LastTimeParsed *time.Time `db:"last_time_parsed"`
}

func NewChat(
	uuid pgtype.UUID,
	telegramChatId string,
	telegramUserId string,
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
		Id:             database.TransformUUIDToString(uuid),
		TelegramChatId: telegramChatId,
		TelegramUserId: telegramUserId,
		Tags:           trimedTags,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		LastTimeParsed: lastTimeParsed,
	}
}

type IChatRepository interface {
	GetAll() []*DBChat
	GetOne(uuid string) *DBChat
	Create(telegramChatId string, telegramUserId string, tags []string) *DBChat
	Delete(uuid string) bool
	UpdateTags(uuid string, tags []string) bool
	UpdateParsed(uuid string) bool
}

type SChatRepository struct{}

// Return all Chats in the DB
func (rep *SChatRepository) GetAll() []*DBChat {
	var chats []*DBChat

	q := `SELECT id, telegram_chat_id, telegram_user_id, tags, created_at, updated_at, last_time_parsed FROM chats`
	rows, err := database.Connection.Query(context.Background(), q)

	if err != nil {
		fmt.Println("Error Query Execute", err.Error())
	}

	var id pgtype.UUID
	var telegramChatId string
	var telegramUserId string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	var lastTimeParsed *time.Time
	params := []any{&id, &telegramChatId, &telegramUserId, &tags, &createdAt, &updatedAt, &lastTimeParsed}
	_, err = pgx.ForEachRow(rows, params, func() error {

		chats = append(
			chats,
			NewChat(
				id,
				telegramChatId,
				telegramUserId,
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

// Return one chat, nil if not found
func (rep *SChatRepository) GetOne(uuid string) *DBChat {
	q := `SELECT id, telegram_chat_id, telegram_user_id, tags, created_at, updated_at, last_time_parsed FROM chats where id=$1`

	var id pgtype.UUID
	var telegramChatId string
	var telegramUserId string
	var tags []string
	var createdAt time.Time
	var updatedAt time.Time
	var lastTimeParsed *time.Time
	err := database.Connection.QueryRow(
		context.Background(),
		q,
		uuid,
	).Scan(
		&id,
		&telegramChatId,
		&telegramUserId,
		&tags,
		&createdAt,
		&updatedAt,
		&lastTimeParsed,
	)

	// if null throw an error
	if err != nil {
		fmt.Println("Fail Retrieved Chat")
		fmt.Println(err.Error())
		return nil
	}

	return NewChat(
		id,
		telegramChatId,
		telegramUserId,
		tags,
		createdAt,
		updatedAt,
		lastTimeParsed,
	)
}

// Create one chat
func (rep *SChatRepository) Create(telegramChatId string, telegramUserId string, tags []string) *DBChat {
	q := `INSERT INTO chats (id, telegram_chat_id, telegram_user_id, tags) VALUES ($1, $2, $3, $4);`

	newId := uuid.New().String()
	_, err := database.Connection.Exec(
		context.Background(),
		q,
		newId,
		telegramChatId,
		telegramUserId,
		pq.Array(tags),
	)
	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	return rep.GetOne(newId)
}

// Delete one chat from the db
func (rep *SChatRepository) Delete(uuid string) bool {
	q := `DELETE FROM chats where id=$1`
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

// Appends tags to chat
func (rep *SChatRepository) UpdateTags(uuid string, tags []string) bool {
	q := `UPDATE chats SET tags=$1 WHERE id=$2;`

	res, err := database.Connection.Exec(
		context.Background(),
		q,
		pq.Array(tags),
		uuid,
	)
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

	res, err := database.Connection.Exec(
		context.Background(),
		q,
		uuid,
	)
	if err != nil {
		fmt.Println("Couldn't create")
		fmt.Println(err)
	}

	if res.RowsAffected() == 1 {
		return true
	}

	return false
}

type MocksChatRepository struct {
	mock.Mock
}

func (m *MocksChatRepository) Create(telegramChatId string, telegramUserId string, tags []string) *DBChat {
	args := m.Called(telegramChatId, telegramUserId, tags)
	return args.Get(0).(*DBChat)
}

func (m *MocksChatRepository) Delete(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}

func (m *MocksChatRepository) GetOne(uuid string) *DBChat {
	args := m.Called(uuid)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*DBChat)
}

func (m *MocksChatRepository) GetAll() []*DBChat {
	args := m.Called()
	return args.Get(0).([]*DBChat)
}

func (m *MocksChatRepository) GetByChatId(uuid string) []string {
	args := m.Called(uuid)
	return args.Get(0).([]string)
}

func (m *MocksChatRepository) UpdateTags(uuid string, tags []string) bool {
	args := m.Called(uuid, tags)
	return args.Get(0).(bool)
}

func (m *MocksChatRepository) UpdateParsed(uuid string) bool {
	args := m.Called(uuid)
	return args.Get(0).(bool)
}
