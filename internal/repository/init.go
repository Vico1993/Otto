package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Chat       IChatRepository
	Feed       IFeedRepository
	Article    IArticleRepository
	connection *pgxpool.Pool
)

func Init() {
	Chat = &SChatRepository{}
	Feed = &SFeedRepository{}
	Article = &SArticleRepository{}

	fmt.Println("Repository Initiated")
}

func getConnection() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(os.Getenv("DB_URI"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create config: %v\n", err)
		os.Exit(1)
	}

	config.MaxConnLifetime = time.Minute * 1

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	connection = conn
	return connection
}
