package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *postgres
	pgOnce     sync.Once
	Chat       IChatRepository
	Feed       IFeedRepository
	Article    IArticleRepository
)

func Init() {
	pg := newConnectionPool(context.Background())

	Chat = &SChatRepository{
		conn: pg.db,
	}
	Feed = &SFeedRepository{
		conn: pg.db,
	}
	Article = &SArticleRepository{
		conn: pg.db,
	}

	fmt.Println("Repository Initiated")
}

type postgres struct {
	db *pgxpool.Pool
}

func newConnectionPool(ctx context.Context) *postgres {
	config, err := pgxpool.ParseConfig(os.Getenv("DB_URI"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create config: %v\n", err)
		os.Exit(1)
	}

	config.MaxConnLifetime = time.Minute * 1

	pgOnce.Do(func() {
		db, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}

		pgInstance = &postgres{db}
	})

	return pgInstance
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}
