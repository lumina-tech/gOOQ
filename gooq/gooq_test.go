package gooq

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testTable struct {
	COLUMN1       StringField
	COLUMN2       StringField
	CREATION_DATE TimeField
	alias         string
	name          string
}

func newTestTable(name string) *testTable {
	instance := &testTable{name: name}
	instance.COLUMN1 = String(instance, "column1")
	instance.COLUMN2 = String(instance, "column2")
	instance.CREATION_DATE = Time(instance, "creation_date")
	return instance
}

func (t *testTable) IsSelectable() {}

func (t *testTable) Name() string {
	return t.name
}

func (t *testTable) String() string {
	return t.name
}

func (t *testTable) As(a string) Selectable {
	instance := newTestTable(t.name)
	t.alias = a
	return instance
}

func (t *testTable) Alias() string {
	return t.alias
}

func (t *testTable) MaybeAlias() string {
	if t.alias == "" {
		return t.name
	}
	return t.alias
}

var (
	TABLE1 = newTestTable("table1")
	TABLE2 = newTestTable("table2")
	TABLE3 = newTestTable("table3")

	TimeBucket5MinutesField = TimeBucket("5 minutes", TABLE1.CREATION_DATE).As("five_min")
)

////////////////////////////////////////////////////////////////////////////////
// Test Cases
////////////////////////////////////////////////////////////////////////////////

