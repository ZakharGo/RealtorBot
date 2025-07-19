package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	UserName string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=%s", cfg.UserName, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
