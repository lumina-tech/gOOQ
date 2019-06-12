package gooq

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var predicateTypes = map[PredicateType]string{
	EqPredicate:         "=",
	GtPredicate:         ">",
	GtePredicate:        ">=",
	ILikePredicate:      "ILIKE",
	InPredicate:         "IN",
	IsNotNullPredicate:  "IS NOT NULL",
	IsNullPredicate:     "IS NULL",
	LikePredicate:       "LIKE",
	LtPredicate:         "<",
	LtePredicate:        "<=",
	NotEqPredicate:      "!=",
	NotInPredicate:      "NOT IN",
	NotILikePredicate:   "NOT ILIKE",
	NotLikePredicate:    "NOT LIKE",
	NotSimilarPredicate: "NOT SIMILAR",
	SimilarPredicate:    "SIMILAR",
}

func resolveParentAlias(alias string, col Field) string {
	if alias != "" {
		return alias
	}
	if tabCol, ok := col.(TableField); ok {
		if tabCol.Parent() != nil {
			return tabCol.Parent().MaybeAlias()
		}
	}
	return ""
}

func renderFilter(w io.Writer, d Dialect, filters []Condition, paramCount int) []interface{} {
	fmt.Fprint(w, " FILTER (")
	placeholders := renderWhereClause("", filters, d, paramCount, w)
	fmt.Fprint(w, " )")
	return placeholders
}

func renderFunction(aliased string, fun FieldFunction) string {

	if &fun == nil {
		return ""
	} else if fun.Child != nil {
		aliased = renderFunction(aliased, *fun.Child)
	}
	var f string

	// TODO(Peter): it is not good to distinguish
	// field.Count() vs Count(*) using fun.Child
	if fun.Name == "Count" && fun.Child == nil {
		f = fun.Expr
	} else if fun.Name == "CaseWhen" || fun.Name == "Coalesce" ||
		fun.Name == "DateTrunc" || fun.Name == "TimeBucket" || fun.Name == "NullIf" {
		f = fun.Expr
	} else {
		if fun.Expr == "" {
			f = aliased
		} else {
			if len(fun.Args) > 0 {
				args := make([]interface{}, 1)
				args[0] = aliased
				args = append(args, fun.Args...)
				f = fmt.Sprintf(fun.Expr, args...)
			} else {
				f = fmt.Sprintf(fun.Expr, aliased)
			}
		}
	}
	return f
}

func renderWhereClause(
	alias string, conds []Condition, d Dialect, paramCount int, w io.Writer,
) []interface{} {
	fmt.Fprint(w, " WHERE ")

	whereFragments := make([]string, len(conds))
	values := make([]interface{}, 0)

	for i, condition := range conds {
		field := condition.Binding.Field
		al := resolveParentAlias(alias, field)
		col := field.Name()
		pred := condition.Predicate
		// No arguments
		if pred == IsNullPredicate || pred == IsNotNullPredicate {
			whereFragments[i] = fmt.Sprintf("%s.%s %s", al, col, predicateTypes[pred])
		} else if pred == InPredicate || pred == NotInPredicate {
			sliceValue := reflect.ValueOf(condition.Binding.Value)
			numValues := sliceValue.Len()
			placeHolders := make([]string, numValues)
			for j := 0; j < numValues; j++ {
				paramCount = paramCount + 1
				placeHolders[j] = d.renderPlaceholder(paramCount)
				values = append(values, sliceValue.Index(j).Interface())
			}
			whereFragments[i] = fmt.Sprintf("%s.%s %s (%s)",
				al, col, predicateTypes[pred], strings.Join(placeHolders, ", "))
		} else {
			paramCount = paramCount + 1
			placeHolder := d.renderPlaceholder(paramCount)
			whereFragments[i] = fmt.Sprintf("%s.%s %s %s", al, col, predicateTypes[pred], placeHolder)
			values = append(values, condition.Binding.Value)
		}
	}
	whereClause := strings.Join(whereFragments, " AND ")
	fmt.Fprint(w, whereClause)
	return values
}

func renderJoinConditions(
	conds []JoinCondition,
) string {
	fragments := make([]string, len(conds))
	for i, cond := range conds {
		fragments[i] = renderJoinFragment(cond)
	}
	clause := strings.Join(fragments, " AND ")
	// add bracket to ON clause if len(conditions) > 1
	if len(fragments) > 1 {
		clause = fmt.Sprintf("(%s)", clause)
	}
	return clause
}

func renderFieldAlias(alias string, f TableField) (string, bool) {
	if alias != "" {
		return fmt.Sprintf("%s.%s", alias, f.Name()), true
	} else if f.Alias() != "" {
		return fmt.Sprintf("%s.%s", f.Alias(), f.Name()), true
	} else {
		return fmt.Sprintf("%s.%s", f.Parent().Alias(), f.Name()), false
	}
}

func renderJoinFragment(cond JoinCondition) string {
	lhsAlias, _ := renderFieldAlias(cond.Lhs.Parent().MaybeAlias(), cond.Lhs)
	rhsAlias, _ := renderFieldAlias(cond.Rhs.Parent().MaybeAlias(), cond.Rhs)
	return fmt.Sprintf("%s %s %s", lhsAlias, predicateTypes[cond.Predicate], rhsAlias)
}

func (d Dialect) renderPlaceholder(n int) string {
	switch d {
	case Postgres:
		return fmt.Sprintf("$%d", n)
	default:
		return "?"
	}
}

func toString(d Dialect, r Renderable) string {
	var buf bytes.Buffer
	r.Render(d, &buf)
	return buf.String()
}
