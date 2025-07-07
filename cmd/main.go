package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/loczek/go-link-shortener/cmd/api"
	"github.com/loczek/go-link-shortener/internal/cache"
	"github.com/loczek/go-link-shortener/internal/config"
	"github.com/loczek/go-link-shortener/internal/db"
	metrics "github.com/loczek/go-link-shortener/internal/telemetry"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("error while loading env file")
	}

	res, err := metrics.NewResource()
	if err != nil {
		panic(err)
	}

	meterProvider, err := metrics.NewMeterProviderPrometheus(res)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	loggerProvider, err := metrics.NewLoggerProvider(res)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	log.Println("initializing")
	database := db.New(context.Background())
	log.Println("connected to postgresql")
	cache := cache.New()
	log.Println("connected to redis")
	server := api.NewApiServer(database, cache).Run()
	log.Printf("started on port %d", config.PORT)
	server.Listen(fmt.Sprintf(":%d", config.PORT))
}
