package gooq

import (
	"fmt"

	"gopkg.in/guregu/null.v3"
)

type Selectable interface {
	As(alias string) Selectable
	GetAlias() null.String
}

type Table interface {
	Selectable
	GetAliasOrQualifiedName() string
	GetQualifiedName() string
	GetName() string
	GetSchema() string
}

type TableImpl struct {
	name   string
	schema string
	alias  null.String
}

func NewTable(schema, name string) *TableImpl {
	return &TableImpl{
		name:   name,
		schema: schema,
	}
}

func (t *TableImpl) Initialize(schema, name string) {
	t.schema = schema
	t.name = name
}

func (t *TableImpl) As(alias string) Selectable {
	return &TableImpl{
		name:   t.name,
		schema: t.schema,
		alias:  null.StringFrom(alias),
	}
}

func (t TableImpl) GetAlias() null.String {
	return t.alias
}

func (t TableImpl) GetAliasOrQualifiedName() string {
	if t.alias.Valid {
		return t.alias.String
	}
	return t.GetQualifiedName()
}

func (t TableImpl) GetName() string {
	return t.name
}

func (t TableImpl) GetQualifiedName() string {
	return fmt.Sprintf("%s.%s", t.schema, t.name)
}

func (t TableImpl) GetSchema() string {
	return t.schema
}
