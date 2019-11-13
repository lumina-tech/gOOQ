package gooq

import "fmt"

type Field interface {
	Named
}

type fieldImpl struct {
	expression Expression
	selectable Selectable
	name       string
}

func (field *fieldImpl) initFieldImpl(
	expression Expression, selectable Selectable, name string,
) {
	field.expression = expression
	field.selectable = selectable
	field.name = name
}

func (field *fieldImpl) GetName() string {
	return field.name
}

func (field *fieldImpl) GetQualifiedName() string {
	var selectableName string
	switch selectable := field.selectable.(type) {
	case Table:
		selectableName = selectable.GetUnqualifiedName()
	case *selection:
		if selectable.GetAlias().Valid {
			selectableName = selectable.GetAlias().String
		}
		// TODO(Peter): can selectable be a anonymous select statement?
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
	builder.Print(field.GetQualifiedName())
}

// BoolField

type BoolField interface {
	BoolExpression
	Field
}

type defaultBoolField struct {
	boolExpressionImpl
	fieldImpl
}

func NewBoolField(
	table Table, name string,
) BoolField {
	field := &defaultBoolField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// DecimalField

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
	field := &defaultDecimalField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// IntField

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
	field := &defaultIntField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// JsonbField

type JsonbField interface {
	StringExpression
	Field
}

type defaultJsonbField struct {
	stringExpressionImpl
	fieldImpl
}

func NewJsonbField(
	table Table, name string,
) JsonbField {
	field := &defaultJsonbField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// StringField

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
	field := &defaultStringField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// StringArrayField

type StringArrayField interface {
	StringExpression
	Field
}

type defaultStringArrayField struct {
	stringExpressionImpl
	fieldImpl
}

func NewStringArrayField(
	table Table, name string,
) StringArrayField {
	field := &defaultStringField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// UUIDField

type UUIDField interface {
	UUIDExpression
	Field
}

type defaultUUIDField struct {
	uuidExpressionImpl
	fieldImpl
}

func NewUUIDField(
	table Table, name string,
) UUIDField {
	field := &defaultUUIDField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}

// TimeField

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
	field := &defaultTimeField{}
	field.expressionImpl.initFieldExpressionImpl(field)
	field.fieldImpl.initFieldImpl(field, table, name)
	return field
}
