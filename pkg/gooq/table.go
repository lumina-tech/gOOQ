package gooq

import (
	"gopkg.in/guregu/null.v3"
)

type Selectable interface {
	GetAlias() null.String
}

type Table interface {
	Selectable
	GetAliasOrName() string
	GetName() string
	GetSchema() string
}

type tableImpl struct {
	name   string
	schema string
	alias  null.String
}

func NewTable(name string) Table {
	return &tableImpl{
		name:   name,
		schema: "public",
	}
}

func (t *tableImpl) initTable(
	name string,
) {
	t.name = name
}

func (t tableImpl) As(alias string) Selectable {
	return tableImpl{
		name:  t.name,
		alias: null.StringFrom(alias),
	}
}

func (t tableImpl) GetAlias() null.String {
	return t.alias
}

func (t tableImpl) GetAliasOrName() string {
	if t.alias.Valid {
		return t.alias.String
	}
	return t.name
}

func (t tableImpl) GetName() string {
	return t.name
}

func (t tableImpl) GetSchema() string {
	return t.schema
}
