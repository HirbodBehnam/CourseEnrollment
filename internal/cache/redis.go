package cache

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/go-redis/redis/v8"
)

// NewRedisClient creates a new Redis client and tests it by pinging it
func NewRedisClient(address, password string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return errors.Wrap(err, "cannot ping database")
	}
	return nil
}
