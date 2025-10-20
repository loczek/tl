package api

import (
	"database/sql"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberutils "github.com/gofiber/fiber/v2/utils"
	"github.com/loczek/tl/internal/cache"
	"github.com/loczek/tl/internal/config"
	api_errors "github.com/loczek/tl/internal/errors"
	"github.com/loczek/tl/internal/middleware"
	"github.com/loczek/tl/services/health"
	"github.com/loczek/tl/services/report"
	"github.com/loczek/tl/services/shortener"
	"github.com/loczek/tl/services/temp"
)

type ApiServer struct {
	db    *sql.DB
	cache *cache.RedisStore
}

func NewApiServer(db *sql.DB, cache *cache.RedisStore) *ApiServer {
	return &ApiServer{db, cache}
}

func (s *ApiServer) Run() *fiber.App {
	log := slog.Default().With(slog.String("app", "tl-server"))

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,http://short.com",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))
	if config.IsProd() {
		app.Use(limiter.New(limiter.Config{
			Max: 50,
		}))
	}
	app.Use(recover.New())
	app.Use(requestid.New(requestid.Config{
		Generator: fiberutils.UUIDv4,
	}))

	// api := app.Group("/api")
	// v1 := api.Group("/v1")
	app.Use(fiberlogger.New())
	app.Use(middleware.AttachTraceContext())
	app.Use(middleware.Logger())
	app.Use(middleware.MetricsByName())

	healthHandler := health.NewHandler(log.With(slog.String("service", "health")))
	app.Get("/health", healthHandler.Health).Name("health")

	tempHandler := temp.NewHandler(log.With(slog.String("service", "temp")), s.cache)
	app.Get("/temp", tempHandler.Temp).Name("temp")
	app.Get("/temp/calc", tempHandler.Calc).Name("temp.calc")
	app.Get("/temp/cache", tempHandler.Cache).Name("temp.cache")
	app.Get("/temp/notfound", tempHandler.NotFound).Name("temp.notfound")
	app.Get("/temp/found", tempHandler.Found).Name("temp.found")
	app.Get("/temp/panic", tempHandler.Panic).Name("temp.panic")
	app.Get("/temp/error", tempHandler.Err).Name("temp.err")

	reportStore := report.NewStore(s.db)
	reportLogger := log.With(slog.String("service", "report"))
	reportHandler := report.NewHandler(reportStore, shortener.NewStore(s.db), s.cache, reportLogger)
	app.Post("/api/report", reportHandler.ReportLink).Name("api.report")

	shortenerHandler := shortener.NewHandler(
		shortener.NewStore(s.db),
		s.cache,
		log.With(slog.String("service", "shortener")),
	)
	app.Post("/api/add", shortenerHandler.AddShortenedLink).Name("api.add")
	app.Get("/:hash", shortenerHandler.GetUnshortenedLink).Name("hash")

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": api_errors.NotFound,
		})
	})

	return app
}
