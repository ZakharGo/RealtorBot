package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisStorage(options *redis.Options) (*redis.Client, string, error) {
	rdb := redis.NewClient(options)
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, "", err
	}
	return rdb, pong, nil
}
