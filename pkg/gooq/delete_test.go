package gooq

import "testing"

var deleteTestCases = []TestCase{
	{
		Delete(Table1).Where(Table1.Column1.Eq(String("foo"))),
		"DELETE FROM table1 WHERE table1.column1 = $1",
	},
	{
		Delete(Table1).Using(Table2).On(Table1.Column1.Eq(Table2.Column2)),
		"DELETE FROM table1 USING table2 table1.column1 = table2.column2 ON table1.column1 = table2.column2",
	},
	{
		Delete(Table1).Using(Select().From(Table2).As("foo")).On(Table1.Column1.Eq(Table2.Column2)),
		"DELETE FROM table1 USING (SELECT * FROM table2) AS foo table1.column1 = table2.column2 ON table1.column1 = table2.column2",
	},
	{
		Delete(Table1).Where(Table1.Column1.Eq(String("foo"))).Returning(Table1.Column1),
		"DELETE FROM table1 WHERE table1.column1 = $1 RETURNING table1.column1",
	},
}

func TestDelete(t *testing.T) {
	runTestCases(t, deleteTestCases)
}
