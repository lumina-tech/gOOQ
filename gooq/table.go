package gooq

type table struct {
	name   string
	fields []Field
	alias  string
}

func Table(name string) TableLike {
	return table{
		name: name,
	}
}

func (t table) IsSelectable() {}

func (t table) Name() string {
	return t.name
}

func (t table) As(alias string) Selectable {
	return table{
		name:   t.name,
		fields: t.fields,
		alias:  alias,
	}
}

func (t table) Alias() string {
	return t.alias
}

func (t table) MaybeAlias() string {
	if t.alias == "" {
		return t.name
	} else {
		return t.alias
	}
}
