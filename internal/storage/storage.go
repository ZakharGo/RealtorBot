package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"realtorBot/internal/storage/entities"
	"time"
)

type Cache interface {
	Create(ctx context.Context, data string, userID string) error
	Get(ctx context.Context, userID string) (string, error)
	Delete(ctx context.Context, userID string) error
}
type Flat interface {
	Create(numb string) error
	GetAll() ([]string, error)
}
type Count interface {
	Create(numb string, count int, date time.Time) error
	GetAll() ([]int, error)
	GetLast(numb string) (int, error)
	GetPenult(numb string) (int, error)
}
type Storage struct {
	Cache
	Flat
	Count
}

func NewStorage(pdb *sqlx.DB, rdb *redis.Client) *Storage {
	return &Storage{
		Cache: entities.NewCacheRedis(rdb),
		Flat:  entities.NewFlatPostgres(pdb),
		Count: entities.NewCountPostgres(pdb),
	}
}
