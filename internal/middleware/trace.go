package middleware

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	fiberutils "github.com/gofiber/fiber/v2/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/loczek/tl")

func AttachTraceContext() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(
			context.Background(),
			fmt.Sprintf("%s",
				string(c.Request().Header.Method()),
			),
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(fiberutils.CopyString(c.Method())),
				// fiber middleware
				semconv.URLScheme(fiberutils.CopyString(c.Protocol())),
				semconv.HTTPRequestBodySize(c.Request().Header.ContentLength()),
				semconv.URLPath(string(fiberutils.CopyBytes(c.Request().URI().Path()))),
				semconv.URLQuery(c.Request().URI().QueryArgs().String()),
				semconv.URLFull(fiberutils.CopyString(c.OriginalURL())),
				semconv.UserAgentOriginal(string(fiberutils.CopyBytes(c.Request().Header.UserAgent()))),
				semconv.ServerAddress(fiberutils.CopyString(c.Hostname())),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		c.SetUserContext(ctx)

		err := c.Next()

		span.SetName(
			fmt.Sprintf("%s %s",
				string(c.Request().Header.Method()),
				c.Route().Path,
			),
		)

		span.SetAttributes(
			semconv.HTTPResponseStatusCode(c.Response().StatusCode()),
		)

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "responded")
		}

		return err
	}
}
