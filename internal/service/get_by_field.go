package service

import (
	"app_burse_backend/pkg/postgres"
	"database/sql"
	"fmt"
)

func GetByField(
	db postgres.Database,
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

	err = db.QueryRowx(
		fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", table, field),
		value,
	).StructScan(scanStruct)
	if err == sql.ErrNoRows {
		return fmt.Errorf("record not found")
	} else if err != nil {
		return err
	}

	return nil
}
