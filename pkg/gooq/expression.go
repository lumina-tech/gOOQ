package gooq

import (
	"fmt"
)

type Expression interface {
	Renderable
	As(alias string) Expression
	Filter(...Expression) Expression

	// Comparison Operators
	// https://www.postgresql.org/docs/11/functions-comparison.html
	Lt(rhs Expression) BoolExpression
	Lte(rhs Expression) BoolExpression
	Gt(rhs Expression) BoolExpression
	Gte(rhs Expression) BoolExpression
	Eq(rhs Expression) BoolExpression
	NotEq(rhs Expression) BoolExpression

	// https://www.postgresql.org/docs/11/functions-comparisons.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	In(value interface{}, rest ...interface{}) BoolExpression
	NotIn(value interface{}, rest ...interface{}) BoolExpression

	// Comparison Predicates
	// https://www.postgresql.org/docs/11/functions-comparison.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression

	// Indexes and ORDER BY
	// https://www.postgresql.org/docs/11/indexes-ordering.html
	// [Good First Issue][Help Wanted] TODO: implement NULLS FIRST, NULL LAST
	Asc() Expression
	Desc() Expression

	// IMPORTANT: this is for internal use only.
	// original returns a reference to the original expression. This allows us to recover
	// an expression's original type. e.g. when we call TableImpl.column1.Eq(String("foo")
	// a newBooleanExpression is created with the operands TableImpl.column1 and
	// String("foo"). In the new boolean expression TableImpl.column1 is stored as a
	// stringExpressionImpl and has lost its original stringFieldImpl type. When
	// we render the expression it renders <nil> = 'foo' instead of TableImpl.column1 = 'foo'
	// because the stringExpressionImpl renderer was used.
	original() Expression
}

type ExpressionType int

const (
	ExpressionTypeField = ExpressionType(iota)
	ExpressionTypeFunction
	ExpressionTypeKeyword
	ExpressionTypeLiteral
	ExpressionTypeLiteralArray
	ExpressionTypeUnaryPrefix
	ExpressionTypeUnaryPostfix
	ExpressionTypeBinary
)

type expressionImpl struct {
	originalExpression Expression
	expressionType     ExpressionType
	operator           Operator
	expressions        []Expression // can be operands or function arguments
	value              interface{}
}

type ExpressionImplOption func(*expressionImpl)

func newKeywordExpression(
	value string,
) Expression {
	return &expressionImpl{
		expressionType: ExpressionTypeKeyword,
		value:          value,
	}
}

func newLiteralExpression(
	value interface{},
) Expression {
	return &expressionImpl{
		expressionType: ExpressionTypeLiteral,
		value:          value,
	}
}

func newLiteralArrayExpression(
	value []interface{},
) Expression {
	return &expressionImpl{
		expressionType: ExpressionTypeLiteralArray,
		value:          value,
	}
}

func newPostfixUnaryExpression(
	operator Operator, operand Expression, options ...ExpressionImplOption,
) Expression {
	instance := &expressionImpl{
		expressionType: ExpressionTypeUnaryPostfix,
		operator:       operator,
		expressions:    getOriginalExpressions([]Expression{operand}),
	}
	instance.apply(options...)
	return instance
}

func (expr *expressionImpl) initFieldExpressionImpl(
	original Expression,
) {
	expr.originalExpression = original
	expr.expressionType = ExpressionTypeField
}

func (expr *expressionImpl) initFunctionExpression(
	original Expression, arguments ...Expression,
) {
	expr.originalExpression = original
	expr.expressionType = ExpressionTypeFunction
	expr.expressions = getOriginalExpressions(arguments)
}

func (expr *expressionImpl) initLiteralExpression(
	original Expression, value interface{},
) {
	expr.originalExpression = original
	expr.expressionType = ExpressionTypeLiteral
	expr.value = value
}

func (expr *expressionImpl) initBinaryExpression(
	operator Operator, lhs, rhs Expression, options ...ExpressionImplOption,
) {
	expr.expressionType = ExpressionTypeBinary
	expr.operator = operator
	expr.expressions = getOriginalExpressions([]Expression{lhs, rhs})
	expr.apply(options...)
}

