package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

func New(ctx context.Context) *Store {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &Store{
		db,
	}
}

type Url struct {
	Id          string `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalUrl string `json:"original_url"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

type UrlStore interface {
	GetUrl(string) (*Url, error)
	AddUrl(string, string) error
}

func (s *Store) GetUrl(short_code string) (*Url, error) {
	response := Url{}

	err := s.db.QueryRow("SELECT id, short_code, original_url, updated_at, created_at FROM urls WHERE short_code = $1", short_code).
		Scan(&response.Id, &response.ShortCode, &response.OriginalUrl, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Store) AddUrl(short_code string, original_url string) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO urls (short_code, original_url) VALUES ($1, $2) ON CONFLICT (short_code) DO NOTHING")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(short_code, original_url)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
