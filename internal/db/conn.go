package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context) *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABSE_URL is not set in .env")
	}
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// sensible defaults
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Cannot connect Database")
	}

	log.Println("Database connected succesfully")

	return pool
}
