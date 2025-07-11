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
		slog.Info("route hit",
			slog.String("path", c.Route().Path),
			slog.String("method", c.Route().Method),
			slog.String("time", end.String()),
			slog.String("ip", c.IP()),
		)
		return err
	}
}
