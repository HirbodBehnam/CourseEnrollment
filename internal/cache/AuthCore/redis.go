package AuthCore

import (
	"CourseEnrollment/internal/shared"
	"context"
	"github.com/go-redis/redis/v8"
)

// Prefix for course enrollment auth
const prefix = "cea:"

// RedisAuth is a token storage for authentication based on Redis
type RedisAuth struct {
	client *redis.Client
}

// NewRedisAuth will create a RedisAuth based on a redis.Client
func NewRedisAuth(client *redis.Client) RedisAuth {
	return RedisAuth{client}
}

func (a RedisAuth) Set(ctx context.Context, token string) error {
	return a.client.Set(ctx, prefix+token, 0, shared.AuthCoreTokenTTL).Err()
}

func (a RedisAuth) Delete(ctx context.Context, token string) error {
	return a.client.Del(ctx, prefix+token).Err()
}

func (a RedisAuth) IsValid(ctx context.Context, token string) (bool, error) {
	result := a.client.Exists(ctx, prefix+token)
	if err := result.Err(); err != nil {
		return false, err
	}
	keys, err := result.Result()
	if err != nil {
		return false, err
	}
	return keys == 1, nil
}

func (a RedisAuth) Replace(ctx context.Context, old, new string) error {
	if err := a.Delete(ctx, old); err != nil {
		return err
	}
	return a.Set(ctx, new)
}
