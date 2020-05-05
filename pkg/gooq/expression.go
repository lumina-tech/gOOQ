package gooq

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Expression interface {
	Renderable

	As(alias string) Expression

	// https://www.postgresql.org/docs/11/functions-subquery.html
	In(subquery Selectable) BoolExpression
	NotIn(subquery Selectable) BoolExpression

	// Comparison Predicates
	// https://www.postgresql.org/docs/11/functions-comparison.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	IsNull() BoolExpression
	IsNotNull() BoolExpression

	// Indexes and ORDER BY
	// https://www.postgresql.org/docs/12/queries-order.html
	// [Good First Issue][Help Wanted] TODO: implement NULLS FIRST, NULL LAST
	Asc() Expression
	Desc() Expression

	// 4.2.7 Aggregate Expressions
	// https://www.postgresql.org/docs/12/sql-expressions.html
	Filter(...Expression) Expression

	// IMPORTANT: these are for internal use only.
	getExpressions() []Expression
	getOperator() Operator
	// getOriginal returns a reference to the getOriginal expression. This allows us to recover
	// an expression's getOriginal type. e.g. when we call TableImpl.column1.Eq(String("foo")
	// a newBooleanExpression is created with the operands TableImpl.column1 and
	// String("foo"). inArray the new boolean expression TableImpl.column1 is stored as a
	// stringExpressionImpl and has lost its original stringFieldImpl type. When
	// we render the expression it renders <nil> = 'foo' instead of TableImpl.column1 = 'foo'
	// because the stringExpressionImpl renderer was used.
	getOriginal() Expression
}

type ExpressionType int

const (
	ExpressionTypeExpressionArray = ExpressionType(iota)
	ExpressionTypeField
	ExpressionTypeFunction
	ExpressionTypeKeyword
	ExpressionTypeLiteral
	ExpressionTypeSubquery
	ExpressionTypeUnaryPrefix
	ExpressionTypeUnaryPostfix
	ExpressionTypeBinary
	// https://en.wikipedia.org/wiki/Plural_quantification
	ExpressionTypeMultigrade
)

type expressionImpl struct {
	originalExpression Expression
	expressionType     ExpressionType
	operator           Operator
	expressions        []Expression // can be operands or function arguments
	value              interface{}
	hasParentheses     bool
}

type ExpressionImplOption func(*expressionImpl)

func HasParentheses(hasParentheses bool) ExpressionImplOption {
	return func(impl *expressionImpl) {
		impl.hasParentheses = hasParentheses
	}
}

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

func newExpressionArray(
	value []Expression,
) Expression {
	return &expressionImpl{
		expressionType: ExpressionTypeExpressionArray,
		value:          value,
	}
}

func newSubquery(
	value Selectable, options ...ExpressionImplOption,
) Expression {
	instance := &expressionImpl{
		expressionType: ExpressionTypeSubquery,
		value:          value,
	}
	instance.apply(options...)
	return instance
}

func newUnaryPostfixExpression(
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

func (expr *expressionImpl) initMultigradeExpression(
	operator Operator, expressions []Expression, options ...ExpressionImplOption,
) {
	expr.expressionType = ExpressionTypeMultigrade
	expr.operator = operator
	expr.expressions = getOriginalExpressions(expressions)
	expr.apply(options...)
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

func (expr *expressionImpl) initUnaryPostfixExpression(
	operator Operator, operand Expression, options ...ExpressionImplOption,
) {
	expr.expressionType = ExpressionTypeUnaryPostfix
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
	return newAliasFunction(expr.getOriginal(), alias)
}

func (expr *expressionImpl) lt(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorLt, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) lte(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorLte, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) gt(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorGt, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) gte(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorGte, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) eq(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorEq, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) notEq(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorNotEq, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) like(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorLike, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) iLike(rhs Expression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorILike, expr.getOriginal(), rhs)
}

func (expr *expressionImpl) inArray(
	value []Expression,
) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorIn,
		expr.getOriginal(), newExpressionArray(value))
}

func (expr *expressionImpl) notInArray(
	value []Expression,
) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorNotIn,
		expr.getOriginal(), newExpressionArray(value))
}

