package gooq

import "fmt"

type Field interface {
	Name() string
	QualifiedName() string
}

type fieldImpl struct {
	name       string
	selectable Selectable
}

func initFieldImpl(
	selectable Selectable, name string,
) fieldImpl {
	return fieldImpl{
		name:       name,
		selectable: selectable,
	}
}

func (field *fieldImpl) Name() string {
	return field.name
}

func (field *fieldImpl) QualifiedName() string {
	var selectableName string
	switch selectable := field.selectable.(type) {
	case Table:
		selectableName = selectable.GetName()
	default:
		if selectable.GetAlias().Valid {
			selectableName = selectable.GetAlias().String
		}
		// TODO(Peter): can selectable is a anonymous select statement?
	}
	if selectableName == "" {
		return field.name
	} else {
		return fmt.Sprintf("%s.%s", selectableName, field.name)
	}
}

func (field *fieldImpl) Render(
	builder *Builder,
) {
	builder.Print(field.QualifiedName())
}

type DecimalField interface {
	NumericExpression
	Field
}

type defaultDecimalField struct {
	numericExpressionImpl
	fieldImpl
}

func NewDecimalField(
	table Table, name string,
) DecimalField {
	field := &defaultDecimalField{
		fieldImpl: initFieldImpl(table, name),
	}
	field.expressionImpl.initFieldExpressionImpl(field)
	return field
}

type IntField interface {
	NumericExpression
	Field
}

type defaultIntField struct {
	numericExpressionImpl
	fieldImpl
}

func NewIntField(
	table Table, name string,
) IntField {
	field := &defaultIntField{
		fieldImpl: initFieldImpl(table, name),
	}
	field.expressionImpl.initFieldExpressionImpl(field)
	return field
}

type StringField interface {
	StringExpression
	Field
}

type defaultStringField struct {
	stringExpressionImpl
	fieldImpl
}

func NewStringField(
	table Table, name string,
) StringField {
	field := &defaultStringField{
		fieldImpl: initFieldImpl(table, name),
	}
	field.expressionImpl.initFieldExpressionImpl(field)
	return field
}

type UUIDField interface {
	StringExpression
	Field
}

type defaultUUIDField struct {
	stringExpressionImpl
	fieldImpl
}

func NewUUIDField(
	table Table, name string,
) StringField {
	field := &defaultStringField{
		fieldImpl: initFieldImpl(table, name),
	}
	field.expressionImpl.initFieldExpressionImpl(field)
	return field
}

type TimeField interface {
	DateTimeExpression
	Field
}

type defaultTimeField struct {
	dateTimeExpressionImpl
	fieldImpl
}

func NewTimeField(
	table Table, name string,
) TimeField {
	field := &defaultTimeField{
		fieldImpl: initFieldImpl(table, name),
	}
	field.expressionImpl.initFieldExpressionImpl(field)
	return field
}
