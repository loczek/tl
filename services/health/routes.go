package health

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger}
}

func (h *Handler) Health(c *fiber.Ctx) error {
	h.logger.Info("Health route hit")
	return c.SendStatus(200)
}
