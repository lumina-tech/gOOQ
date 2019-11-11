package gooq

import "testing"

var functionTestCases = []TestCase{
	{
		Constructed:  Select(Coalesce(Table1.Column1, Table1.Column2)).From(Table1),
		ExpectedStmt: "SELECT COALESCE(table1.column1, table1.column2) FROM public.table1",
	},
	{
		Constructed:  Select(Coalesce(Table1.Column1, Int64(0))).From(Table1),
		ExpectedStmt: "SELECT COALESCE(table1.column1, $1) FROM public.table1",
	},
	{
		Constructed:  Select(Count(Asterisk)).From(Table1),
		ExpectedStmt: "SELECT COUNT(*) FROM public.table1",
	},
	{
		Constructed:  Greatest(Int64(10), Int64(2), Int64(23)),
		ExpectedStmt: "GREATEST($1, $2, $3)",
	},
	{
		Constructed:  Least(String("a"), String("b")),
		ExpectedStmt: "LEAST($1, $2)",
	},
	{
		Constructed:  Ascii(String("abc")),
		ExpectedStmt: "ASCII($1)",
	},
	{
		Constructed:  Ascii(Table1.Column1),
		ExpectedStmt: "ASCII(table1.column1)",
	},
	{
		Constructed:  BTrim(String("    abc    ")),
		ExpectedStmt: "BTRIM($1)",
	},
	{
		Constructed:  LTrim(Table1.Column1, String("xyz")),
		ExpectedStmt: "LTRIM(table1.column1, $1)",
	},
	{
		Constructed:  RTrim(String("xyzxyzabcxyz"), Table1.Column1),
		ExpectedStmt: "RTRIM($1, table1.column1)",
	},
	{
		Constructed:  LessThan(Int64(10), Int64(2)),
		ExpectedStmt: "$1 < $2",
	},
	{
		Constructed:  GreaterThan(Int64(10), Int64(2)),
		ExpectedStmt: "$1 > $2",
	},
	{
		Constructed:  LessOrEqual(Int64(10), Int64(2)),
		ExpectedStmt: "$1 <= $2",
	},
	{
		Constructed:  GreaterOrEqual(Int64(10), Int64(2)),
		ExpectedStmt: "$1 >= $2",
	},
	{
		Constructed:  Equal(Int64(10), Int64(2)),
		ExpectedStmt: "$1 = $2",
	},
	{
		Constructed:  NotEqual(Int64(10), Int64(2)),
		ExpectedStmt: "$1 <> $2",
	},
}

func TestFunctions(t *testing.T) {
	runTestCases(t, functionTestCases)
}
