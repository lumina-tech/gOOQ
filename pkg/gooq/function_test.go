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
}

func TestFunctions(t *testing.T) {
	runTestCases(t, functionTestCases)
}