func (expr *expressionImpl) In(
	subquery Selectable,
) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorIn,
		expr.getOriginal(), newSubquery(subquery, HasParentheses(true)))
}

func (expr *expressionImpl) NotIn(
	subquery Selectable,
) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorNotIn,
		expr.getOriginal(), newSubquery(subquery, HasParentheses(true)))
}

func (expr *expressionImpl) IsNull() BoolExpression {
	return newUnaryPostfixBooleanExpressionImpl(OperatorIsNull, expr.getOriginal())
}

func (expr *expressionImpl) IsNotNull() BoolExpression {
	return newUnaryPostfixBooleanExpressionImpl(OperatorIsNotNull, expr.getOriginal())
}

func (expr *expressionImpl) Asc() Expression {
	return newUnaryPostfixExpression(OperatorAsc, expr.getOriginal())
}

func (expr *expressionImpl) Desc() Expression {
	return newUnaryPostfixExpression(OperatorDesc, expr.getOriginal())
}

func (expr *expressionImpl) Filter(
	expressions ...Expression,
) Expression {
	return newFilterWhereFunction(expr.getOriginal(), expressions...)
}

func (expr *expressionImpl) Render(
	builder *Builder,
) {
	if expr.hasParentheses {
		builder.Print("(")
	}
	switch expr.expressionType {
	case ExpressionTypeExpressionArray:
		array := expr.value.([]Expression)
		builder.RenderExpressionArray(array)
	case ExpressionTypeKeyword:
		builder.Print(expr.value.(string))
	case ExpressionTypeLiteral:
		builder.RenderLiteral(expr.value)
	case ExpressionTypeSubquery:
		expr.value.(Selectable).Render(builder)
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
	case ExpressionTypeMultigrade:
		for index, expression := range expr.expressions {
			isLast := index == len(expr.expressions)-1
			builder.RenderExpression(expression)
			if !isLast {
				builder.Print(" ").Print(expr.operator.String()).Print(" ")
			}
		}
	default:
		panic(fmt.Errorf("invalid operatorType=%v", expr.operator))
	}
	if expr.hasParentheses {
		builder.Print(")")
	}
}

func (expr *expressionImpl) getExpressions() []Expression {
	return expr.expressions
}

func (expr *expressionImpl) getOriginal() Expression {
	if expr.originalExpression != nil {
		return expr.originalExpression
	}
	return expr
}

func (expr *expressionImpl) getOperator() Operator {
	return expr.operator
}

func getOriginalExpressions(
	expressions []Expression,
) []Expression {
	var results []Expression
	for _, expr := range expressions {
		results = append(results, expr.getOriginal())
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
	// Not() BoolExpression - should be a function since it is an unary prefix operator. It is more natural that way

	// comparison operators helper methods
	// https://www.postgresql.org/docs/11/functions-comparison.html
	Eq(rhs BoolExpression) BoolExpression
	NotEq(rhs BoolExpression) BoolExpression
	IsEq(rhs bool) BoolExpression
	IsNotEq(rhs bool) BoolExpression

	// https://www.postgresql.org/docs/11/functions-comparisons.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	IsIn(value ...bool) BoolExpression
	IsNotIn(value ...bool) BoolExpression

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

func newMultigradeBooleanExpressionImpl(
	operator Operator, expressions []Expression, options ...ExpressionImplOption,
) *boolExpressionImpl {
	instance := &boolExpressionImpl{}
	instance.expressionImpl.initMultigradeExpression(operator, expressions, options...)
	return instance
}

func newUnaryPostfixBooleanExpressionImpl(
	operator Operator, expr Expression,
) *boolExpressionImpl {
	instance := &boolExpressionImpl{}
	instance.expressionImpl.initUnaryPostfixExpression(operator, expr)
	return instance
}

func (expr *boolExpressionImpl) And(rhs BoolExpression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorAnd, expr, rhs,
		HasParentheses(true))
}

func (expr *boolExpressionImpl) Or(rhs BoolExpression) BoolExpression {
	return newBinaryBooleanExpressionImpl(OperatorOr, expr, rhs,
		HasParentheses(true))
}

func (expr *boolExpressionImpl) Eq(rhs BoolExpression) BoolExpression {
	return expr.expressionImpl.eq(rhs)
}

func (expr *boolExpressionImpl) NotEq(rhs BoolExpression) BoolExpression {
	return expr.expressionImpl.notEq(rhs)
}

func (expr *boolExpressionImpl) IsEq(rhs bool) BoolExpression {
	return expr.expressionImpl.eq(Bool(rhs))
}

func (expr *boolExpressionImpl) IsNotEq(rhs bool) BoolExpression {
	return expr.expressionImpl.notEq(Bool(rhs))
}

func (expr *boolExpressionImpl) IsIn(value ...bool) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.inArray(expressions)
}

