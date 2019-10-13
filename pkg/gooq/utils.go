package gooq

import (
	"github.com/jmoiron/sqlx"
)

func ScanRow(
	db DBInterface, stmt Fetchable, results interface{},
) error {
	row := stmt.FetchRow(Postgres, db)
	return row.StructScan(results)
}

func ScanRows(
	db DBInterface, stmt Fetchable, results interface{},
) error {
	rows, err := stmt.Fetch(Postgres, db)
	if err != nil {
		return err
	}
	defer rows.Close()
	return sqlx.StructScan(rows, results)
}

func ScanCount(
	db DBInterface, stmt Fetchable,
) (int, error) {
	row := stmt.FetchRow(Postgres, db)
	count := 0
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
