package generator

import (
	"github.com/lumina-tech/gooq/pkg/generator/metadata"
	"gopkg.in/guregu/null.v3"
)

type TablesTemplateArgs struct {
	Timestamp string
	Package   string
	Schema    string
	Tables    []TableTemplateArgs
}

type TableTemplateArgs struct {
	Name                   string
	DatabaseName           string
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

type EnumsTemplateArgs struct {
	Timestamp string
	Package   string
	Schema    string
	Enums     []EnumType
}

type EnumType struct {
	Name                    string
	Values                  []metadata.EnumValueMetadata
	IsReferenceTable        bool
	ReferenceTableModelType string
}

// TODO(Peter): refactor this
// TODO: this is quite jank
var (
	typesMap map[string]TypeInfo
)

func init() {
	typesMap = make(map[string]TypeInfo)
	for _, typeInfo := range Types {
		typesMap[typeInfo.Prefix] = typeInfo
	}
}

type TypeInfo struct {
	Prefix          string
	Literal         string
	NullableLiteral string
}

var Types = []TypeInfo{
	TypeInfo{Prefix: "Bool", Literal: "bool", NullableLiteral: "null.Bool"},
	TypeInfo{Prefix: "Float32", Literal: "float32", NullableLiteral: "null.Float"},
	TypeInfo{Prefix: "Float64", Literal: "float64", NullableLiteral: "null.Float"},
	TypeInfo{Prefix: "Int", Literal: "int", NullableLiteral: "null.Int"},
	TypeInfo{Prefix: "Int64", Literal: "int64", NullableLiteral: "null.Int"},
	TypeInfo{Prefix: "Jsonb", Literal: "[]byte", NullableLiteral: "nullable.Jsonb"},
	TypeInfo{Prefix: "String", Literal: "string", NullableLiteral: "null.String"},
	TypeInfo{Prefix: "StringArray", Literal: "pq.StringArray", NullableLiteral: "pq.StringArray"},
	TypeInfo{Prefix: "Time", Literal: "time.Time", NullableLiteral: "null.Time"},
	TypeInfo{Prefix: "UUID", Literal: "uuid.UUID", NullableLiteral: "nullable.UUID"},
}
