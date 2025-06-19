package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionRepositoryRedis struct {
	client *redis.Client
}

func NewSessionRepositoryRedis(client *redis.Client) *SessionRepositoryRedis {
	return &SessionRepositoryRedis{client}
}

func (r *SessionRepositoryRedis) SetSession(ctx context.Context, token string, userID string, expired time.Duration) error {
	return r.client.Set(ctx, token, userID, expired).Err()
}

func (r *SessionRepositoryRedis) DeleteSession(ctx context.Context, token string) error {
	return r.client.Del(ctx, token).Err()
}

func (r *SessionRepositoryRedis) CheckSession(ctx context.Context, token string) (string, error) {
	return r.client.Get(ctx, token).Result()
}
