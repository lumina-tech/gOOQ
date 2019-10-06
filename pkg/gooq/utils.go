package gooq

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func GetQualifiedName(schema, name string) string {
	return fmt.Sprintf("%s.%s", schema, name)
}

func ScanRow(
	db DBInterface, stmt Fetchable, results interface{},
) error {
	row, err := stmt.FetchRow(Postgres, db)
	if err != nil {
		return err
	}
	if err = row.StructScan(results); err != nil {
		return err
	}
	return nil
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
	row, err := stmt.FetchRow(Postgres, db)
	if err != nil {
		return 0, err
	}
	count := 0
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
