package entities

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func maxDateSlice(slice []time.Time) time.Time {
	if len(slice) == 0 {
		return time.Time{} // или panic("empty slice")
	}
	maxt := slice[0]
	for _, val := range slice[1:] {
		if val.After(maxt) {
			maxt = val
		}
	}
	return maxt
}

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
	queryCount := fmt.Sprint("SELECT r.count from records as r inner join flat_records as fr on r.id = fr.record_id where fr.flat_id=(select id from flats where flat = $1);")
	var count []int
	if err := c.Count.Select(&count, queryCount, numb); err != nil {
		return 0, err
	}

	var bestMaxDate time.Time
	var maxCount int
	for _, count := range count {
		var dates []time.Time
		queryDate := fmt.Sprint("SELECT date from records where count = $1")
		if err := c.Count.Select(&dates, queryDate, count); err != nil {
			return 0, err
		}
		maxDate := maxDateSlice(dates)
		if maxDate.After(bestMaxDate) {
			bestMaxDate = maxDate
			maxCount = count
		}
	}
	return maxCount, nil
}

func removeElementMap(slice []map[int]time.Time, value time.Time) []map[int]time.Time {
	for i, maps := range slice {
		for _, date := range maps {
			if date == value {
				return append(slice[:i], slice[i+1:]...)
			}
		}
	}
	return slice // если элемент не найден
}
func maxDateFromMap(maps []map[int]time.Time) (time.Time, int) {
	var maxDate time.Time
	var maxCount int
	for _, m := range maps {
		for count, date := range m {
			if date.After(maxDate) {
				maxDate = date
				maxCount = count
			}
		}
	}
	return maxDate, maxCount
}
func dateAndCountMap(dates []time.Time, count int) []map[int]time.Time {
	sliceMaps := make([]map[int]time.Time, len(dates))
	for _, date := range dates {
		slice := map[int]time.Time{
			count: date,
		}
		sliceMaps = append(sliceMaps, slice)
	}
	return sliceMaps
}
func maxMap(fullCountDateMap []map[int]time.Time) map[int]time.Time {
	maxMap := make(map[int]time.Time)
	maxdate := time.Time{}
	for _, maps := range fullCountDateMap {
		for _, date := range maps {
			if date.After(maxdate) {
				maxdate = date
				maxMap = maps
			}
		}
	}
	return maxMap
}
func (c *CountPostgres) GetPenult(numb string) (int, error) {
	queryCount := fmt.Sprint("SELECT r.count from records as r inner join flat_records as fr on r.id = fr.record_id where fr.flat_id=(select id from flats where flat = $1);")
	var count []int
	if err := c.Count.Select(&count, queryCount, numb); err != nil {
		return 0, err
	}

	fullCountDateMap := make([]map[int]time.Time, len(count))
	for _, count := range count {
		var dates []time.Time
		queryDate := fmt.Sprint("SELECT date from records where count = $1")
		if err := c.Count.Select(&dates, queryDate, count); err != nil {
			return 0, err
		}
		sliceMaps := dateAndCountMap(dates, count)
		fullCountDateMap = append(fullCountDateMap, sliceMaps...)
	}
	maxDate, _ := maxDateFromMap(fullCountDateMap)
	fullCountDateMap = removeElementMap(fullCountDateMap, maxDate)
	_, maxCount := maxDateFromMap(fullCountDateMap)
	return maxCount, nil
}

func (c *CountPostgres) DeleteLastCount() (string, int, error) {
	tx, err := c.Count.Beginx()
	if err != nil {
		return "", 0, fmt.Errorf("error creating transaction: %v", err)
	}
	querySelectCountId := fmt.Sprint("Select id, count from records where date=(select max(date) from records)")
	row := tx.QueryRow(querySelectCountId)
	var countId int64
	var count int
	if err := row.Scan(&countId, &count); err != nil {
		tx.Rollback()
		return "", 0, fmt.Errorf("error scan received countId: %v", err)
	}

	querySelectFlat := fmt.Sprintf("SELECT flat FROM flats as f INNER JOIN flat_records as fr ON f.id =fr.flat_id WHERE fr.record_id=$1")
	row = tx.QueryRow(querySelectFlat, countId)
	if row.Err() != nil {
		tx.Rollback()
		return "", 0, fmt.Errorf("error getting flat deleted count: %v", row.Err())
	}
	var flat string
	if err := row.Scan(&flat); err != nil {
		tx.Rollback()
		return "", 0, fmt.Errorf("error scan flat: %v", err)
	}

	queryDeleteRecord := fmt.Sprint("DELETE from records where date=(select max(date) from records) RETURNING id")
	if _, err := tx.Exec(queryDeleteRecord); err != nil {
		tx.Rollback()
		return "", 0, fmt.Errorf("error deleting last count: %v", err)
	}

	tx.Commit()
	return flat, count, nil
}
