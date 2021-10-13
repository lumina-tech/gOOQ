package gooq

import "testing"

var insertTestCases = []TestCase{
	{
		Constructed:  InsertInto(Table1).Set(Table1.Column1, "foo"),
		ExpectedStmt: `INSERT INTO public.table1 (column1) VALUES ($1)`,
	},
	{
		Constructed:  InsertInto(Table1).Set(Table1.Column1, "foo").Set(Table1.Column2, "bar"),
		ExpectedStmt: `INSERT INTO public.table1 (column1, column2) VALUES ($1, $2)`,
	},
	{
		Constructed: InsertInto(Table1).
			Values("1", "2", 3, 4).
			Values("2", "3", 4, 5),
		ExpectedStmt: `INSERT INTO public.table1 VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`,
	},
	{
		Constructed: InsertInto(Table1).
			Columns(Table1.Column1, Table1.Column2, Table1.Column3, Table1.Column4).
			Values("1", "2", 3, 4).
			Values("2", "3", 4, 5),
		ExpectedStmt: `INSERT INTO public.table1 (column1, column2, column3, column4) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)`,
	},
	{
		Constructed:  InsertInto(Table1).Select(Select(Table1.Column1).From(Table1)),
		ExpectedStmt: `INSERT INTO public.table1 (SELECT "table1".column1 FROM public.table1)`,
	},
	{
		Constructed:  InsertInto(Table1).Set(Table1.Column1, "foo").Returning(Table1.Column1),
		ExpectedStmt: `INSERT INTO public.table1 (column1) VALUES ($1) RETURNING "table1".column1`,
	},
	{
		Constructed: InsertInto(Table1).
			Set(Table1.Column1, "foo").Set(Table1.Column2, "bar").
			OnConflictDoUpdate(&Table1Constraint).
			SetUpdates(Table1.Column2, String("bar")),
		ExpectedStmt: `INSERT INTO public.table1 (column1, column2) VALUES ($1, $2) ON CONFLICT ON CONSTRAINT table1_pkey DO UPDATE SET column2 = $3`,
	},
	{
		Constructed: InsertInto(Table1).
			Set(Table1.Column1, "foo").Set(Table1.Column2, "bar").
			OnConflictDoUpdate(&Table1Constraint).
			SetUpdateColumns(Table1.Column2),
		ExpectedStmt: `INSERT INTO public.table1 (column1, column2) VALUES ($1, $2) ON CONFLICT ON CONSTRAINT table1_pkey DO UPDATE SET column2 = "excluded".column2`,
	},
	{
		Constructed:  InsertInto(Table1).Set(Table1.Column1, "foo").OnConflictDoNothing(),
		ExpectedStmt: `INSERT INTO public.table1 (column1) VALUES ($1) ON CONFLICT DO NOTHING`,
	},
	{
		Constructed: InsertInto(Table2).
			Set(Table2.Column1, "foo").
			Set(Table2.Column2, "bar").
			Set(Table2.Column3, 1).
			OnConflictDoUpdate(&Table2Constraint).
			SetUpdateColumns(Table2.Column3),
		ExpectedStmt: `INSERT INTO public.table2 (column1, column2, column3) VALUES ($1, $2, $3) ON CONFLICT ( column1, column2 ) WHERE ((bool_column)::bool <> 'true'::bool) DO UPDATE SET column3 = "excluded".column3`,
	},
}

func TestInsert(t *testing.T) {
	runTestCases(t, insertTestCases)
}
