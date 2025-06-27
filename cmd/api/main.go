package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/loczek/go-link-shortener/internal/metrics"
	"github.com/loczek/go-link-shortener/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	slog.Info("Starting app")
	ctx := context.Background()
	server := server.New(ctx)
	server.RegisterRoutes()
	server.Listen(":3000")
}
