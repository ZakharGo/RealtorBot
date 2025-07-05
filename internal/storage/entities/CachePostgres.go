package entities

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CachePostgres struct {
	Cache *sqlx.DB
}

func NewCachePostgres(cache *sqlx.DB) *CachePostgres {
	return &CachePostgres{Cache: cache}
}
func (c *CachePostgres) Get() (string, error) {
	query := fmt.Sprintf("SELECT script FROM scripts")
	var script string
	if err := c.Cache.Get(&script, query); err != nil {
		return "", err
	}
	return script, nil
}
func (c *CachePostgres) Create(data string) error {
	query := fmt.Sprintf("INSERT INTO scripts (script) values ($1)")
	_, err := c.Cache.Exec(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (c *CachePostgres) Delete(data string) error {
	query := fmt.Sprintf("DELETE FROM scripts where script = $1")
	_, err := c.Cache.Exec(query, data)
	if err != nil {
		return err
	}
	return nil
}
