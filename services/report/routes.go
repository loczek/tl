package report

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/loczek/tl/internal/cache"
	"github.com/loczek/tl/internal/telemetry/metrics"
	"github.com/loczek/tl/services/shortener"
)

type Handler struct {
	reportStore ReportStore
	urlStore    shortener.UrlStore
	cache       cache.Cache
}

func NewHandler(db ReportStore, urlStore shortener.UrlStore, cache cache.Cache) *Handler {
	return &Handler{db, urlStore, cache}
}

type Payload struct {
	ShortCode string `json:"short_code" validate:"required,gte=6,lte=8"`
}

func (h *Handler) ReportLink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	body := new(Payload)
	if err := c.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		return err
	}

	val, err := h.urlStore.GetUrl(ctx, body.ShortCode)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if val == nil {
		return fiber.ErrNotFound
	}

	_, err = h.reportStore.CreateReport(ctx, val.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	metrics.ReportCounter.Add(context.Background(), 1)

	return c.SendStatus(200)
}