var rendered = []struct {
	Constructed Renderable
	Expected    string
}{
	{
		Select().From(TABLE1),
		"SELECT * FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Distinct()).From(TABLE1),
		"SELECT DISTINCT(table1.column1) FROM table1",
	},
	{
		Select(Coalesce(TABLE1.COLUMN1, 0)).From(TABLE1),
		"SELECT COALESCE(table1.column1, 0) FROM table1",
	},
	{
		Select(NullIf(TABLE1.COLUMN1, 0)).From(TABLE1),
		"SELECT NULLIF(table1.column1, 0) FROM table1",
	},
	{
		Select(Coalesce(NullIf(TABLE1.COLUMN1, 0), 0)).From(TABLE1),
		"SELECT COALESCE(NULLIF(table1.column1, 0), 0) FROM table1",
	},
	{
		Select(Coalesce(NullIf(TABLE1.COLUMN1.Sum(), 0), 1)).From(TABLE1),
		"SELECT COALESCE(NULLIF(SUM(table1.column1), 0), 1) FROM table1",
	},
	{
		Select(DateTrunc("day", TABLE1.CREATION_DATE)).From(TABLE1),
		"SELECT DATE_TRUNC('day', table1.creation_date) FROM table1",
	},
	{
		Select(TABLE1.COLUMN1, TABLE1.COLUMN2).From(TABLE1),
		"SELECT table1.column1, table1.column2 FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Div(5)).From(TABLE1),
		"SELECT table1.column1 / 5 FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Div(1.72)).From(TABLE1),
		"SELECT table1.column1 / 1.72 FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Div(TABLE1.COLUMN2)).From(TABLE1),
		"SELECT table1.column1 / column2 FROM table1",
	},
	{
		Select(Count().Cast("REAL")).From(TABLE1),
		"SELECT CAST(COUNT(*) AS REAL) FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Div(5).As("result")).From(TABLE1),
		"SELECT table1.column1 / 5 AS result FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Div(5).Sum().As("result")).From(TABLE1),
		"SELECT SUM(table1.column1 / 5) AS result FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Mult(TABLE1.COLUMN2).Sum().Div(TABLE1.COLUMN2.Sum()).As("result")).From(TABLE1),
		"SELECT SUM(table1.column1 * column2) / SUM(column2) AS result FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Div(5).Sum().Div(TABLE1.COLUMN2).As("result")).From(TABLE1),
		"SELECT SUM(table1.column1 / 5) / column2 AS result FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.Distinct().Count()).From(TABLE1),
		"SELECT COUNT(DISTINCT(table1.column1)) FROM table1",
	},
	{
		Select(Count().Cast("REAL").Div(20)).From(TABLE1),
		"SELECT CAST(COUNT(*) AS REAL) / 20 FROM table1",
	},
	{
		Select(Count().Cast("REAL").Div(20).Ceil().Cast("INT").As("calc")).From(TABLE1),
		"SELECT CAST(CEIL(CAST(COUNT(*) AS REAL) / 20) AS INT) AS calc FROM table1",
	},
	{
		Select(TABLE1.COLUMN1.As("x"), TABLE1.COLUMN2.As("y")).From(TABLE1),
		"SELECT table1.column1 AS x, table1.column2 AS y FROM table1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Eq("foo")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 = $1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Lt("foo")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 < $1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Lte("foo")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 <= $1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Gt("foo")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 > $1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Gte("foo")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 >= $1",
	},
	{
		SelectCount().From(TABLE1).Where(TABLE1.COLUMN2.Eq("foo")),
		"SELECT COUNT(*) FROM table1 WHERE table1.column2 = $1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Eq("foo")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 = $1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.IsNull()),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 IS NULL",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.IsNotNull()),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 IS NOT NULL",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.In([]string{"quix", "foo"}), TABLE1.COLUMN2.Eq("quack")),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 IN ($1, $2) AND table1.column2 = $3",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Eq("quix"), TABLE1.COLUMN2.Eq("quack")).
			Union(Select(TABLE1.COLUMN1).From(TABLE1).Where(TABLE1.COLUMN2.Eq("foo"), TABLE1.COLUMN2.Eq("quack"))).
			OrderBy(String(Table(""), "column2").Asc()),
		"SELECT table1.column1 FROM table1 WHERE table1.column2 = $1 AND table1.column2 = $2 UNION (SELECT table1.column1 FROM table1 WHERE table1.column2 = $3 AND table1.column2 = $4) ORDER BY column2 ASC",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).GroupBy(TABLE1.COLUMN1).OrderBy(TABLE1.COLUMN1.Asc()),
		"SELECT table1.column1 FROM table1 GROUP BY table1.column1 ORDER BY table1.column1 ASC",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).GroupBy(TABLE1.COLUMN1).OrderBy(TABLE1.COLUMN1.Desc()),
		"SELECT table1.column1 FROM table1 GROUP BY table1.column1 ORDER BY table1.column1 DESC",
	},
	{
		Select(TABLE1.COLUMN1, Count()).From(TABLE1).GroupBy(TABLE1.COLUMN1),
		"SELECT table1.column1, COUNT(*) FROM table1 GROUP BY table1.column1",
	},
	{
		Select(TABLE1.COLUMN1.Sum().As("total"), Count()).From(TABLE1).GroupBy(TABLE1.COLUMN1),
		"SELECT SUM(table1.column1) AS total, COUNT(*) FROM table1 GROUP BY table1.column1",
	},
	{
		Select(TABLE1.COLUMN1.Sum().FilterWhere(TABLE1.COLUMN2.Eq("quix")).As("total"), Count()).From(TABLE1).GroupBy(TABLE1.COLUMN1),
		"SELECT SUM(table1.column1) FILTER ( WHERE table1.column2 = $1 ) AS total, COUNT(*) FROM table1 GROUP BY table1.column1",
	},
	{
		Select(TABLE1.COLUMN1, Count()).From(TABLE1).Limit(10),
		"SELECT table1.column1, COUNT(*) FROM table1 LIMIT 10",
	},
	{
		Select(TABLE1.COLUMN1.Md5().Hex().Lower()).From(TABLE1),
		"SELECT LOWER(HEX(MD5(table1.column1))) FROM table1",
	},
	{
		Select(TABLE1.COLUMN1, GroupConcat(TABLE1.COLUMN2).Separator("/").As("grouped")).
			From(TABLE1).GroupBy(TABLE1.COLUMN1),
		"SELECT table1.column1, GROUP_CONCAT(table1.column2, '/') AS grouped FROM table1 GROUP BY table1.column1",
	},
	{
		Select(TABLE1.COLUMN1,
			GroupConcat(TABLE1.COLUMN2).OrderBy(TABLE1.COLUMN1).Separator("/").As("grouped")).
			From(TABLE1).GroupBy(TABLE1.COLUMN1),
		"SELECT table1.column1, GROUP_CONCAT(table1.column2 ORDER BY table1.column1 ASC SEPARATOR '/') AS grouped FROM table1 GROUP BY table1.column1",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Join(TABLE2).On(TABLE2.COLUMN1.IsEq(TABLE1.COLUMN1)),
		"SELECT table1.column1 FROM table1 JOIN table2 ON table2.column1 = table1.column1",
	},
	{
		Select(TABLE1.COLUMN1, TABLE2.COLUMN1).From(TABLE1).
			Join(TABLE2).On(TABLE2.COLUMN1.IsEq(TABLE1.COLUMN1), TABLE2.COLUMN2.IsEq(TABLE1.COLUMN2)),
		"SELECT table1.column1, table2.column1 FROM table1 JOIN table2 ON (table2.column1 = table1.column1 AND table2.column2 = table1.column2)",
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).
			LeftOuterJoin(TABLE2).On(TABLE2.COLUMN1.IsEq(TABLE1.COLUMN1)).
			LeftOuterJoin(TABLE3).On(TABLE3.COLUMN1.IsEq(TABLE1.COLUMN1)),
		"SELECT table1.column1 FROM table1 LEFT OUTER JOIN table2 ON table2.column1 = table1.column1 LEFT OUTER JOIN table3 ON table3.column1 = table1.column1",
	},
	{
		Select().From(TABLE1).
			LeftOuterJoin(Select(TABLE1.COLUMN1).From(TABLE1).As("boo")).
			On(String(Table("boo"), "column1").IsEq(TABLE1.COLUMN1)),
		"SELECT * FROM table1 LEFT OUTER JOIN (SELECT table1.column1 FROM table1) AS boo ON boo.column1 = table1.column1",
	},
	{
		Select().From(Select(TABLE1.COLUMN1).From(TABLE1).As("boo")),
		"SELECT * FROM (SELECT table1.column1 FROM table1) AS boo",
	},
	{
		Select(TABLE1.COLUMN1, TABLE2.COLUMN1).From(
			Select(TABLE1.COLUMN1).From(TABLE1).As("boo")).
			Join(TABLE2).On(TABLE2.COLUMN1.IsEq(TABLE1.COLUMN1)),
		"SELECT table1.column1, table2.column1 FROM (SELECT table1.column1 FROM table1) AS boo JOIN table2 ON table2.column1 = table1.column1",
	},
	{
		InsertInto(TABLE1).Set(TABLE1.COLUMN1, "foo"),
		"INSERT INTO table1 (column1) VALUES ($1)",
	},
	{
		InsertInto(TABLE1).Select(Select(TABLE1.COLUMN1).From(TABLE1)),
		"INSERT INTO table1 (SELECT table1.column1 FROM table1)",
	},
	{
		InsertInto(TABLE1).Set(TABLE1.COLUMN1, "foo").OnConflictDoNothing(),
		"INSERT INTO table1 (column1) VALUES ($1) ON CONFLICT DO NOTHING",
	},
	{
		InsertInto(TABLE1).Set(TABLE1.COLUMN1, "foo").Returning(TABLE1.COLUMN1),
		"INSERT INTO table1 (column1) VALUES ($1) RETURNING column1",
	},
	{
		Update(TABLE1).Set(TABLE1.COLUMN1, "foo").Where(TABLE1.COLUMN2.Eq("table2")),
		"UPDATE table1 SET column1 = $1 WHERE table1.column2 = $2",
	},
	{
		Update(TABLE1).Set(TABLE1.COLUMN1, "foo").Where(TABLE1.COLUMN2.Eq("table2")).OnConflict(ConflictActionDoNothing),
		"UPDATE table1 SET column1 = $1 WHERE table1.column2 = $2 ON CONFLICT DO NOTHING",
	},
	{
		Update(TABLE1).Set(TABLE1.COLUMN1, "foo").Join(TABLE2).On(TABLE1.COLUMN1.IsEq(TABLE2.COLUMN1)),
		"UPDATE table1 SET column1 = $1 FROM table2 WHERE table1.column1 = table2.column1",
	},
	{
		Update(TABLE1).Set(TABLE1.COLUMN1, TABLE2.COLUMN1).
			Join(TABLE2).On(TABLE1.COLUMN1.IsEq(TABLE2.COLUMN1)),
		"UPDATE table1 SET column1 = table2.column1 FROM table2 WHERE table1.column1 = table2.column1",
	},
	{
		Update(TABLE1).Set(TABLE1.COLUMN1, TABLE2.COLUMN1).
			Join(Select(TABLE1.COLUMN1).From(TABLE1).As("column1")).On(TABLE1.COLUMN1.IsEq(TABLE2.COLUMN1)),
		"UPDATE table1 SET column1 = table2.column1 FROM (SELECT table1.column1 FROM table1) AS column1 WHERE table1.column1 = table2.column1",
	},
	{
		Delete(TABLE1).Where(TABLE1.COLUMN2.Eq("table2")),
		"DELETE FROM table1 WHERE table1.column2 = $1",
	},
	{
		Delete(TABLE1).Using(Select(TABLE1.COLUMN1).From(TABLE1).As("column2")).On(TABLE1.COLUMN1.IsEq(TABLE1.COLUMN2)),
		"DELETE FROM table1 USING (SELECT table1.column1 FROM table1) AS column2 WHERE table1.column1 = table1.column2",
	},
	{
		Select(TimeBucket5MinutesField, TABLE1.COLUMN2.Avg()).From(TABLE1),
		"SELECT time_bucket('5 minutes', table1.creation_date) AS five_min, AVG(table1.column2) FROM table1",
	},
	{
		Select(TimeBucket("5 minutes", TABLE1.CREATION_DATE), TABLE1.COLUMN1.Last(TABLE1.CREATION_DATE.Name()).As("last"), TABLE1.COLUMN1.First(TABLE1.CREATION_DATE.Name()).As("first")).From(TABLE1),
		"SELECT time_bucket('5 minutes', table1.creation_date), last(table1.column1, creation_date) AS last, first(table1.column1, creation_date) AS first FROM table1",
	},
	{
		Select(TimeBucket5MinutesField, TABLE1.COLUMN2.Avg()).From(TABLE1).GroupBy(TimeBucket5MinutesField),
		"SELECT time_bucket('5 minutes', table1.creation_date) AS five_min, AVG(table1.column2) FROM table1 GROUP BY five_min",
	},
}

