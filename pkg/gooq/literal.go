package gooq

import (
	"time"

	"github.com/google/uuid"
)

func keyword(value string) Expression {
	return newKeywordExpression(value)
}

func Literal(value interface{}) Expression {
	return newLiteralExpression(value)
}

func Bool(value bool) BoolExpression {
	expr := &boolExpressionImpl{}
	expr.expressionImpl.initLiteralExpression(expr, value)
	return expr
}

func DateTime(value time.Time) DateTimeExpression {
	expr := &dateTimeExpressionImpl{}
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

func String(value string) StringExpression {
	expr := &stringExpressionImpl{}
	expr.expressionImpl.initLiteralExpression(expr, value)
	return expr
}

func UUID(value uuid.UUID) UUIDExpression {
	expr := &uuidExpressionImpl{}
	expr.expressionImpl.initLiteralExpression(expr, value)
	return expr
}

///////////////////////////////////////////////////////////////////////////////
// Other literals
///////////////////////////////////////////////////////////////////////////////

var (
	Asterisk = keyword("*")
)