func (expr *expressionImpl) initUnaryPrefixExpression(
	operator Operator, operand Expression, options ...ExpressionImplOption,
) {
	expr.expressionType = ExpressionTypeUnaryPrefix
	expr.operator = operator
	expr.expressions = getOriginalExpressions([]Expression{operand})
	expr.apply(options...)
}

func (expr *expressionImpl) apply(
	options ...ExpressionImplOption,
) *expressionImpl {
	for _, option := range options {
		option(expr)
	}
	return expr
}

func (expr *expressionImpl) As(alias string) Expression {
	return newAliasFunction(expr.original(), alias)
}

func (expr *expressionImpl) Filter(
	expressions ...Expression,
) Expression {
	return newFilterWhereFunction(expr.original(), expressions...)
}

func (expr *expressionImpl) Lt(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorLt, expr, rhs)
}

func (expr *expressionImpl) Lte(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorLte, expr, rhs)
}

func (expr *expressionImpl) Gt(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorGt, expr, rhs)
}

func (expr *expressionImpl) Gte(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorGte, expr, rhs)
}

func (expr *expressionImpl) Eq(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorEq, expr, rhs)
}

func (expr *expressionImpl) NotEq(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorNotEq, expr, rhs)
}

func (expr *expressionImpl) In(
	value interface{}, rest ...interface{},
) BoolExpression {
	array := LiteralArray(append([]interface{}{value}, rest...))
	return newBinaryBooleanExpressionImpl(OperatorIn, expr, array)
}

func (expr *expressionImpl) NotIn(
	value interface{}, rest ...interface{},
) BoolExpression {
	array := LiteralArray(append([]interface{}{value}, rest...))
	return newBinaryBooleanExpressionImpl(OperatorNotIn, expr, array)
}

func (expr *expressionImpl) Asc() Expression {
	return newPostfixUnaryExpression(OperatorAsc, expr)
}

func (expr *expressionImpl) Desc() Expression {
	return newPostfixUnaryExpression(OperatorDesc, expr)
}

func (expr *expressionImpl) Render(
	builder *Builder,
) {
	switch expr.expressionType {
	case ExpressionTypeKeyword:
		builder.Print(expr.value.(string))
	case ExpressionTypeLiteral:
		builder.RenderLiteral(expr.value)
	case ExpressionTypeLiteralArray:
		array := expr.value.([]interface{})
		builder.RenderLiteralArray(array)
	case ExpressionTypeUnaryPrefix:
		operand := expr.expressions[0]
		builder.Print(expr.operator.String()).Print(" ").
			RenderExpression(operand)
	case ExpressionTypeUnaryPostfix:
		operand := expr.expressions[0]
		builder.RenderExpression(operand).Print(" ").
			Print(expr.operator.String())
	case ExpressionTypeBinary:
		lhs := expr.expressions[0]
		rhs := expr.expressions[1]
		builder.RenderExpression(lhs).
			Print(" ").Print(expr.operator.String()).Print(" ").
			RenderExpression(rhs)
	default:
		panic(fmt.Errorf("invalid operatorType=%v", expr.operator))
	}
}

func (expr *expressionImpl) original() Expression {
	if expr.originalExpression != nil {
		return expr.originalExpression
	}
	return expr
}

func getOriginalExpressions(
	expressions []Expression,
) []Expression {
	var results []Expression
	for _, expr := range expressions {
		results = append(results, expr.original())
	}
	return results
}

///////////////////////////////////////////////////////////////////////////////
// Boolean
///////////////////////////////////////////////////////////////////////////////

type BoolExpression interface {
	Expression

	// logical operators
	// https://www.postgresql.org/docs/11/functions-logical.html
	And(rhs BoolExpression) BoolExpression
	Or(rhs BoolExpression) BoolExpression
	Not() BoolExpression

	// comparison Predicates
	// https://www.postgresql.org/docs/11/functions-logical.html
	// [Good First Issue][Help Wanted] TODO: implement remaining relevant for boolean_expression

}

type boolExpressionImpl struct {
	expressionImpl
}

func newBinaryBooleanExpressionImpl(
	operator Operator, lhs, rhs Expression, options ...ExpressionImplOption,
) *boolExpressionImpl {
	instance := &boolExpressionImpl{}
	instance.expressionImpl.initBinaryExpression(operator, lhs, rhs, options...)
	return instance
}

