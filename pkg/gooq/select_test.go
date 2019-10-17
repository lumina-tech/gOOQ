package gooq

import (
	"testing"
)

var selectTestCases = []TestCase{
	{
		Constructed:  Select().From(Table1),
		ExpectedStmt: "SELECT * FROM public.table1",
	},
	{
		Constructed:  SelectCount().From(Table1),
		ExpectedStmt: "SELECT COUNT(*) FROM public.table1",
	},
	{
		Constructed:  Select().Distinct().From(Table1),
		ExpectedStmt: "SELECT DISTINCT * FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column1).From(Table1),
		ExpectedStmt: "SELECT table1.column1 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column1, Table1.Column2).From(Table1),
		ExpectedStmt: "SELECT table1.column1, table1.column2 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column1.As("result")).From(Table1),
		ExpectedStmt: "SELECT table1.column1 AS result FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column1).From(Table1).Where(Table1.Column2.Eq(String("foo"))),
		ExpectedStmt: "SELECT table1.column1 FROM public.table1 WHERE table1.column2 = $1",
	},
	{
		Constructed:  Select(Table1.Column1.Filter(Table1.Column2.Eq(String("foo")))).From(Table1),
		ExpectedStmt: "SELECT table1.column1 FILTER (WHERE table1.column2 = $1) FROM public.table1",
	},
	//{
	//	Constructed: Select(Table1.Column1).From(Table1).Where(
	//		Table1.Column2.IsIn("quix", "foo"),
	//		Table1.Column2.Eq(String("quack"))),
	//	ExpectedStmt: "SELECT table1.column1 FROM public.table1 WHERE table1.column2 IN ($1, $2) AND table1.column2 = $3",
	//},
	{
		Constructed:  Select(Table1.Column3.Add(Int64(5))).From(Table1),
		ExpectedStmt: "SELECT table1.column3 + $1 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column3.Add(Float64(1.72))).From(Table1),
		ExpectedStmt: "SELECT table1.column3 + $1 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column3.Add(Table1.Column4)).From(Table1),
		ExpectedStmt: "SELECT table1.column3 + table1.column4 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column3.Sub(Table1.Column4)).From(Table1),
		ExpectedStmt: "SELECT table1.column3 - table1.column4 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column3.Mult(Table1.Column4)).From(Table1),
		ExpectedStmt: "SELECT table1.column3 * table1.column4 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column3.Div(Table1.Column4)).From(Table1),
		ExpectedStmt: "SELECT table1.column3 / table1.column4 FROM public.table1",
	},
	{
		Constructed:  Select(Table1.Column3.Div(Table1.Column4).As("result")).From(Table1),
		ExpectedStmt: "SELECT table1.column3 / table1.column4 AS result FROM public.table1",
	},
	{
		Constructed: Select(Table1.Column1).From(Table1).Where(
			Table1.Column2.Eq(String("quix")),
			Table1.Column2.Eq(String("quack"))).
			Union(
				Select(Table1.Column1).From(Table1).Where(
					Table1.Column2.Eq(String("foo")),
					Table1.Column2.Eq(String("quack")))).
			OrderBy(NewStringField(NewTable("", ""), "column2").Asc()),
		ExpectedStmt: "SELECT table1.column1 FROM public.table1 WHERE table1.column2 = $1 AND table1.column2 = $2 UNION (SELECT table1.column1 FROM public.table1 WHERE table1.column2 = $3 AND table1.column2 = $4) ORDER BY column2 ASC",
	},
	{
		Constructed:  Select().From(Table1).OrderBy(Table1.Column1.Asc()),
		ExpectedStmt: "SELECT * FROM public.table1 ORDER BY table1.column1 ASC",
	},
	{
		Constructed:  Select().From(Table1).OrderBy(Table1.Column1.Desc()),
		ExpectedStmt: "SELECT * FROM public.table1 ORDER BY table1.column1 DESC",
	},
	{
		Constructed:  Select().From(Table1).GroupBy(Table1.Column1),
		ExpectedStmt: "SELECT * FROM public.table1 GROUP BY table1.column1",
	},
	{
		Constructed:  Select(Table1.Column1.Filter(Table1.Column2.Eq(String("foo")))).From(Table1),
		ExpectedStmt: "SELECT table1.column1 FILTER (WHERE table1.column2 = $1) FROM public.table1",
	},
	{
		Constructed:  Select(Coalesce(Table1.Column1.Filter(Table1.Column2.Eq(String("foo"))), Int64(0)).As("total")).From(Table1),
		ExpectedStmt: "SELECT COALESCE(table1.column1 FILTER (WHERE table1.column2 = $1), $2) AS total FROM public.table1",
	},
	{
		Constructed:  Select().From(Table1).Limit(10),
		ExpectedStmt: "SELECT * FROM public.table1 LIMIT 10",
	},
	{
		Constructed:  Select(Table1.Column1).From(Table1).Join(Table2).On(Table2.Column1.Eq(Table1.Column1)),
		ExpectedStmt: "SELECT table1.column1 FROM public.table1 JOIN public.table2 ON table2.column1 = table1.column1",
	},
	{
		Constructed: Select(Table1.Column1, Table2.Column1).From(Table1).
			Join(Table2).On(Table2.Column1.Eq(Table1.Column1), Table2.Column2.Eq(Table1.Column2)),
		ExpectedStmt: "SELECT table1.column1, table2.column1 FROM public.table1 JOIN public.table2 ON table2.column1 = table1.column1 AND table2.column2 = table1.column2",
	},
	{
		Constructed: Select(Table1.Column1).From(Table1).
			LeftOuterJoin(Table2).On(Table2.Column1.Eq(Table1.Column1)).
			LeftOuterJoin(Table3).On(Table3.Column1.Eq(Table1.Column1)),
		ExpectedStmt: "SELECT table1.column1 FROM public.table1 LEFT OUTER JOIN public.table2 ON table2.column1 = table1.column1 LEFT OUTER JOIN public.table3 ON table3.column1 = table1.column1",
	},
	{
		Constructed: Select().From(Table1).
			LeftOuterJoin(Select(Table1.Column1).From(Table1).As("boo")).
			On(NewStringField(NewTable("", "boo"), "column1").Eq(Table1.Column1)),
		ExpectedStmt: "SELECT * FROM public.table1 LEFT OUTER JOIN (SELECT table1.column1 FROM public.table1) AS boo ON boo.column1 = table1.column1",
	},
	{
		Constructed:  Select().From(Select(Table1.Column1).From(Table1).As("boo")),
		ExpectedStmt: "SELECT * FROM (SELECT table1.column1 FROM public.table1) AS boo",
	},
	{
		Constructed: Select(Table1.Column1, Table2.Column1).From(
			Select(Table1.Column1).From(Table1).As("boo")).
			Join(Table2).On(Table2.Column1.Eq(Table1.Column1)),
		ExpectedStmt: "SELECT table1.column1, table2.column1 FROM (SELECT table1.column1 FROM public.table1) AS boo JOIN public.table2 ON table2.column1 = table1.column1",
	},
	//{
	//	Select(TimeBucket5MinutesField, Table1.Column2.Avg()).From(Table1),
	//	"SELECT time_bucket('5 minutes', table1.creation_date) AS five_min, AVG(table1.column2) FROM public.table1",
	//},
	//{
	//	Select(TimeBucket("5 minutes", Table1.CreationDate), Table1.Column1.Last(Table1.CreationDate.GetName()).As("last"), Table1.Column1.First(Table1.CreationDate.GetName()).As("first")).From(Table1),
	//	"SELECT time_bucket('5 minutes', table1.creation_date), last(table1.column1, creation_date) AS last, first(table1.column1, creation_date) AS first FROM public.table1",
	//},
	//{
	//	Select(TimeBucket5MinutesField, Table1.Column2.Avg()).From(Table1).GroupBy(TimeBucket5MinutesField),
	//	"SELECT time_bucket('5 minutes', table1.creation_date) AS five_min, AVG(table1.column2) FROM public.table1 GROUP BY five_min",
	//},
}

func TestSelects(t *testing.T) {
	runTestCases(t, selectTestCases)
}
