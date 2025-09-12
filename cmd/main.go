package main

import (
	"context"
	"fmt"
	"log"

	"github.com/loczek/tl/cmd/api"
	"github.com/loczek/tl/internal/cache"
	"github.com/loczek/tl/internal/config"
	"github.com/loczek/tl/internal/db"
	"github.com/loczek/tl/internal/telemetry"
)

func main() {
	res, err := telemetry.NewResource()
	if err != nil {
		panic(err)
	}

	meterProvider, err := telemetry.NewMeterProviderHttp(res)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := meterProvider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	loggerProvider, err := telemetry.NewLoggerProvider(res)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := loggerProvider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	tracerProvider, err := telemetry.NewTracerProvider(res)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
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
