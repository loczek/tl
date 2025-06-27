package config

import (
	"log"
	"os"
	"strconv"
)

var Env = New()

type config struct {
	PORT         int
	DATABASE_URL string
	REDIS_URL    string
}

func New() config {
	return config{
		PORT:         getEnvAsIntWithFallback("PORT", 3000),
		DATABASE_URL: getEnv("DATABASE_URL"),
		REDIS_URL:    getEnv("REDIS_URL"),
	}
}

func getEnvWithFallback(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("ERROR: missing env var %s", key)
	}

	return val
}

func getEnvAsInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("ERROR: missing env var %s", key)
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("ERROR: failed to parse key %s with value %s as int", key, val)
	}

	return num
}

func getEnvAsIntWithFallback(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("ERROR: failed to parse key %s with value %s as int", key, val)
		}
		return num
	}
	return fallback
}
