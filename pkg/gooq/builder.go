package gooq

import (
	"bytes"
	"fmt"
)

type Builder struct {
	isDebug   bool
	buffer    bytes.Buffer
	arguments []interface{}
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

func (builder *Builder) RenderLiteralArray(
	array []interface{},
) {
	builder.Print("(")
	for index, item := range array {
		builder.RenderLiteral(item)
		if index != len(array)-1 {
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

func (builder *Builder) RenderFields(
	projections []Expression,
) {
	for index, expression := range projections {
		expression.Render(builder)
		if index != len(projections)-1 {
			builder.Print(", ")
		}
	}
}

func (builder *Builder) RenderProjections(
	projections []Expression,
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
		builder.Printf("%s = ", item.field.QualifiedName())
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
