package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TransformUUIDToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

func Init() {
	_ = migrations()
}

func migrations() error {
	m, err := migrate.New(
		"file://internal/database/migrations",
		os.Getenv("DB_URI"),
	)
	if err != nil {
		fmt.Println("Couldn't start the migrations process")
		fmt.Println(err)
		return err
	}

	if err := m.Up(); err != nil {
		fmt.Println("UP")
		fmt.Println(err)
	}

	fmt.Println("Done migrations")

	return nil
}