var selectTrees = []struct {
	Constructed Selectable
	Expected    selection
}{
	{
		Select().From(TABLE1),
		selection{
			selection: TABLE1,
		},
	},
	{
		Select(TABLE1.COLUMN1, TABLE1.COLUMN2).From(TABLE1),
		selection{
			selection:  TABLE1,
			projection: []Field{TABLE1.COLUMN1, TABLE1.COLUMN2},
		},
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).Join(TABLE2).On(TABLE2.COLUMN1.IsEq(TABLE1.COLUMN1)),
		selection{
			selection:  TABLE1,
			projection: []Field{TABLE1.COLUMN1},
			joinTarget: nil,
			joinType:   NotJoined,
			joins: []join{
				join{
					target:   TABLE2,
					joinType: Join,
					conds:    []JoinCondition{TABLE2.COLUMN1.IsEq(TABLE1.COLUMN1)},
				},
			},
		},
	},
	{
		Select(TABLE1.COLUMN1).From(TABLE1).GroupBy(TABLE1.COLUMN1).OrderBy(TABLE1.COLUMN1),
		selection{
			selection:  TABLE1,
			projection: []Field{TABLE1.COLUMN1},
			groups:     []Field{TABLE1.COLUMN1},
			ordering:   []Field{TABLE1.COLUMN1},
		},
	},
	{
		Select().From(Select(TABLE1.COLUMN1).From(TABLE1)),
		selection{
			selection: &selection{
				selection:  TABLE1,
				projection: []Field{TABLE1.COLUMN1},
			},
		},
	},
}

