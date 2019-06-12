package gooq

func keyword(value string) Expression {
	return newKeywordExpression(value)
}

func Literal(value interface{}) Expression {
	return newLiteralExpression(value)
}

func LiteralArray(value []interface{}) Expression {
	return newLiteralArrayExpression(value)
}

func String(value string) StringExpression {
	expr := &stringExpressionImpl{}
	expr.expressionImpl.initLiteralExpression(expr, value)
	return expr
}

func Int64(value int64) NumericExpression {
	expr := &numericExpressionImpl{}
	expr.expressionImpl.initLiteralExpression(expr, value)
	return expr
}

func Float64(value float64) NumericExpression {
	expr := &numericExpressionImpl{}
	expr.expressionImpl.initLiteralExpression(expr, value)
	return expr
}

///////////////////////////////////////////////////////////////////////////////
// Other literals
///////////////////////////////////////////////////////////////////////////////

var (
	Asterisk = keyword("*")
)
