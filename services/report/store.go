package report

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type ReportStore interface {
	GetReportByID(int) (*Report, error)
	CreateReport(string) (int64, error)
}

type Report struct {
	Id        string `json:"id"`
	UrlId     string `json:"url_id"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
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

func (s *Store) CreateReport(short_code string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO reports (short_code) VALUES ($1)")
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
