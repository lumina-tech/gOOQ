package gooq

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

type testTable struct {
	TableImpl
	ID            UUIDField
	Column1       StringField
	Column2       StringField
	Column3       IntField
	Column4       DecimalField
	BoolColumn    BoolField
	DecimalColumn DecimalField
	StringColumn  StringField
	TimeColumn    TimeField
}

func newTestTable(name string) *testTable {
	instance := &testTable{}
	instance.TableImpl.Initialize("public", name)
	instance.ID = NewUUIDField(instance, "id")
	instance.Column1 = NewStringField(instance, "column1")
	instance.Column2 = NewStringField(instance, "column2")
	instance.Column3 = NewIntField(instance, "column3")
	instance.Column4 = NewDecimalField(instance, "column4")
	instance.BoolColumn = NewBoolField(instance, "bool_column")
	instance.DecimalColumn = NewDecimalField(instance, "decimal_column")
	instance.StringColumn = NewStringField(instance, "string_column")
	instance.TimeColumn = NewTimeField(instance, "time_column")
	return instance
}

func (t *testTable) As(alias string) Table {
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
	Arguments    interface{}
	Errors       []error
	//ExpectedPreparedStmt string
}

func runTestCases(t *testing.T, testCases []TestCase) {
	for _, rendered := range testCases {
		t.Run(rendered.ExpectedStmt, func(t *testing.T) {
			builder := Builder{}
			rendered.Constructed.Render(&builder)
			require.Equal(t, rendered.ExpectedStmt, builder.String())
			if rendered.Errors != nil {
				require.Equal(t, rendered.Errors, builder.errors)
			}
			if rendered.Arguments != nil {
				require.Equal(t, rendered.Arguments, builder.arguments)
			}
		})
	}
}
