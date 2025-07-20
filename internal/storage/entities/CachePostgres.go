package entities

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	redisKeyTTL = 1 * time.Minute
)

type CacheRedis struct {
	Cache *redis.Client
}

func NewCacheRedis(cache *redis.Client) *CacheRedis {
	return &CacheRedis{Cache: cache}
}

func (c *CacheRedis) Create(ctx context.Context, data string, userID string) error {
	err := c.Cache.Set(ctx, userID, data, redisKeyTTL).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheRedis) Get(ctx context.Context, userID string) (string, error) {
	val, err := c.Cache.GetDel(ctx, userID).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s does not exist", userID)
	} else if err != nil {
		return "", fmt.Errorf("error getting item: %v", err)
	}
	return val, nil
}

//
//func (c *CacheRedis) Delete(ctx context.Context, userID string) error {
//	_, err := c.Cache.Del(ctx, userID).Result()
//	if err != nil {
//		return fmt.Errorf("error deleting item: %v", err)
//	}
//	return nil
//}
