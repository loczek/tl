package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/loczek/go-link-shortener/internal/config"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("github.com/loczek/go-link-shortener")

type RedisStore struct {
	*redis.Client
}

func New() *RedisStore {
	opts, err := redis.ParseURL(config.REDIS_URL)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to parse redis url: %w", err))
	}

	return &RedisStore{
		redis.NewClient(opts),
	}
}

// Returns empty string when the key doesn't exist instead of an error
func (r *RedisStore) GetCacheKey(ctx context.Context, key string) (string, error) {
	ctx, span := tracer.Start(
		ctx,
		"GET",
		trace.WithAttributes(
			semconv.DBSystemNameRedis,
			attribute.String("key", key),
		),
	)
	defer span.End()

	val, err := r.Get(context.Background(), key).Result()
	if err == nil || err == redis.Nil {
		span.SetStatus(codes.Ok, "")
		return val, nil
	} else {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return val, err
	}
}
