package gooq

import (
	"fmt"

	"gopkg.in/guregu/null.v3"
)

type Selectable interface {
	Alias() null.String
	QualifiedName() string
}

type Table interface {
	Selectable
	Name() string
}

type tableImpl struct {
	name  string
	alias null.String
}

func NewTable(name string) Table {
	return &tableImpl{
		name: name,
	}
}

func (t *tableImpl) initTable(
	name string,
) {
	t.name = name
}

func (t tableImpl) Name() string {
	return t.name
}

func (t tableImpl) As(alias string) Selectable {
	return tableImpl{
		name:  t.name,
		alias: null.StringFrom(alias),
	}
}

func (t tableImpl) Alias() null.String {
	return t.alias
}

func (t tableImpl) QualifiedName() string {
	if t.alias.Valid {
		return fmt.Sprintf("%s AS %s", t.Name(), t.alias.String)
	}
	return t.Name()
}
