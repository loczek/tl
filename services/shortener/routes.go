package shortener

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/loczek/go-link-shortener/internal/base62"
	"github.com/loczek/go-link-shortener/internal/cache"
	metrics "github.com/loczek/go-link-shortener/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Handler struct {
	db     UrlStore
	cache  *cache.RedisStore
	logger *slog.Logger
}

func NewHandler(store UrlStore, cache *cache.RedisStore, logger *slog.Logger) *Handler {
	return &Handler{store, cache, logger}
}

func (h *Handler) GetUnshortenedLink(c *fiber.Ctx) error {
	hash := c.Params("hash")

	h.logger.Info("test", slog.String("name", c.Route().Name), slog.String("path", c.Route().Path))

	val, err := h.cache.GetCacheKey(context.Background(), fmt.Sprintf("get:%s", hash))
	if err != nil {
		return err
	}

	if val != "" {
		metrics.CacheRequestsCounter.Add(c.Context(), 1, metric.WithAttributes(attribute.String("type", "hit")))
		return c.Redirect(val)
	} else {
		metrics.CacheRequestsCounter.Add(c.Context(), 1, metric.WithAttributes(attribute.String("type", "miss")))

		data, err := h.db.GetUrl(hash)
		if err != nil {
			return err
		}

		err = h.cache.SetEx(context.Background(), fmt.Sprintf("get:%s", hash), data.OriginalUrl, time.Minute).Err()
		if err != nil {
			return err
		}

		return c.Redirect(data.OriginalUrl)
	}
}

type Payload struct {
	Url string `json:"url" validate:"required,gte=8,lte=1024,url"`
}

type Response struct {
	ShortCode string `json:"short_code"`
}

func (h *Handler) AddShortenedLink(c *fiber.Ctx) error {
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
		rowsAffectedInner, err := h.db.AddUrl(seqInner, u.String())
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

	if seq == "" {
		return fiber.ErrInternalServerError
	}

	return c.JSON(&Response{
		ShortCode: seq,
	})
}
