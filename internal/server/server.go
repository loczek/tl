package server

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/loczek/go-link-shortener/internal/cache"
	database "github.com/loczek/go-link-shortener/internal/db"
)

type Server struct {
	*fiber.App
	db    *database.Store
	cache *cache.RedisStore
}

func New(ctx context.Context) *Server {
	db := database.New(ctx)
	cache := cache.New()
	// cache.Ping(ctx)
	app := fiber.New()
	return &Server{
		app,
		db,
		cache,
	}
}
