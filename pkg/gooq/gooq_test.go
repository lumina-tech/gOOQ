package gooq

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

type testTable struct {
	TableImpl
	Column1      StringField
	Column2      StringField
	Column3      IntField
	Column4      DecimalField
	CreationDate TimeField
}

func newTestTable(name string) *testTable {
	instance := &testTable{}
	instance.TableImpl.Initialize("public", name)
	instance.Column1 = NewStringField(instance, "column1")
	instance.Column2 = NewStringField(instance, "column2")
	instance.Column3 = NewIntField(instance, "column3")
	instance.Column4 = NewDecimalField(instance, "column4")
	instance.CreationDate = NewTimeField(instance, "creation_date")
	return instance
}

func (t *testTable) As(alias string) Selectable {
	instance := newTestTable(t.name)
	instance.alias = null.StringFrom(alias)
	return instance
}

var (
	Table1           = newTestTable("table1")
	Table2           = newTestTable("table2")
	Table3           = newTestTable("table3")
	Table1Constraint = DatabaseConstraint{
		Name:    "table1_pkey",
		Columns: []Field{Table1.Column1},
	}
	//TimeBucket5MinutesField = TimeBucket("5 minutes", Table1.CreationDate).As("five_min")
)

type TestCase struct {
	Constructed  Renderable
	ExpectedStmt string
	//ExpectedPreparedStmt string
}

func runTestCases(t *testing.T, testCases []TestCase) {
	for _, rendered := range testCases {
		t.Run(rendered.ExpectedStmt, func(t *testing.T) {
			builder := Builder{}
			rendered.Constructed.Render(&builder)
			require.Equal(t, rendered.ExpectedStmt, builder.String())
		})
	}
}