func NewUnaryPrefixBooleanExpressionImpl(
	operator Operator, expr Expression,
) *boolExpressionImpl {
	instance := &boolExpressionImpl{}
	instance.expressionImpl.initUnaryPrefixExpression(operator, expr)
	return instance
}

func (expr *boolExpressionImpl) And(expression BoolExpression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorAnd, expr, expression)
}

func (expr *boolExpressionImpl) Or(expression BoolExpression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorOr, expr, expression)
}

func (expr *boolExpressionImpl) Not() BoolExpression {
	return NewUnaryPrefixBooleanExpressionImpl(OperatorOr, expr)
}

///////////////////////////////////////////////////////////////////////////////
// Numeric
///////////////////////////////////////////////////////////////////////////////

type NumericExpression interface {
	Expression

	// mathematical operators (non exhaustive)
	// https://www.postgresql.org/docs/11/functions-math.html
	// [Good First Issue][Help Wanted] TODO: implement remaining
	Add(rhs NumericExpression) NumericExpression
	Sub(rhs NumericExpression) NumericExpression
	Mult(rhs NumericExpression) NumericExpression
	Div(rhs NumericExpression) NumericExpression
}

type numericExpressionImpl struct {
	expressionImpl
}

func NewBinaryNumericExpressionImpl(
	operator Operator, lhs, rhs NumericExpression,
) NumericExpression {
	instance := &numericExpressionImpl{}
	instance.expressionImpl.initBinaryExpression(operator, lhs, rhs)
	return instance
}

func (expr *numericExpressionImpl) Add(rhs NumericExpression) NumericExpression {
	return NewBinaryNumericExpressionImpl(OperatorAdd, expr, rhs)
}

func (expr *numericExpressionImpl) Sub(rhs NumericExpression) NumericExpression {
	return NewBinaryNumericExpressionImpl(OperatorSub, expr, rhs)
}

func (expr *numericExpressionImpl) Mult(rhs NumericExpression) NumericExpression {
	return NewBinaryNumericExpressionImpl(OperatorMult, expr, rhs)
}

func (expr *numericExpressionImpl) Div(rhs NumericExpression) NumericExpression {
	return NewBinaryNumericExpressionImpl(OperatorDiv, expr, rhs)
}

///////////////////////////////////////////////////////////////////////////////
// String
///////////////////////////////////////////////////////////////////////////////

type StringExpression interface {
	Expression
}

type stringExpressionImpl struct {
	expressionImpl
}

///////////////////////////////////////////////////////////////////////////////
// DateTime
///////////////////////////////////////////////////////////////////////////////

type DateTimeExpression interface {
	Expression

	// Table 9.29. Date/Time Operators
	// https://www.postgresql.org/docs/11/functions-datetime.html
	Add(rhs DateTimeExpression) DateTimeExpression
	Sub(rhs DateTimeExpression) DateTimeExpression
	Mult(rhs DateTimeExpression) DateTimeExpression
	Div(rhs DateTimeExpression) DateTimeExpression
}

type dateTimeExpressionImpl struct {
	expressionImpl
}

func newDateTimeExpressionImpl(
	operator Operator, lhs, rhs DateTimeExpression,
) DateTimeExpression {
	instance := &dateTimeExpressionImpl{}
	instance.expressionImpl.initBinaryExpression(operator, lhs, rhs)
	return instance
}

func (expr *dateTimeExpressionImpl) Add(rhs DateTimeExpression) DateTimeExpression {
	return newDateTimeExpressionImpl(OperatorAdd, expr, rhs)
}

func (expr *dateTimeExpressionImpl) Sub(rhs DateTimeExpression) DateTimeExpression {
	return newDateTimeExpressionImpl(OperatorSub, expr, rhs)
}

func (expr *dateTimeExpressionImpl) Mult(rhs DateTimeExpression) DateTimeExpression {
	return newDateTimeExpressionImpl(OperatorMult, expr, rhs)
}

func (expr *dateTimeExpressionImpl) Div(rhs DateTimeExpression) DateTimeExpression {
	return newDateTimeExpressionImpl(OperatorDiv, expr, rhs)
}

///////////////////////////////////////////////////////////////////////////////
// Other Data Types
// https://www.postgresql.org/docs/11/datatype.html
// [Help Wanted] TODO: implement
///////////////////////////////////////////////////////////////////////////////
