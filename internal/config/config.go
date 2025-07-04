package config

import (
	"log"
	"os"
	"strconv"
)

var (
	ENV           string = getEnvWithFallback("ENV", "production")
	PORT          int    = getEnvAsIntWithFallback("PORT", 3000)
	LOG_TO_STDOUT bool   = true
	DATABASE_URL  string = getEnv("DATABASE_URL")
	REDIS_URL     string = getEnv("REDIS_URL")
)

func IsProd() bool {
	return ENV != "development"
}

func IsDev() bool {
	return !IsProd()
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
		log.Fatalf("missing env var %s\n", key)
	}

	return val
}

func getEnvAsInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing env var %s\n", key)
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("failed to parse key %s with value %s as int\n", key, val)
	}

	return num
}

func getEnvAsIntWithFallback(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("failed to parse key %s with value %s as int\n", key, val)
		}
		return num
	}
	return fallback
}
