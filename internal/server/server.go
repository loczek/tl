package server

import (
	"context"
	"fmt"
	"log"

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
	err := db.Ping()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to ping postgres: %w", err))
	}

	log.Println("connected to postgres")

	cache := cache.New()
	_, err = cache.Ping(ctx).Result()
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to ping redis: %w", err))
	}

	log.Println("connected to redis")

	app := fiber.New()
	return &Server{
		app,
		db,
		cache,
	}
}
