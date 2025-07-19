package entities

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type FlatPostgres struct {
	Flat *sqlx.DB
}

func NewFlatPostgres(flat *sqlx.DB) *FlatPostgres {
	return &FlatPostgres{Flat: flat}
}

func (f *FlatPostgres) Create(numb string) error {
	query := "INSERT INTO flats (flat) values ($1)"
	_, err := f.Flat.Exec(query, numb)
	if err != nil {
		return err
	}
	return nil
}

func (f *FlatPostgres) Delete(numb string) error {
	query := "DELETE FROM flats WHERE flat = $1"
	res, err := f.Flat.Exec(query, numb)
	if err != nil {
		return err
	}
	k, err := res.RowsAffected()
	if k == 0 {
		return errors.New("Flat not found")
	}
	return nil
}

func (f *FlatPostgres) GetAll() ([]string, error) {
	query := fmt.Sprintf("SELECT flat FROM flats")
	var numbers []string
	if err := f.Flat.Select(&numbers, query); err != nil {
		return nil, err
	}
	return numbers, nil
}
