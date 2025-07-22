package entities

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
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

func (f *FlatPostgres) GetAll() ([]string, []time.Time, error) {
	tx, err := f.Flat.Beginx()
	if err != nil {
		return nil, nil, err
	}
	querySelectAllFlat := fmt.Sprintf("SELECT id, flat FROM flats")
	var numbers []string
	var flatIds []int64
	row := tx.QueryRow(querySelectAllFlat)
	if row.Err() != nil {
		tx.Rollback()
		return nil, nil, row.Err()
	}
	if err := row.Scan(&flatIds, &numbers); err != nil {
		tx.Rollback()
		return nil, nil, row.Err()
	}
	var dates []time.Time
	querySelectDate := fmt.Sprintf("SELECT r.date FROM records as r INNER JOIN flat_records as fr ON fr.record_id = r.id WHERE fr.flat_id = $1 ")
	for _, flatId := range flatIds {
		var date time.Time
		if err := tx.Select(&date, querySelectDate, flatId); err != nil {
			tx.Rollback()
			return nil, nil, err
			break
		}
		dates = append(dates, date)
	}
	tx.Commit()
	return numbers, dates, nil
}