func (expr *boolExpressionImpl) IsNotIn(value ...bool) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.notInArray(expressions)
}

///////////////////////////////////////////////////////////////////////////////
// Numeric
///////////////////////////////////////////////////////////////////////////////

type NumericExpression interface {
	Expression

	// comparison operators helper methods
	// https://www.postgresql.org/docs/11/functions-comparison.html
	Lt(rhs NumericExpression) BoolExpression
	Lte(rhs NumericExpression) BoolExpression
	Gt(rhs NumericExpression) BoolExpression
	Gte(rhs NumericExpression) BoolExpression
	Eq(rhs NumericExpression) BoolExpression
	NotEq(rhs NumericExpression) BoolExpression
	IsLt(rhs float64) BoolExpression
	IsLte(rhs float64) BoolExpression
	IsGt(rhs float64) BoolExpression
	IsGte(rhs float64) BoolExpression
	IsEq(rhs float64) BoolExpression
	IsNotEq(rhs float64) BoolExpression

	// https://www.postgresql.org/docs/11/functions-comparisons.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	IsIn(value ...float64) BoolExpression
	IsNotIn(value ...float64) BoolExpression

	// mathematical operators (non exhaustive)
	// https://www.postgresql.org/docs/11/functions-math.html
	// [Good First Issue][Help Wanted] TODO: implement remaining
	Add(rhs NumericExpression) NumericExpression
	Sub(rhs NumericExpression) NumericExpression
	Mult(rhs NumericExpression) NumericExpression
	Div(rhs NumericExpression) NumericExpression
	Sqrt() NumericExpression
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

func NewUnaryPrefixNumericExpressionImpl(
	operator Operator, operand NumericExpression,
) NumericExpression {
	instance := &numericExpressionImpl{}
	instance.expressionImpl.initUnaryPrefixExpression(operator, operand)
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

func (expr *numericExpressionImpl) Sqrt() NumericExpression {
	return NewUnaryPrefixNumericExpressionImpl(OperatorSqrt, expr)
}

func (expr *numericExpressionImpl) IsIn(value ...float64) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.inArray(expressions)
}

func (expr *numericExpressionImpl) IsNotIn(value ...float64) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.notInArray(expressions)
}

func (expr *numericExpressionImpl) Lt(rhs NumericExpression) BoolExpression {
	return expr.expressionImpl.lt(rhs)
}

func (expr *numericExpressionImpl) Lte(rhs NumericExpression) BoolExpression {
	return expr.expressionImpl.lte(rhs)
}

func (expr *numericExpressionImpl) Gt(rhs NumericExpression) BoolExpression {
	return expr.expressionImpl.gt(rhs)
}

func (expr *numericExpressionImpl) Gte(rhs NumericExpression) BoolExpression {
	return expr.expressionImpl.gte(rhs)
}

func (expr *numericExpressionImpl) Eq(rhs NumericExpression) BoolExpression {
	return expr.expressionImpl.eq(rhs)
}

func (expr *numericExpressionImpl) NotEq(rhs NumericExpression) BoolExpression {
	return expr.expressionImpl.notEq(rhs)
}

func (expr *numericExpressionImpl) IsLt(rhs float64) BoolExpression {
	return expr.expressionImpl.lt(Float64(rhs))
}

