package cache

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	*redis.Client
}

func New() *RedisStore {
	url := os.Getenv("REDIS_URL")

	log.Infof("url: %s", url)
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
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
