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
}

func TestFunctions(t *testing.T) {
	runTestCases(t, functionTestCases)
}