func (expr *numericExpressionImpl) IsLte(rhs float64) BoolExpression {
	return expr.expressionImpl.lte(Float64(rhs))
}

func (expr *numericExpressionImpl) IsGt(rhs float64) BoolExpression {
	return expr.expressionImpl.gt(Float64(rhs))
}

func (expr *numericExpressionImpl) IsGte(rhs float64) BoolExpression {
	return expr.expressionImpl.gte(Float64(rhs))
}

func (expr *numericExpressionImpl) IsEq(rhs float64) BoolExpression {
	return expr.expressionImpl.eq(Float64(rhs))
}

func (expr *numericExpressionImpl) IsNotEq(rhs float64) BoolExpression {
	return expr.expressionImpl.notEq(Float64(rhs))
}

///////////////////////////////////////////////////////////////////////////////
// String
///////////////////////////////////////////////////////////////////////////////

type StringExpression interface {
	Expression

	// comparison operators helper methods
	// https://www.postgresql.org/docs/11/functions-comparison.html
	Eq(rhs StringExpression) BoolExpression
	NotEq(rhs StringExpression) BoolExpression
	IsEq(rhs string) BoolExpression
	IsNotEq(rhs string) BoolExpression

	Like(value string) BoolExpression
	ILike(value string) BoolExpression

	// https://www.postgresql.org/docs/11/functions-comparisons.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	IsIn(value ...string) BoolExpression
	IsNotIn(value ...string) BoolExpression
}

type stringExpressionImpl struct {
	expressionImpl
}

func (expr *stringExpressionImpl) Eq(rhs StringExpression) BoolExpression {
	return expr.expressionImpl.eq(rhs)
}

func (expr *stringExpressionImpl) NotEq(rhs StringExpression) BoolExpression {
	return expr.expressionImpl.notEq(rhs)
}

func (expr *stringExpressionImpl) IsEq(rhs string) BoolExpression {
	return expr.expressionImpl.eq(String(rhs))
}

func (expr *stringExpressionImpl) IsNotEq(rhs string) BoolExpression {
	return expr.expressionImpl.notEq(String(rhs))
}

func (expr *stringExpressionImpl) Like(value string) BoolExpression {
	return expr.expressionImpl.like(String(value))
}

func (expr *stringExpressionImpl) ILike(value string) BoolExpression {
	return expr.expressionImpl.iLike(String(value))
}

func (expr *stringExpressionImpl) IsIn(value ...string) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.inArray(expressions)
}

func (expr *stringExpressionImpl) IsNotIn(value ...string) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.notInArray(expressions)
}

///////////////////////////////////////////////////////////////////////////////
// DateTime
///////////////////////////////////////////////////////////////////////////////

