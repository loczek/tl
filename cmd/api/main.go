package main

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/loczek/go-link-shortener/internal/metrics"
	"github.com/loczek/go-link-shortener/internal/server"
	"go.opentelemetry.io/otel"
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
			log.Error(err)
		}
	}()

	otel.SetMeterProvider(meterProvider)

	log.Info("Starting app")
	ctx := context.Background()
	server := server.New(ctx)
	server.RegisterRoutes()
	log.Info("Started")
	server.Listen(":3000")
}
