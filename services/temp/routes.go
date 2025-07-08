package temp

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

func (h *Handler) Temp(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func (h *Handler) Found(c *fiber.Ctx) error {
	// h.logger.Info("Temp route hit")
	return c.SendStatus(200)
}

func (h *Handler) NotFound(c *fiber.Ctx) error {
	return c.SendStatus(404)
}

func (h *Handler) Panic(c *fiber.Ctx) error {
	panic("wtf")
}

func (h *Handler) Err(c *fiber.Ctx) error {
	return fiber.ErrBadGateway
}
