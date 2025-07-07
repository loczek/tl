package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/loczek/go-link-shortener/internal/config"
)

func New(ctx context.Context) *sql.DB {
	db, err := sql.Open("pgx", config.DATABASE_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return db
}
