package report

import (
	"database/sql"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type ReportStore interface {
	GetReportByID(int) (*Report, error)
	CreateReport(int) (int64, error)
}

type Report struct {
	Id        int       `json:"id"`
	UrlId     int       `json:"url_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) GetReportByID(id int) (*Report, error) {
	response := Report{}

	err := s.db.QueryRow("SELECT id, url_id, updated_at, created_at FROM reports WHERE id = $1", id).
		Scan(&response.Id, &response.UrlId, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Store) CreateReport(id int) (int64, error) {
	res, err := s.db.Exec("INSERT INTO reports (url_id) VALUES ($1)", id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
