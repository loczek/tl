package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		end := time.Since(start)
		slog.Info("",
			slog.String("method", c.Route().Method),
			slog.String("path", c.Route().Path),
			slog.String("time", end.String()),
			slog.String("ip", c.IP()),
			slog.String("id", c.Locals("requestid").(string)),
		)
		return err
	}
}
