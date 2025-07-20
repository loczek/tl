package shortener

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type UrlStore interface {
	GetUrl(context.Context, string) (*Url, error)
	AddUrl(context.Context, string, string) (int64, error)
}

type Url struct {
	Id          int       `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalUrl string    `json:"original_url"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Store) GetUrl(ctx context.Context, short_code string) (*Url, error) {
	ctx, span := tracer.Start(
		ctx,
		"SELECT url",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	response := Url{}

	err := s.db.QueryRow("SELECT id, short_code, original_url, updated_at, created_at FROM urls WHERE short_code = $1", short_code).
		Scan(&response.Id, &response.ShortCode, &response.OriginalUrl, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "selected")

	return &response, nil
}

func (s *Store) AddUrl(ctx context.Context, short_code string, original_url string) (int64, error) {
	ctx, span := tracer.Start(
		ctx,
		"INSERT url",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	res, err := s.db.Exec("INSERT INTO urls (short_code, original_url) VALUES ($1, $2) ON CONFLICT (short_code) DO NOTHING", short_code, original_url)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	span.SetStatus(codes.Ok, "inserted")

	return count, nil
}
