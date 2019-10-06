package gooq

import "testing"

var expressionTestCases = []TestCase{
	{
		Table1.Column1.Eq(Table2.Column1).Or(Table1.Column2.Eq(Table2.Column2)),
		"table1.column1 = table2.column1 Or table1.column2 = table2.column2",
	},
}

func TestExpressions(t *testing.T) {
	runTestCases(t, expressionTestCases)
}
