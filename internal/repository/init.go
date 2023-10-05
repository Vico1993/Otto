package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var (
	Chat    IChatRepository
	Feed    IFeedRepository
	Article IArticleRepository
)

func Init() {
	Chat = &SChatRepository{}
	Feed = &SFeedRepository{}
	Article = &SArticleRepository{}

	fmt.Println("Repository Initiated")
}

func getConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URI"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
