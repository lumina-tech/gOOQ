package gooq

import "testing"

var updateTestCases = []TestCase{
	{
		Constructed:  Update(Table1).Set(Table1.Column1, "10"),
		ExpectedStmt: "UPDATE public.table1 SET column1 = $1",
	},
	{
		Constructed:  Update(Table1).Set(Table1.Column3, Table1.Column4.Add(Int64(10))),
		ExpectedStmt: "UPDATE public.table1 SET column3 = table1.column4 + $1",
	},
	{
		Constructed:  Update(Table1).Set(Table1.Column3, Select().From(Table2)),
		ExpectedStmt: "UPDATE public.table1 SET column3 = (SELECT * FROM public.table2)",
	},
	{
		Constructed:  Update(Table1).Set(Table1.Column1, "10").Where(Table1.Column2.Eq(String("foo"))),
		ExpectedStmt: "UPDATE public.table1 SET column1 = $1 WHERE table1.column2 = $2",
	},
	{
		Constructed: Update(Table1).Set(Table1.Column1, Table2.Column1).
			From(Table2).Where(Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: "UPDATE public.table1 SET column1 = table2.column1 FROM public.table2 WHERE table1.column2 = table2.column2",
	},
	{
		Constructed: Update(Table1).Set(Table1.Column1, Table2.Column1).
			From(Select().From(Table2).As("foo")).Where(Table1.Column2.Eq(Table2.Column2)),
		ExpectedStmt: "UPDATE public.table1 SET column1 = table2.column1 FROM (SELECT * FROM public.table2) AS foo WHERE table1.column2 = table2.column2",
	},
	{
		Constructed:  Update(Table1).Set(Table1.Column1, "10").OnConflictDoNothing(),
		ExpectedStmt: "UPDATE public.table1 SET column1 = $1 ON CONFLICT DO NOTHING",
	},
	{
		Constructed:  Update(Table1).Set(Table1.Column1, "10").Returning(Table1.Column1),
		ExpectedStmt: "UPDATE public.table1 SET column1 = $1 RETURNING table1.column1",
	},
}

func TestUpdate(t *testing.T) {
	runTestCases(t, updateTestCases)
}
