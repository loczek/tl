package api

import (
	"database/sql"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberutils "github.com/gofiber/fiber/v2/utils"
	"github.com/loczek/go-link-shortener/internal/cache"
	"github.com/loczek/go-link-shortener/internal/middleware"
	"github.com/loczek/go-link-shortener/services/health"
	"github.com/loczek/go-link-shortener/services/report"
	"github.com/loczek/go-link-shortener/services/shortener"
	"github.com/loczek/go-link-shortener/services/temp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ApiServer struct {
	db    *sql.DB
	cache *cache.RedisStore
}

func NewApiServer(db *sql.DB, cache *cache.RedisStore) *ApiServer {
	return &ApiServer{db, cache}
}

func (s *ApiServer) Run() *fiber.App {
	log := slog.Default().With(slog.String("app", "go-link-shortener"))

	app := fiber.New()

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler())).Name("metrics")

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))
	app.Use(limiter.New(limiter.Config{
		Max: 50,
	}))
	app.Use(recover.New())
	app.Use(requestid.New(requestid.Config{
		Generator: fiberutils.UUIDv4,
	}))

	// api := app.Group("/api")
	// v1 := api.Group("/v1")
	app.Use(fiberlogger.New())
	// app.Use(middleware.Logger())
	app.Use(middleware.MetricsByName())

	app.Static("/", "./website/dist").Name("index")

	app.Use(favicon.New(favicon.Config{
		File: "./website/dist/favicon.ico",
	}))

	healthHandler := health.NewHandler(log.With(slog.String("service", "health")))
	app.Get("/health", healthHandler.Health).Name("health")

	tempHandler := temp.NewHandler(log.With(slog.String("service", "temp")))
	app.Get("/temp", tempHandler.Temp).Name("temp")
	app.Get("/temp/notfound", tempHandler.NotFound).Name("temp.notfound")
	app.Get("/temp/found", tempHandler.Found).Name("temp.found")
	app.Get("/temp/panic", tempHandler.Panic).Name("temp.panic")
	app.Get("/temp/error", tempHandler.Err).Name("temp.err")

	reportStore := report.NewStore(s.db)
	reportHandler := report.NewHandler(reportStore, shortener.NewStore(s.db), s.cache)
	app.Post("/api/report", reportHandler.ReportLink).Name("api.report")

	// shortenerStore := shortener.NewStore(s.db)
	shortenerHandler := shortener.NewHandler(
		shortener.NewStore(s.db),
		s.cache,
		log.With(slog.String("service", "shortener")),
	)
	app.Post("/api/add", shortenerHandler.AddShortenedLink).Name("api.add")
	app.Get("/:hash", shortenerHandler.GetUnshortenedLink).Name("hash")

	return app
}
