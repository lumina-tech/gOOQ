package modelgen

import (
	"github.com/lumina-tech/gooq/pkg/generator/metadata"
	"gopkg.in/guregu/null.v3"
)

type TemplateArgs struct {
	Timestamp string
	Package   string
	Schema    string
	Tables    []TableTemplateArgs
}

type TableTemplateArgs struct {
	TableName              string
	TableType              string
	TableSingletonName     string
	ModelType              string
	QualifiedModelType     string
	ReferenceTableEnumType string
	IsReferenceTable       bool
	Fields                 []FieldTemplateArgs
	Constraints            []ConstraintTemplateArgs
	ForeignKeyConstraints []ForeignKeyConstraintTemplateArgs
}

type ConstraintTemplateArgs struct {
	Name      string
	Columns   string
	Predicate null.String
}

type ForeignKeyConstraintTemplateArgs struct {
	Name string
	ColumnName string
	ForeignTableName string
	ForeignColumnName string
}

type FieldTemplateArgs struct {
	GooqType string
	Name     string
	Type     string
}

type EnumType struct {
	Name                    string
	Values                  []metadata.EnumValueMetadata
	IsReferenceTable        bool
	ReferenceTableModelType string
}
