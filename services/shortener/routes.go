package shortener

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/loczek/tl/internal/base62"
	"github.com/loczek/tl/internal/cache"
	"github.com/loczek/tl/internal/telemetry/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/loczek/tl")

type Handler struct {
	urlStore UrlStore
	cache    cache.Cache
	logger   *slog.Logger
}

func NewHandler(store UrlStore, cache cache.Cache, logger *slog.Logger) *Handler {
	return &Handler{store, cache, logger}
}

func (h *Handler) GetUnshortenedLink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	hash := c.Params("hash")

	h.logger.InfoContext(ctx, "test", slog.String("name", c.Route().Name), slog.String("path", c.Route().Path))

	val, err := h.cache.GetCacheKey(ctx, fmt.Sprintf("get:%s", hash))
	if err != nil {
		return err
	}

	if val != "" {
		metrics.CacheRequestsCounter.Add(c.Context(), 1, metric.WithAttributes(attribute.String("type", "hit")))
		return c.Redirect(val)
	}

	metrics.CacheRequestsCounter.Add(c.Context(), 1, metric.WithAttributes(attribute.String("type", "miss")))

	data, err := h.urlStore.GetUrl(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.ErrNotFound
		} else {
			h.logger.ErrorContext(ctx, err.Error())
			return fiber.ErrInternalServerError
		}
	}

	err = h.cache.SetCacheKey(ctx, fmt.Sprintf("get:%s", hash), data.OriginalURL, time.Minute)
	if err != nil {
		return err
	}

	return c.Redirect(data.OriginalURL)
}

type Payload struct {
	Url string `json:"url" validate:"required,gte=8,lte=1024,url"`
}

type Response struct {
	ShortCode string `json:"short_code"`
}

func (h *Handler) AddShortenedLink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	span := trace.SpanFromContext(ctx)

	body := new(Payload)
	if err := c.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		return err
	}

	u, err := url.ParseRequestURI(body.Url)
	if err != nil {
		return fiber.ErrBadRequest
	}

	i := 0
	var seq string

	for i < 5 {
		seqInner := base62.RandomSeqRange(6, 8)
		rowsAffectedInner, err := h.urlStore.AddUrl(ctx, seqInner, u.String())
		if err != nil {
			return err
		}
		if rowsAffectedInner == 0 {
			metrics.CollisionsCounter.Add(c.Context(), 1)
		} else {
			seq = seqInner
			break
		}
		i += 1
	}

	if i > 0 {
		span.SetAttributes(attribute.Int("retries", i))
	}

	if seq == "" {
		return fiber.ErrInternalServerError
	}

	return c.JSON(&Response{
		ShortCode: seq,
	})
}
