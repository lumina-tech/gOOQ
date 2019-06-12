package gooq

import (
	"fmt"
	"strings"
	"time"
)

func Count() IntField {
	return &intField{name: "*", fun: FieldFunction{Name: "Count", Expr: "COUNT(*)"}}
}

func DateTrunc(format string, field interface{}) TimeField {
	parsedField := ""
	switch value := field.(type) {
	case Field:
		table := resolveParentAlias("", value)
		parsedField = fmt.Sprintf("%s.%s", table, value.Name())
	default:
		parsedField = fmt.Sprintf("%v", value)
	}
	return &timeField{
		name: "DateTrunc",
		fun: FieldFunction{
			Name: "DateTrunc",
			Expr: fmt.Sprintf("DATE_TRUNC('%s', %s)", format, parsedField),
		},
	}
}

type groupConcat struct {
	stringField
}

// This may indicate that the rendering pipeline needs to get adjusted so that things like can be less stateful
func (g *groupConcat) OrderBy(f Field) *groupConcat {
	al := resolveParentAlias(f.Alias(), f)
	g.stringField.fun.Expr = "GROUP_CONCAT(%s ORDER BY %s.%s ASC)"
	g.stringField.fun.Args = append(g.stringField.fun.Args, al, f.Name())
	return g
}

func (g *groupConcat) Separator(s string) *groupConcat {
	if len(g.stringField.fun.Args) > 0 {
		g.stringField.fun.Expr = "GROUP_CONCAT(%s ORDER BY %s.%s ASC SEPARATOR '%s')" // TODO ASC is hard coded
	} else {
		g.stringField.fun.Expr = "GROUP_CONCAT(%s, '%s')" // TODO sqlite specific (i.e. no SEPARATOR keyword)
	}

	g.stringField.fun.Args = append(g.stringField.fun.Args, s)
	return g
}

func CaseWhen(expr string, fields ...Field) StringField {
	args := make([]interface{}, 0)
	for _, field := range fields {
		table := resolveParentAlias("", field)
		qualifiedField := fmt.Sprintf("%s.%s", table, field.Name())
		args = append(args, qualifiedField)
	}

	return &stringField{
		name: "CaseWhen",
		fun: FieldFunction{
			Name: "CaseWhen",
			Expr: fmt.Sprintf(expr, args...),
		},
	}
}

func Coalesce(args ...interface{}) StringField {
	strs := []string{}
	for _, arg := range args {
		switch value := arg.(type) {
		case Field:
			table := resolveParentAlias("", value)
			qualifiedField := value.String()
			if table != "" {
				qualifiedField = fmt.Sprintf("%s.%s", table, value.Name())
			}
			strs = append(strs, qualifiedField)
		default:
			strs = append(strs, fmt.Sprintf("%v", value))
		}
	}
	return &stringField{
		name: "Coalesce",
		fun: FieldFunction{
			Name: "Coalesce",
			Expr: fmt.Sprintf("COALESCE(%s)", strings.Join(strs, ", ")),
		},
	}
}

func NullIf(field Field, value interface{}) StringField {
	table := resolveParentAlias("", field)
	qualifiedField := field.String()
	if table != "" {
		qualifiedField = fmt.Sprintf("%s.%s", table, field.Name())
		qualifiedField = renderFunction(qualifiedField, field.Function())
	}
	return &stringField{
		name: "NullIf",
		fun: FieldFunction{
			Name: "NullIf",
			Expr: fmt.Sprintf("NULLIF(%s, %v)", qualifiedField, value),
		},
	}
}

func GroupConcat(field Field) *groupConcat {
	var s Selectable
	if tf, ok := field.(TableField); ok {
		s = tf.Parent()
	}
	return &groupConcat{
		stringField: stringField{
			name:      field.Name(),
			selection: s,
			fun: FieldFunction{
				Name: "GroupConcat",
				Expr: "GROUP_CONCAT(%s)",
			},
		},
	}
}

// timescaledb stuff
func TimeBucket(bucket_width string, field Field) TimeField {
	table := resolveParentAlias("", field)
	parsedField := fmt.Sprintf("%s.%s", table, field.Name())
	return &timeField{
		name: "TimeBucket",
		fun: FieldFunction{
			Name: "TimeBucket",
			Expr: fmt.Sprintf("time_bucket('%s', %s)", bucket_width, parsedField),
		},
	}
}

func TimeBucketGapFill(
	bucket_width string, field Field, start, end time.Time,
) TimeField {
	table := resolveParentAlias("", field)
	parsedField := fmt.Sprintf("%s.%s", table, field.Name())
	return &timeField{
		name: "TimeBucket",
		fun: FieldFunction{
			Name: "TimeBucket",
			Expr: fmt.Sprintf("time_bucket('%s', %s, %s, %s)", bucket_width, parsedField, start.String(), end.String()),
		},
	}
}
