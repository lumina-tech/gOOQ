package metadata

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	ReferenceTableSuffix = "_reference_table"
)

type Data struct {
	Schema              string
	Tables              []Table
	Enums               []Enum
	ReferenceTableEnums []Enum
	Loader              *Loader
}

type Enum struct {
	Name               string
	Values             []EnumValueMetadata
	IsReferenceTable   bool
	ReferenceTableName string
}

type Table struct {
	Table                 TableMetadata
	Columns               []ColumnMetadata
	Constraints           []ConstraintMetadata
	ForeignKeyConstraints []ForeignKeyConstraintMetadata
}

func NewData(
	db *sqlx.DB, loader *Loader,
) (*Data, error) {
	schema, err := loader.Schema()
	if err != nil {
		return nil, err
	}
	tables, err := getDatabaseTables(db, loader, schema)
	if err != nil {
		return nil, err
	}
	dbEnums, err := getDatabaseEnums(db, loader, schema)
	if err != nil {
		return nil, err
	}
	refTableEnums, err := getReferenceTableEnums(db, loader, schema)
	if err != nil {
		return nil, err
	}
	return &Data{
		Schema:              schema,
		Tables:              tables,
		Enums:               dbEnums,
		ReferenceTableEnums: refTableEnums,
		Loader:              loader,
	}, nil
}

func getDatabaseEnums(
	db *sqlx.DB,
	loader *Loader,
	schema string,
) ([]Enum, error) {
	enums, err := loader.EnumList(db, schema)
	if err != nil {
		return nil, err
	}
	var result []Enum
	for _, enum := range enums {
		enumValues, err := loader.EnumValueList(db, schema, enum.EnumName)
		if err != nil {
			return nil, err
		}
		result = append(result, Enum{
			Name:             enum.EnumName,
			Values:           enumValues,
			IsReferenceTable: false,
		})
	}
	return result, nil
}

func getReferenceTableEnums(
	db *sqlx.DB,
	loader *Loader,
	schema string,
) ([]Enum, error) {
	tables, err := loader.TableList(db, schema)
	if err != nil {
		return nil, err
	}
	var result []Enum
	for _, table := range tables {
		if !strings.HasSuffix(table.TableName, ReferenceTableSuffix) {
			continue
		}
		name := strings.ReplaceAll(table.TableName, ReferenceTableSuffix, "")
		enumValues, err := loader.ReferenceTableValueList(db, schema, table.TableName)
		if err != nil {
			return nil, err
		}
		result = append(result, Enum{
			Name:               name,
			Values:             enumValues,
			IsReferenceTable:   true,
			ReferenceTableName: table.TableName,
		})
	}
	return result, nil
}

func getDatabaseTables(
	db *sqlx.DB,
	loader *Loader,
	schema string,
) ([]Table, error) {
	tables, err := loader.TableList(db, schema)
	if err != nil {
		return nil, err
	}
	var result []Table
	for _, table := range tables {
		columns, err := loader.ColumnList(db, schema, table.TableName)
		if err != nil {
			return nil, err
		}
		constraints, err := loader.ConstraintList(db, schema, table.TableName)
		if err != nil {
			return nil, err
		}
		foreignConstraints, err := loader.ForeignKeyConstraintList(db, table.TableName)
		if err != nil {
			return nil, err
		}
		result = append(result, Table{
			Table:                 table,
			Columns:               columns,
			Constraints:           constraints,
			ForeignKeyConstraints: foreignConstraints,
		})
	}
	return result, nil
}
