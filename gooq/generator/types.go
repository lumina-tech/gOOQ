package generator

import (
	"github.com/lumina-tech/gooq/meta"
	"github.com/lumina-tech/gooq/loader"
)

type TablesTemplateArgs struct {
	Timestamp string
	Package   string
	Schema    string
	Tables    []TableTemplateArgs
}

type TableTemplateArgs struct {
	Name        string
	Fields      []FieldTemplateArgs
	Constraints []ConstraintTemplateArgs
}

type ConstraintTemplateArgs struct {
	Name    string
	Columns string
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
	Name   string
	Values []loader.EnumValueMetadata
}

var (
	typesMap map[string]meta.TypeInfo
)

func init() {
	typesMap = make(map[string]meta.TypeInfo)
	for _, typeInfo := range meta.Types {
		typesMap[typeInfo.Prefix] = typeInfo
	}
}
