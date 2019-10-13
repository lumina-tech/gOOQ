package gooq

import "testing"

var updateTestCases = []TestCase{
	{
		Update(Table1).Set(Table1.Column1, "10"),
		"UPDATE public.table1 SET table1.column1 = $1",
	},
	{
		Update(Table1).Set(Table1.Column3, Table1.Column4.Add(Int64(10))),
		"UPDATE public.table1 SET table1.column3 = table1.column4 + $1",
	},
	{
		Update(Table1).Set(Table1.Column3, Select().From(Table2)),
		"UPDATE public.table1 SET table1.column3 = (SELECT * FROM public.table2)",
	},
	{
		Update(Table1).Set(Table1.Column1, "10").Where(Table1.Column2.Eq(String("foo"))),
		"UPDATE public.table1 SET table1.column1 = $1 WHERE table1.column2 = $2",
	},
	{
		Update(Table1).Set(Table1.Column1, Table2.Column1).
			From(Table2).On(Table1.Column2.Eq(Table2.Column2)),
		"UPDATE public.table1 SET table1.column1 = table2.column1 FROM public.table2 WHERE table1.column2 = table2.column2",
	},
	{
		Update(Table1).Set(Table1.Column1, Table2.Column1).
			From(Select().From(Table2)).On(Table1.Column2.Eq(Table2.Column2)),
		"UPDATE public.table1 SET table1.column1 = table2.column1 FROM (SELECT * FROM public.table2) WHERE table1.column2 = table2.column2",
	},
	{
		Update(Table1).Set(Table1.Column1, "10").OnConflictDoNothing(),
		"UPDATE public.table1 SET table1.column1 = $1 ON CONFLICT DO NOTHING",
	},
	{
		Update(Table1).Set(Table1.Column1, "10").Returning(Table1.Column1),
		"UPDATE public.table1 SET table1.column1 = $1 RETURNING table1.column1",
	},
}

func TestUpdate(t *testing.T) {
	runTestCases(t, updateTestCases)
}
