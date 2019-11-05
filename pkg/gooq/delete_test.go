package gooq

import "testing"

var deleteTestCases = []TestCase{
	{
		Constructed:  Delete(Table1).Where(Table1.Column1.Eq(String("foo"))),
		ExpectedStmt: `DELETE FROM public.table1 WHERE table1.column1 = $1`,
	},
	{
		Constructed:  Delete(Table1).Using(Table2).On(Table1.Column1.Eq(Table2.Column2)),
		ExpectedStmt: `DELETE FROM public.table1 USING public.table2 WHERE table1.column1 = table2.column2`,
	},
	{
		Constructed:  Delete(Table1).Using(Table2).On(Table1.Column1.Eq(Table2.Column2)).Where(Table1.Column1.Eq(String("foo"))),
		ExpectedStmt: `DELETE FROM public.table1 USING public.table2 WHERE table1.column1 = table2.column2 AND table1.column1 = $1`,
	},
	{
		Constructed:  Delete(Table1).Using(Select().From(Table2).As("foo")).On(Table1.Column1.Eq(Table2.Column2)),
		ExpectedStmt: `DELETE FROM public.table1 USING (SELECT * FROM public.table2) AS "foo" WHERE table1.column1 = table2.column2`,
	},
	{
		Constructed:  Delete(Table1).Where(Table1.Column1.Eq(String("foo"))).Returning(Table1.Column1),
		ExpectedStmt: `DELETE FROM public.table1 WHERE table1.column1 = $1 RETURNING table1.column1`,
	},
}

func TestDelete(t *testing.T) {
	runTestCases(t, deleteTestCases)
}
