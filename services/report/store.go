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

var tracer = otel.Tracer("github.com/loczek/go-link-shortener")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

type ReportStore interface {
	GetReportByID(context.Context, int) (*Report, error)
	CreateReport(context.Context, int) (int64, error)
}

type Report struct {
	Id        int       `json:"id"`
	UrlId     int       `json:"url_id"`
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
		Scan(&response.Id, &response.UrlId, &response.UpdatedAt, &response.CreatedAt)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "")

	return &response, nil
}

func (s *Store) CreateReport(ctx context.Context, id int) (int64, error) {
	ctx, span := tracer.Start(
		ctx,
		"INSERT report",
		trace.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	res, err := s.db.Exec("INSERT INTO reports (url_id) VALUES ($1)", id)
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

	span.SetStatus(codes.Ok, "")

	return count, nil
}
