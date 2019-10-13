package gooq

import "testing"

var insertTestCases = []TestCase{
	{
		InsertInto(Table1).Set(Table1.Column1, "foo"),
		"INSERT INTO public.table1 (column1) VALUES ($1)",
	},
	{
		InsertInto(Table1).Set(Table1.Column1, "foo").Set(Table1.Column2, "bar"),
		"INSERT INTO public.table1 (column1, column2) VALUES ($1, $2)",
	},
	{
		InsertInto(Table1).
			Values("1", "2", 3, 4).
			Values("2", "3", 4, 5),
		"INSERT INTO public.table1 VALUES ($1, $2, $3, $4) ($5, $6, $7, $8)",
	},
	{
		InsertInto(Table1).
			Columns(Table1.Column1, Table1.Column2, Table1.Column3, Table1.Column4).
			Values("1", "2", 3, 4).
			Values("2", "3", 4, 5),
		"INSERT INTO public.table1 (column1, column2, column3, column4) VALUES ($1, $2, $3, $4) ($5, $6, $7, $8)",
	},
	{
		InsertInto(Table1).Select(Select(Table1.Column1).From(Table1)),
		"INSERT INTO public.table1 (SELECT table1.column1 FROM public.table1)",
	},
	{
		InsertInto(Table1).Set(Table1.Column1, "foo").OnConflictDoNothing(),
		"INSERT INTO public.table1 (column1) VALUES ($1) ON CONFLICT DO NOTHING",
	},
	{
		InsertInto(Table1).Set(Table1.Column1, "foo").Returning(Table1.Column1),
		"INSERT INTO public.table1 (column1) VALUES ($1) RETURNING table1.column1",
	},
	{
		InsertInto(Table1).
			Set(Table1.Column1, "foo").Set(Table1.Column2, "bar").
			OnConflictDoUpdate(&Table1Constraint).
			SetUpdates(Table1.Column2, String("bar")),
		"INSERT INTO public.table1 (column1, column2) VALUES ($1, $2) ON CONFLICT DO UPDATE (table1.column1) SET table1.column2 = $3",
	},
	{
		InsertInto(Table1).
			Set(Table1.Column1, "foo").Set(Table1.Column2, "bar").
			OnConflictDoUpdate(&Table1Constraint).
			SetUpdateColumns(Table1.Column2),
		"INSERT INTO public.table1 (column1, column2) VALUES ($1, $2) ON CONFLICT DO UPDATE (table1.column1) SET table1.column2 = EXCLUDED.column2",
	},
}

func TestInsert(t *testing.T) {
	runTestCases(t, insertTestCases)
}
