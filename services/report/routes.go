package report

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/loczek/tl/internal/cache"
	api_errors "github.com/loczek/tl/internal/errors"
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

type Response struct {
	ReportID int `json:"report_id"`
}

func (h *Handler) ReportLink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	body := new(Payload)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": api_errors.InvalidBodyShape,
		})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": api_errors.InvalidBodyData,
		})
	}

	val, err := h.urlStore.GetUrl(ctx, body.ShortCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": api_errors.DatabaseQuery,
		})
	}

	if val == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": api_errors.NotFound,
		})
	}

	data, err := h.reportStore.CreateReport(ctx, val.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": api_errors.DatabaseQuery,
		})
	}

	metrics.ReportCounter.Add(context.Background(), 1)

	response := Response{
		ReportID: data.ID,
	}

	return c.JSON(response)
}
