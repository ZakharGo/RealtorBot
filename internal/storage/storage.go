package storage

import (
	"github.com/jmoiron/sqlx"
	"realtorBot/internal/storage/entities"
	"time"
)

type Cache interface {
	Create(data string) error
	Get() (string, error)
	Delete(data string) error
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

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		Cache: entities.NewCachePostgres(db),
		Flat:  entities.NewFlatPostgres(db),
		Count: entities.NewCountPostgres(db),
	}
}
