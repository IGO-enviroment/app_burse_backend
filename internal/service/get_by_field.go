package service

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetByField(
	db *sqlx.DB,
	field,
	value interface{},
	table string,
	scanStruct interface{},
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	err = db.QueryRow(
		fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", table, field),
		value,
	).Scan(scanStruct)
	if err == sql.ErrNoRows {
		return fmt.Errorf("record not found")
	} else if err != nil {
		return err
	}

	return nil
}
