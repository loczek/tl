package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/loczek/go-link-shortener/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	*redis.Client
}

func New() *RedisStore {
	url := config.Env.REDIS_URL

	opts, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to parse redis url: %w", err))
	}

	return &RedisStore{
		redis.NewClient(opts),
	}
}

// Returns empty string when the key doesn't exist instead of an error
func (r *RedisStore) GetCacheKey(ctx context.Context, key string) (string, error) {
	val, err := r.Get(ctx, key).Result()
	if err == nil || err == redis.Nil {
		return val, nil
	} else {
		return val, err
	}
}