type DateTimeExpression interface {
	Expression

	// comparison operators helper methods
	// https://www.postgresql.org/docs/11/functions-comparison.html
	Lt(rhs DateTimeExpression) BoolExpression
	Lte(rhs DateTimeExpression) BoolExpression
	Gt(rhs DateTimeExpression) BoolExpression
	Gte(rhs DateTimeExpression) BoolExpression
	Eq(rhs DateTimeExpression) BoolExpression
	NotEq(rhs DateTimeExpression) BoolExpression
	IsLt(rhs time.Time) BoolExpression
	IsLte(rhs time.Time) BoolExpression
	IsGt(rhs time.Time) BoolExpression
	IsGte(rhs time.Time) BoolExpression
	IsEq(rhs time.Time) BoolExpression
	IsNotEq(rhs time.Time) BoolExpression

	// https://www.postgresql.org/docs/11/functions-comparisons.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	IsIn(value ...time.Time) BoolExpression
	IsNotIn(value ...time.Time) BoolExpression

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

func (expr *dateTimeExpressionImpl) Lt(rhs DateTimeExpression) BoolExpression {
	return expr.expressionImpl.lt(rhs)
}

func (expr *dateTimeExpressionImpl) Lte(rhs DateTimeExpression) BoolExpression {
	return expr.expressionImpl.lte(rhs)
}

func (expr *dateTimeExpressionImpl) Gt(rhs DateTimeExpression) BoolExpression {
	return expr.expressionImpl.gt(rhs)
}

func (expr *dateTimeExpressionImpl) Gte(rhs DateTimeExpression) BoolExpression {
	return expr.expressionImpl.gte(rhs)
}

func (expr *dateTimeExpressionImpl) Eq(rhs DateTimeExpression) BoolExpression {
	return expr.expressionImpl.eq(rhs)
}

func (expr *dateTimeExpressionImpl) NotEq(rhs DateTimeExpression) BoolExpression {
	return expr.expressionImpl.notEq(rhs)
}

func (expr *dateTimeExpressionImpl) IsLt(rhs time.Time) BoolExpression {
	return expr.expressionImpl.lt(DateTime(rhs))
}

func (expr *dateTimeExpressionImpl) IsLte(rhs time.Time) BoolExpression {
	return expr.expressionImpl.lte(DateTime(rhs))
}

func (expr *dateTimeExpressionImpl) IsGt(rhs time.Time) BoolExpression {
	return expr.expressionImpl.gt(DateTime(rhs))
}

func (expr *dateTimeExpressionImpl) IsGte(rhs time.Time) BoolExpression {
	return expr.expressionImpl.gte(DateTime(rhs))
}

func (expr *dateTimeExpressionImpl) IsEq(rhs time.Time) BoolExpression {
	return expr.expressionImpl.eq(DateTime(rhs))
}

func (expr *dateTimeExpressionImpl) IsNotEq(rhs time.Time) BoolExpression {
	return expr.expressionImpl.notEq(DateTime(rhs))
}

func (expr *dateTimeExpressionImpl) IsIn(value ...time.Time) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.inArray(expressions)
}

func (expr *dateTimeExpressionImpl) IsNotIn(value ...time.Time) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.notInArray(expressions)
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
// UUID
///////////////////////////////////////////////////////////////////////////////

type UUIDExpression interface {
	Expression

	// comparison operators helper methods
	// https://www.postgresql.org/docs/11/functions-comparison.html
	Eq(rhs UUIDExpression) BoolExpression
	NotEq(rhs UUIDExpression) BoolExpression
	IsEq(rhs uuid.UUID) BoolExpression
	IsNotEq(rhs uuid.UUID) BoolExpression

	// https://www.postgresql.org/docs/11/functions-comparisons.html
	// [Good First Issue][Help Wanted] TODO: implement remaining operators relevant for expression
	IsIn(value ...uuid.UUID) BoolExpression
	IsNotIn(value ...uuid.UUID) BoolExpression
}

type uuidExpressionImpl struct {
	expressionImpl
}

func (expr *uuidExpressionImpl) Eq(rhs UUIDExpression) BoolExpression {
	return expr.expressionImpl.eq(rhs)
}

func (expr *uuidExpressionImpl) NotEq(rhs UUIDExpression) BoolExpression {
	return expr.expressionImpl.notEq(rhs)
}

func (expr *uuidExpressionImpl) IsEq(rhs uuid.UUID) BoolExpression {
	return expr.expressionImpl.eq(UUID(rhs))
}

func (expr *uuidExpressionImpl) IsNotEq(rhs uuid.UUID) BoolExpression {
	return expr.expressionImpl.notEq(UUID(rhs))
}

func (expr *uuidExpressionImpl) IsIn(value ...uuid.UUID) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.inArray(expressions)
}

func (expr *uuidExpressionImpl) IsNotIn(value ...uuid.UUID) BoolExpression {
	expressions := getExpressionSlice(value)
	return expr.expressionImpl.notInArray(expressions)
}

///////////////////////////////////////////////////////////////////////////////
// Other Data Types
// https://www.postgresql.org/docs/11/datatype.html
// [Help Wanted] TODO: implement
///////////////////////////////////////////////////////////////////////////////

func getExpressionSlice(
	value interface{},
) []Expression {
	result := make([]Expression, 0)
	array := reflect.ValueOf(value)
	for i := 0; i < array.Len(); i++ {
		result = append(result, newLiteralExpression(array.Index(i).Interface()))
	}
	return result
}
