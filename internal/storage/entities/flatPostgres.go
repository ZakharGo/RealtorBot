package entities

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
func maxMapCountDate(maps []map[int]time.Time) map[int]time.Time {
	maxMapDate := make(map[int]time.Time)
	var maxDate time.Time
	for _, m := range maps {
		for _, date := range m {
			if date.After(maxDate) {
				maxDate = date
				maxMapDate = m
			}
		}
	}
	return maxMapDate
}
func (f *FlatPostgres) GetAll() ([]string, error) {
	tx, err := f.Flat.Beginx()
	if err != nil {
		return nil, err
	}
	var flats []string
	queryGetAllFlat := fmt.Sprintf("SELECT flat FROM flats")
	if err = tx.Select(&flats, queryGetAllFlat); err != nil {
		tx.Rollback()
		return nil, err
	}
	var res []string
	for _, flat := range flats {
		queryCount := fmt.Sprint("SELECT r.count,r.date from records as r inner join flat_records as fr on r.id = fr.record_id where fr.flat_id=(select id from flats where flat = $1);")
		row, err := tx.Queryx(queryCount, flat)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		maps := make([]map[int]time.Time, 0)

		for row.Next() {
			m := make(map[int]time.Time)
			var c int
			var t time.Time
			if err := row.Scan(&c, &t); err != nil {
				tx.Rollback()
				return nil, err
			}
			m[c] = t
			maps = append(maps, m)

		}

		maxmap := maxMapCountDate(maps)
		for count, date := range maxmap {
			getAllStr := fmt.Sprintf("%s %v  %v", flat, count, date.Format("2006-01-02 15:04:05"))
			res = append(res, getAllStr)
		}

	}
	return res, nil

}
