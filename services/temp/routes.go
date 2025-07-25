package temp

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/loczek/go-link-shortener")

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger}
}

func (h *Handler) Temp(c *fiber.Ctx) error {
	ctx, span := tracer.Start(context.Background(), "before")
	defer span.End()

	h.logger.Info("log from temp handler")
	h.logger.InfoContext(ctx, "log from temp handler with ctx")

	expensiveCalculation(ctx, 1, 100_000)
	expensiveCalculationTwo(ctx, 1, 100_000)

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

func expensiveCalculation(ctx context.Context, one int, two int) int {
	ctx, span := tracer.Start(ctx, "calc")
	defer span.End()

	ans := 0

	for i := one; i < two; i++ {
		ans += i
	}

	return ans
}

func expensiveCalculationTwo(ctx context.Context, one int, two int) int {
	ctx, span := tracer.Start(ctx, "lookup user")
	defer span.End()

	ans := 0

	for i := one; i < two; i++ {
		ans += i
	}

	return ans
}
