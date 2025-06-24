package server

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/loczek/go-link-shortener/internal/base62"
	"github.com/loczek/go-link-shortener/internal/metrics"
	"github.com/loczek/go-link-shortener/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (s *Server) RegisterRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	s.App.Use(otelfiber.Middleware())

	s.Static("/", "./website/dist")
	s.Post("/api/add", middleware.MetricsByPath("/api/add"), s.AddShortenedLink)
	s.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	s.Get("/temp", middleware.MetricsByPath("/temp"), s.Temp)
	s.Get("/health", middleware.MetricsByPath("/health"), s.Health)
	s.Get("/:hash", middleware.MetricsByPath("/hash"), s.GetUnshortenedLink)
}

func (s *Server) Health(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func (s *Server) Temp(c *fiber.Ctx) error {
	log.Infof("got request!")
	return c.SendStatus(200)
}

func (s *Server) GetUnshortenedLink(c *fiber.Ctx) error {
	hash := c.Params("hash")

	val, err := s.cache.GetCacheKey(context.Background(), fmt.Sprintf("get:%s", hash))
	if err != nil {
		return err
	}

	if val != "" {
		metrics.CacheRequestsCounter.Add(c.Context(), 1, metric.WithAttributes(attribute.String("type", "hit")))
		return c.Redirect(val)
	} else {
		metrics.CacheRequestsCounter.Add(c.Context(), 1, metric.WithAttributes(attribute.String("type", "miss")))

		data, err := s.db.GetUrl(hash)
		if err != nil {
			return err
		}

		err = s.cache.SetEx(context.Background(), fmt.Sprintf("get:%s", hash), data.OriginalUrl, time.Minute).Err()
		if err != nil {
			return err
		}

		return c.Redirect(data.OriginalUrl)
	}
}

type BodyData struct {
	Url string `json:"url" validate:"required,gte=8,lte=1024,url"`
}

type Response struct {
	ShortCode string `json:"short_code"`
}

func (s *Server) AddShortenedLink(c *fiber.Ctx) error {
	body := new(BodyData)
	if err := c.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		return err
	}

	u, err := url.ParseRequestURI(body.Url)
	if err != nil {
		return fiber.ErrBadRequest
	}

	i := 0
	var seq string

	for i < 5 {
		seqInner := base62.RandomSeqRange(6, 8)
		rowsAffectedInner, err := s.db.AddUrl(seqInner, u.String())
		if err != nil {
			return err
		}
		if rowsAffectedInner == 0 {
			metrics.CollisionsCounter.Add(c.Context(), 1)
		} else {
			seq = seqInner
			break
		}
		i += 1
	}

	if seq == "" {
		return fiber.ErrInternalServerError
	}

	return c.JSON(&Response{
		ShortCode: seq,
	})
}