var insertTrees = []struct {
	Constructed InsertSetMoreStep
	Expected    insert
}{
	{
		InsertInto(TABLE1).Set(TABLE1.COLUMN1, "foo"),
		insert{
			table: TABLE1,
			bindings: []TableFieldBinding{
				TableFieldBinding{
					Field: TABLE1.COLUMN1,
					Value: "foo",
				},
			},
		},
	},
}

var updateTrees = []struct {
	Constructed Executable
	Expected    update
}{
	{
		Update(TABLE1).Set(TABLE1.COLUMN1, "foo").Where(TABLE1.COLUMN2.Eq("table2")),
		update{
			table: TABLE1,
			bindings: []TableFieldBinding{
				TableFieldBinding{
					Field: TABLE1.COLUMN1,
					Value: "foo",
				},
			},
			predicate: []Condition{
				Condition{
					Binding: FieldBinding{
						Field: TABLE1.COLUMN2,
						Value: "table2",
					},
					Predicate: EqPredicate,
				},
			},
		},
	},
}

var deleteTrees = []struct {
	Constructed Executable
	Expected    deletion
}{
	{
		Delete(TABLE1).Where(TABLE1.COLUMN2.Eq("table2")),
		deletion{
			table: TABLE1,
			predicate: []Condition{
				Condition{
					Binding: FieldBinding{
						Field: TABLE1.COLUMN2,
						Value: "table2",
					},
					Predicate: EqPredicate,
				},
			},
		},
	},
}

func TestDeleteTrees(t *testing.T) {
	for _, tree := range deleteTrees {
		require.Equal(t, &tree.Expected, tree.Constructed)
	}
}

func TestUpdateTrees(t *testing.T) {
	for _, tree := range updateTrees {
		require.Equal(t, &tree.Expected, tree.Constructed)
	}
}

func TestInsertTrees(t *testing.T) {
	for _, tree := range insertTrees {
		require.Equal(t, &tree.Expected, tree.Constructed)
	}
}

func TestSelectTrees(t *testing.T) {
	for _, tree := range selectTrees {
		require.Equal(t, &tree.Expected, tree.Constructed)
	}
}

func TestRendered(t *testing.T) {
	for _, rendered := range rendered {
		t.Run(rendered.Expected, func(t *testing.T) {
			r := rendered.Constructed.String(Postgres)
			require.Equal(t, rendered.Expected, r)
		})
	}
}

func TestRenderedSingle(t *testing.T) {
	constructed := Select(TimeBucket5MinutesField, TABLE1.COLUMN2.Avg()).From(TABLE1).GroupBy(TimeBucket5MinutesField)
	r := constructed.String(Postgres)
	fmt.Println(r)
}
