package report

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/loczek/go-link-shortener/internal/cache"
	metrics "github.com/loczek/go-link-shortener/internal/telemetry"
	"github.com/loczek/go-link-shortener/services/shortener"
)

type Handler struct {
	reportStore ReportStore
	urlStore    shortener.UrlStore
	cache       *cache.RedisStore
}

func NewHandler(db ReportStore, urlStore shortener.UrlStore, cache *cache.RedisStore) *Handler {
	return &Handler{db, urlStore, cache}
}

type Payload struct {
	ShortCode string `json:"short_code" validate:"required,gte=6,lte=8"`
}

func (h *Handler) ReportLink(c *fiber.Ctx) error {
	body := new(Payload)
	if err := c.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		return err
	}

	val, err := h.urlStore.GetUrl(body.ShortCode)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if val == nil {
		return fiber.ErrNotFound
	}

	_, err = h.reportStore.CreateReport(val.Id)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	metrics.ReportCounter.Add(context.Background(), 1)

	return c.SendStatus(200)
}
