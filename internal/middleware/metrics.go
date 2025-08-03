package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/loczek/go-link-shortener/internal/telemetry/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func MetricsByName() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		metrics.HttpRequestDuration.Record(c.Context(), float64(time.Since(start).Microseconds())/1000)
		metrics.HttpRequestsCounter.Add(
			c.Context(),
			1,
			metric.WithAttributes(
				attribute.String("name", c.Route().Name),
				attribute.String("path", c.Route().Path),
				attribute.Int("status", c.Response().StatusCode()),
			),
		)

		return err
	}
}
