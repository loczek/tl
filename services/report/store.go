package report

import (
	"context"
	"database/sql"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/loczek/tl")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type ReportStore interface {
	GetReportByID(ctx context.Context, id int) (*Report, error)
	CreateReport(ctx context.Context, urlID int) (*Report, error)
}

type Report struct {
	ID        int       `json:"id"`
	URLID     int       `json:"url_id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Store) GetReportByID(ctx context.Context, id int) (*Report, error) {
	ctx, span := tracer.Start(
		ctx,
		"SELECT report",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	response := Report{}

	err := s.db.QueryRow("SELECT id, url_id, updated_at, created_at FROM reports WHERE id = $1", id).
		Scan(&response.ID, &response.URLID, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "")

	return &response, nil
}

func (s *Store) CreateReport(ctx context.Context, urlID int) (*Report, error) {
	ctx, span := tracer.Start(
		ctx,
		"INSERT report",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	response := Report{}

	err := s.db.QueryRow("INSERT INTO reports (url_id) VALUES ($1) RETURNING id, url_id, updated_at, created_at", urlID).
		Scan(&response.ID, &response.URLID, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "")

	return &response, nil
}
