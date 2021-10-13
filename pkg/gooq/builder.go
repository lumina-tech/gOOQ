package gooq

import (
	"bytes"
	"fmt"
)

var (
	invalidLiteralValueError = fmt.Errorf("literal value cannot be of kind slice")
)

type Builder struct {
	isDebug   bool
	buffer    bytes.Buffer
	arguments []interface{}
	errors    []error
}

func (builder *Builder) Printf(
	str string, args ...interface{},
) *Builder {
	return builder.Print(fmt.Sprintf(str, args...))
}

func (builder *Builder) Print(
	str string,
) *Builder {
	builder.buffer.Write([]byte(str))
	return builder
}

func (builder *Builder) RenderExpression(
	expression Expression,
) *Builder {
	expression.Render(builder)
	return builder
}

func (builder *Builder) RenderLiteral(
	value interface{},
) {
	if builder.isDebug {
		builder.Print(fmt.Sprintf("%v", value))
	} else {
		placeholder := fmt.Sprintf("$%d", len(builder.arguments)+1)
		builder.Print(placeholder)
		builder.arguments = append(builder.arguments, value)
	}
}

func (builder *Builder) RenderExpressionArray(
	array []Expression,
) {
	builder.Print("(")
	for index, expression := range array {
		builder.RenderExpression(expression)
		if index != len(array)-1 {
			builder.Print(", ")
		}
	}
	builder.Print(")")
}

func (builder *Builder) RenderFieldArray(
	fields []Field,
) {
	builder.Print("(")
	for index, field := range fields {
		builder.Print(field.GetName())
		if index != len(fields)-1 {
			builder.Print(", ")
		}
	}
	builder.Print(")")
}

func (builder *Builder) RenderConditions(
	conditions []Expression,
) {
	for index, expression := range conditions {
		expression.Render(builder)
		if index != len(conditions)-1 {
			builder.Print(" AND ")
		}
	}
}

func (builder *Builder) RenderExpressions(
	expressions []Expression,
) {
	for index, expression := range expressions {
		expression.Render(builder)
		if index != len(expressions)-1 {
			builder.Print(", ")
		}
	}
}

func (builder *Builder) RenderProjections(
	projections []Selectable,
) {
	for index, expression := range projections {
		expression.Render(builder)
		if index != len(projections)-1 {
			builder.Print(", ")
		}
	}
}

func (builder *Builder) RenderSetPredicates(
	predicates []setPredicate,
) *Builder {
	for index := range predicates {
		item := &predicates[index]
		// https://www.postgresql.org/docs/12/sql-update.html
		// do not include the table's name in the specification of a target column — for example, UPDATE table_name SET table_name.col = 1 is invalid
		builder.Printf("%s = ", item.field.GetName())
		switch predicate := item.value.(type) {
		case *selection:
			builder.Print("(")
			predicate.Render(builder)
			builder.Printf(")")
		case Expression:
			builder.RenderExpression(predicate)
		default:
			builder.RenderExpression(newLiteralExpression(item.value))
		}
		if index != len(predicates)-1 {
			builder.Print(", ")
		}
	}
	return builder
}

func (builder *Builder) String() string {
	return builder.buffer.String()
}
