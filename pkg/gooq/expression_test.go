package gooq

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var expressionTestCases = []TestCase{
	{
		Constructed:  Count(Table1.Column1),
		ExpectedStmt: `COUNT("table1".column1)`,
	},
	{
		Constructed:  Count(Table1.Column1).IsGt(5),
		ExpectedStmt: `COUNT("table1".column1) > $1`,
		Arguments:    []interface{}{float64(5)},
	},
	{
		Constructed:  Table1.Column1.Asc(),
		ExpectedStmt: `"table1".column1 ASC`,
	},
	{
		Constructed:  Table1.Column1.Desc(),
		ExpectedStmt: `"table1".column1 DESC`,
	},
	{
		Constructed:  Table1.Column1.IsNull(),
		ExpectedStmt: `"table1".column1 IS NULL`,
	},
	{
		Constructed:  Table1.Column1.IsNotNull(),
		ExpectedStmt: `"table1".column1 IS NOT NULL`,
	},
	{
		Constructed:  Table1.Column1.Eq(Table2.Column1).Or(Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: `("table1".column1 = "table2".column1 OR "table1".column2 = "table2".column2)`,
	},
	{
		Constructed:  Table1.Column1.Eq(Table2.Column1).And(Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: `("table1".column1 = "table2".column1 AND "table1".column2 = "table2".column2)`,
	},
	{
		Constructed:  Table1.Column1.Eq(Table2.Column1).And(Table1.Column2.Eq(Table2.Column2)).And(Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: `(("table1".column1 = "table2".column1 AND "table1".column2 = "table2".column2) AND "table1".column2 = "table2".column2)`,
	},
	{
		Constructed:  Table1.Column1.Eq(Table2.Column1).And(Table1.Column2.Eq(Table2.Column2)).Or(Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: `(("table1".column1 = "table2".column1 AND "table1".column2 = "table2".column2) OR "table1".column2 = "table2".column2)`,
	},
	{
		Constructed:  Table1.BoolColumn.IsEq(true),
		ExpectedStmt: `"table1".bool_column = $1`,
	},
	{
		Constructed:  Table1.BoolColumn.IsNotEq(true),
		ExpectedStmt: `"table1".bool_column != $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.IsGt(1.0),
		ExpectedStmt: `"table1".decimal_column > $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.IsGte(1.0),
		ExpectedStmt: `"table1".decimal_column >= $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.IsLt(1.0),
		ExpectedStmt: `"table1".decimal_column < $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.IsLte(1.0),
		ExpectedStmt: `"table1".decimal_column <= $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.IsEq(1.0),
		ExpectedStmt: `"table1".decimal_column = $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.IsNotEq(1.0),
		ExpectedStmt: `"table1".decimal_column != $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.Add(Int64(42)),
		ExpectedStmt: `"table1".decimal_column + $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.Sub(Int64(-42)),
		ExpectedStmt: `"table1".decimal_column - $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.Mult(Int64(42)),
		ExpectedStmt: `"table1".decimal_column * $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.Div(Int64(0)),
		ExpectedStmt: `"table1".decimal_column / $1`,
	},
	{
		Constructed:  Table1.DecimalColumn.Sqrt(),
		ExpectedStmt: `|/ "table1".decimal_column`,
	},
	{
		Constructed:  Table1.StringColumn.IsEq("foo"),
		ExpectedStmt: `"table1".string_column = $1`,
	},
	{
		Constructed:  Table1.StringColumn.IsNotEq("foo"),
		ExpectedStmt: `"table1".string_column != $1`,
	},
	{
		Constructed:  Table1.TimeColumn.IsGt(time.Now()),
		ExpectedStmt: `"table1".time_column > $1`,
	},
	{
		Constructed:  Table1.TimeColumn.IsGte(time.Now()),
		ExpectedStmt: `"table1".time_column >= $1`,
	},
	{
		Constructed:  Table1.TimeColumn.IsLt(time.Now()),
		ExpectedStmt: `"table1".time_column < $1`,
	},
	{
		Constructed:  Table1.TimeColumn.IsLte(time.Now()),
		ExpectedStmt: `"table1".time_column <= $1`,
	},
	{
		Constructed:  Table1.TimeColumn.IsEq(time.Now()),
		ExpectedStmt: `"table1".time_column = $1`,
	},
	{
		Constructed:  Table1.TimeColumn.IsNotEq(time.Now()),
		ExpectedStmt: `"table1".time_column != $1`,
	},
	{
		Constructed:  Table1.ID.IsEq(uuid.Nil),
		ExpectedStmt: `"table1".id = $1`,
	},
	{
		Constructed:  Table1.ID.IsNotEq(uuid.Nil),
		ExpectedStmt: `"table1".id != $1`,
	},
	{
		Constructed:  Table1.ID.In(Select().From(Table1)),
		ExpectedStmt: `"table1".id IN (SELECT * FROM public.table1)`,
	},
	{
		Constructed:  Table1.ID.NotIn(Select().From(Table1)),
		ExpectedStmt: `"table1".id NOT IN (SELECT * FROM public.table1)`,
	},
	{
		Constructed:  Table1.ID.IsIn(uuid.Nil, uuid.Nil),
		ExpectedStmt: `"table1".id IN ($1, $2)`,
		Arguments:    []interface{}{uuid.Nil, uuid.Nil},
	},
	{
		Constructed:  Table1.ID.IsNotIn(uuid.Nil, uuid.Nil),
		ExpectedStmt: `"table1".id NOT IN ($1, $2)`,
		Arguments:    []interface{}{uuid.Nil, uuid.Nil},
	},
}

func TestExpressions(t *testing.T) {
	runTestCases(t, expressionTestCases)
}
