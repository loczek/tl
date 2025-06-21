package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/loczek/go-link-shortener/internal/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func MetricsByPath(path string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		metrics.HttpRequestsCounter.Add(
			c.Context(),
			1,
			metric.WithAttributes(
				attribute.String("path", path),
				attribute.Int("status", c.Response().StatusCode()),
			),
		)

		return err
	}
}

func MetricsByRoutePath() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		metrics.HttpRequestsCounter.Add(
			c.Context(),
			1,
			metric.WithAttributes(
				attribute.String("path", c.Route().Path),
				attribute.Int("status", c.Response().StatusCode()),
			),
		)

		return err
	}
}
