package entities

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type CountPostgres struct {
	Count *sqlx.DB
}

func NewCountPostgres(count *sqlx.DB) *CountPostgres {
	return &CountPostgres{Count: count}
}

func (c *CountPostgres) Create(numb string, count int, date time.Time) error {
	tx, err := c.Count.Beginx()
	var recordId int64
	queryInsertCount := fmt.Sprint("INSERT INTO records(count, date) VALUES($1, $2) RETURNING id;")
	row := tx.QueryRow(queryInsertCount, count, date)
	if err = row.Scan(&recordId); err != nil {
		tx.Rollback()
		return err
	}
	queryGetFlatId := fmt.Sprint("SELECT id FROM flats WHERE flat = $1")
	var flatId int64
	if err := tx.Get(&flatId, queryGetFlatId, numb); err != nil {
		tx.Rollback()
		return err
	}
	query := fmt.Sprint("INSERT INTO flat_records(flat_Id, record_Id) VALUES($1, $2)")
	_, err = tx.Exec(query, flatId, recordId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (c *CountPostgres) GetAll() ([]int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CountPostgres) GetLast(numb string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CountPostgres) GetPenult(numb string) (int, error) {
	//TODO implement me
	panic("implement me")
}
