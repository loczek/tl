package temp

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/loczek/tl/internal/cache"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/loczek/tl")

type Handler struct {
	logger *slog.Logger
	cache  cache.Cache
}

func NewHandler(logger *slog.Logger, cache cache.Cache) *Handler {
	return &Handler{logger, cache}
}

func (h *Handler) Temp(c *fiber.Ctx) error {
	ctx := c.UserContext()

	ctx, span := tracer.Start(ctx, "before")
	defer span.End()

	h.logger.Info("log from temp handler")
	h.logger.InfoContext(ctx, "log from temp handler with ctx")

	out := expensiveCalculation(ctx, 1, rand.IntN(100_000_000))

	return c.SendString(fmt.Sprint(out))
}

func (h *Handler) Calc(c *fiber.Ctx) error {
	out := rand.IntN(100_000_000)
	return c.SendString(fmt.Sprint(out))
}

func (h *Handler) Cache(c *fiber.Ctx) error {
	ctx := c.UserContext()

	val, err := h.cache.GetCacheKey(ctx, "temp:num")
	if err != nil {
		return err
	}

	if val != "" {
		return c.SendString(val)
	}

	out := rand.IntN(100_000_000)

	err = h.cache.SetCacheKey(ctx, "temp:num", fmt.Sprint(out), time.Minute*3)
	if err != nil {
		return err
	}

	return c.SendString(fmt.Sprint(out))
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
