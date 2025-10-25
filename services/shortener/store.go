package shortener

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type UrlStore interface {
	GetUrl(ctx context.Context, shortCode string) (*URL, error)
	AddUrl(ctx context.Context, shortCode string, originalURL string) (int64, error)
}

type URL struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Store) GetUrl(ctx context.Context, shortCode string) (*URL, error) {
	ctx, span := tracer.Start(
		ctx,
		"SELECT url",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	response := URL{}

	err := s.db.QueryRow("SELECT id, short_code, original_url, updated_at, created_at FROM urls WHERE short_code = $1", shortCode).
		Scan(&response.ID, &response.ShortCode, &response.OriginalURL, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "selected")

	return &response, nil
}

func (s *Store) AddUrl(ctx context.Context, shortCode string, originalURL string) (int64, error) {
	ctx, span := tracer.Start(
		ctx,
		"INSERT url",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	res, err := s.db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2) ON CONFLICT (short_code) DO NOTHING", shortCode, originalURL)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return 0, err
	}

	span.SetStatus(codes.Ok, "inserted")

	return count, nil
}
