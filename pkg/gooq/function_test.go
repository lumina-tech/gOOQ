package gooq

import "testing"

var functionTestCases = []TestCase{
	{
		Select(Coalesce(Table1.Column1, Table1.Column2)).From(Table1),
		"SELECT COALESCE(table1.column1, table1.column2) FROM public.table1",
	},
	{
		Select(Coalesce(Table1.Column1, Int64(0))).From(Table1),
		"SELECT COALESCE(table1.column1, $1) FROM public.table1",
	},
	{
		Select(Count(Asterisk)).From(Table1),
		"SELECT COUNT(*) FROM public.table1",
	},
}

func TestFunctions(t *testing.T) {
	runTestCases(t, functionTestCases)
}
