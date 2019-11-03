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
		Constructed:  StartsWith(String("alphabet"), String("alph")),
		ExpectedStmt: "STARTS_WITH($1, $2)",
	},
	{
		Constructed:  Chr(Table1.Column3),
		ExpectedStmt: "CHR(table1.column3)",
	},
	{
		Constructed:  Concat(String("xyzxyzabcxyz"), Table1.Column3, Int64(3)),
		ExpectedStmt: "CONCAT($1, table1.column3, $2)",
	},
	{
		Constructed:  ConcatWs(String("x"), Table1.Column3, Int64(3), String("four")),
		ExpectedStmt: "CONCAT_WS($1, table1.column3, $2, $3)",
	},
	{
		Constructed:  Format(String("Hello %s, %1$s"), Table1.Column3),
		ExpectedStmt: "FORMAT($1, table1.column3)",
	},
	{
		Constructed:  Format(String("no formatting to be done")),
		ExpectedStmt: "FORMAT($1)",
	},
}

func TestFunctions(t *testing.T) {
	runTestCases(t, functionTestCases)
}
