package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/loczek/go-link-shortener/internal/config"
	"github.com/loczek/go-link-shortener/internal/metrics"
	"github.com/loczek/go-link-shortener/internal/server"
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
	ctx := context.Background()
	server := server.New(ctx)
	server.RegisterRoutes()
	log.Printf("started on port %d", config.Env.PORT)
	server.Listen(fmt.Sprintf(":%d", config.Env.PORT))
}
