package shortener

import (
	"database/sql"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type UrlStore interface {
	GetUrl(string) (*Url, error)
	AddUrl(string, string) (int64, error)
}

type Url struct {
	Id          string `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalUrl string `json:"original_url"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
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
	res, err := s.db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2) ON CONFLICT (short_code) DO NOTHING", original_url)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
