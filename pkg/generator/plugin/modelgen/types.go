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
	Name                   string
	ModelTableName         string
	RepositoryName         string
	ModelType              string
	ModelTypeWithPackage   string
	ReferenceTableEnumType string
	IsReferenceTable       bool
	Fields                 []FieldTemplateArgs
	Constraints            []ConstraintTemplateArgs
}

type ConstraintTemplateArgs struct {
	Name      string
	Columns   string
	Predicate null.String
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
