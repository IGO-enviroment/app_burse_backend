package service

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetById(db *sqlx.DB, id int, table string, scanStruct interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	err = db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table), id).Scan(scanStruct)
	if err == sql.ErrNoRows {
		return fmt.Errorf("record not found")
	} else if err != nil {
		return err
	}

	return nil
}
