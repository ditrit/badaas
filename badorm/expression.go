package badorm

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm/expressions"
)

var (
	ErrNotRelated      = errors.New("value type not related with T")
	ErrValueCantBeNull = errors.New("value to compare can't be null")
)

type Expression[T any] interface {
	// Transform the Expression to a SQL string and a list of values to use in the query
	// columnName is used by the expression to determine which is the objective column.
	ToSQL(query *Query, columnName string) (string, []any, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T],
	// since if no method receives by parameter a type T,
	// any other Expression[T2] would also be considered a Expression[T].
	InterfaceVerificationMethod(T)
}

type DynamicExpression[T any] interface {
	Expression[T]

	SelectJoin(joinNumber uint) DynamicExpression[T]
}

// Expression that compares the value of the column against a fixed value
// If ExpressionsAndValues has multiple entries, comparisons will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueExpression[T any] struct {
	ExpressionsAndValues []SQLExpressionAndValue
}

type SQLExpressionAndValue struct {
	SQLExpression string
	Value         any
}

func (expr ValueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr ValueExpression[T]) ToSQL(_ *Query, columnName string) (string, []any, error) {
	exprString := columnName
	values := []any{}

	for _, sqlExprAndValue := range expr.ExpressionsAndValues {
		exprString += " " + sqlExprAndValue.SQLExpression + " ?"
		values = append(values, sqlExprAndValue.Value)
	}

	return exprString, values, nil
}

func NewValueExpression[T any](value any, sqlExpression expressions.SQLExpression) ValueExpression[T] {
	expr := ValueExpression[T]{}

	return expr.AddSQLExpression(value, sqlExpression)
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

func NewCantBeNullValueExpression[T any](value any, sqlExpression expressions.SQLExpression) Expression[T] {
	if value == nil || mapsToNull(value) {
		return NewInvalidExpression[T](ErrValueCantBeNull)
	}

	return NewValueExpression[T](value, sqlExpression)
}

func NewMustBePOSIXValueExpression[T string | sql.NullString](pattern string, sqlExpression expressions.SQLExpression) Expression[T] {
	_, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return NewInvalidExpression[T](err)
	}

	return NewValueExpression[T](pattern, sqlExpression)
}

func (expr *ValueExpression[T]) AddSQLExpression(value any, sqlExpression expressions.SQLExpression) ValueExpression[T] {
	expr.ExpressionsAndValues = append(
		expr.ExpressionsAndValues,
		SQLExpressionAndValue{
			Value:         value,
			SQLExpression: expressions.ToSQL[sqlExpression],
		},
	)

	return *expr
}

// Expression that compares the value of the column against multiple values
// Example: value IN (v1, v2, v3, ..., vN)
type MultivalueExpression[T any] struct {
	Values        []T    // the values to compare with
	SQLExpression string // the expression used to compare, example: IN
	SQLConnector  string // the connector between values, example: ', '
	SQLPrefix     string // something to put before the values, example: (
	SQLSuffix     string // something to put after the values, example: )
}

func (expr MultivalueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr MultivalueExpression[T]) ToSQL(_ *Query, columnName string) (string, []any, error) {
	placeholders := strings.Join(pie.Map(expr.Values, func(value T) string {
		return "?"
	}), " "+expr.SQLConnector+" ")

	values := pie.Map(expr.Values, func(value T) any {
		return value
	})

	return fmt.Sprintf(
		"%s %s %s"+placeholders+"%s",
		columnName,
		expr.SQLExpression,
		expr.SQLPrefix,
		expr.SQLSuffix,
	), values, nil
}

func NewMultivalueExpression[T any](sqlExpression expressions.SQLExpression, sqlConnector, sqlPrefix, sqlSuffix string, values ...T) MultivalueExpression[T] {
	return MultivalueExpression[T]{
		Values:        values,
		SQLExpression: expressions.ToSQL[sqlExpression],
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
	}
}

// Expression that verifies a predicate
// Example: value IS TRUE
type PredicateExpression[T any] struct {
	SQLExpression string
}

func (expr PredicateExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr PredicateExpression[T]) ToSQL(_ *Query, columnName string) (string, []any, error) {
	return fmt.Sprintf("%s %s", columnName, expr.SQLExpression), []any{}, nil
}

func NewPredicateExpression[T any](sqlExpression string) PredicateExpression[T] {
	return PredicateExpression[T]{
		SQLExpression: sqlExpression,
	}
}

// Expression used to return an error
type InvalidExpression[T any] struct {
	Err error
}

func (expr InvalidExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr InvalidExpression[T]) SelectJoin(joinNumber uint) DynamicExpression[T] {
	return expr
}

func (expr InvalidExpression[T]) ToSQL(_ *Query, _ string) (string, []any, error) {
	return "", nil, expr.Err
}

func NewInvalidExpression[T any](err error) InvalidExpression[T] {
	return InvalidExpression[T]{
		Err: err,
	}
}

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
// EqOrIsNull must be used in cases where value can be NULL
func Eq[T any](value T) Expression[T] {
	return NewCantBeNullValueExpression[T](value, expressions.Eq)
}

// if value is not NULL returns a Eq expression
// but if value is NULL returns a IsNull expression
// this must be used because ANSI SQL-92 standard defines:
// NULL = NULL evaluates to unknown, which is later considered a false
//
// this behavior can be also avoided in other ways:
//   - in SQLServer you can:
//     ** set ansi_nulls setting to off and use sqlserver.EqNullable
//     ** use the IS NOT DISTINCT operator (implemented in IsNotDistinct)
//   - in MySQL you can use the equal_to operator (implemented in mysql.IsEqual)
//   - in PostgreSQL you can use the IS NOT DISTINCT operator (implemented in IsNotDistinct)
//   - in SQLite you can use the IS NOT DISTINCT operator (implemented in IsNotDistinct)
func EqOrIsNull[T any](value any) Expression[T] {
	return expressionFromValueOrNil[T](value, Eq[T], IsNull[T]())
}

// NotEqualTo
// NotEqOrNotIsNull must be used in cases where value can be NULL
func NotEq[T any](value T) Expression[T] {
	return NewCantBeNullValueExpression[T](value, expressions.NotEq)
}

// if value is not NULL returns a NotEq expression
// but if value is NULL returns a IsNotNull expression
// this must be used because ANSI SQL-92 standard defines:
// NULL = NULL evaluates to unknown, which is later considered a false
//
// this behavior can be also avoided in other ways:
//   - in SQLServer you can:
//     ** set ansi_nulls setting to off and use sqlserver.NotEqNullable
//     ** use the IS DISTINCT operator (implemented in IsDistinct)
//   - in PostgreSQL you can use the IS DISTINCT operator (implemented in IsDistinct)
//   - in SQLite you can use the IS DISTINCT operator (implemented in IsDistinct)
func NotEqOrIsNotNull[T any](value any) Expression[T] {
	return expressionFromValueOrNil[T](value, NotEq[T], IsNotNull[T]())
}

func expressionFromValueOrNil[T any](value any, notNilFunc func(T) Expression[T], nilExpression Expression[T]) Expression[T] {
	if value == nil {
		return nilExpression
	}

	valueTPointer, isTPointer := value.(*T)
	if isTPointer {
		if valueTPointer == nil {
			return nilExpression
		}

		return notNilFunc(*valueTPointer)
	}

	valueT, isT := value.(T)
	if isT {
		if mapsToNull(value) {
			return nilExpression
		}

		return notNilFunc(valueT)
	}

	return NewInvalidExpression[T](ErrNotRelated)
}

func mapsToNull(value any) bool {
	reflectVal := reflect.ValueOf(value)
	isNullableKind := pie.Contains(nullableKinds, reflectVal.Kind())
	// avoid nil is not nil behavior of go
	if isNullableKind && reflectVal.IsNil() {
		return true
	}

	valuer, isValuer := value.(driver.Valuer)
	if isValuer {
		valuerValue, err := valuer.Value()
		if err == nil && valuerValue == nil {
			return true
		}
	}

	return false
}

// LessThan
func Lt[T any](value T) Expression[T] {
	return NewCantBeNullValueExpression[T](value, expressions.Lt)
}

// LessThanOrEqualTo
func LtOrEq[T any](value T) Expression[T] {
	return NewCantBeNullValueExpression[T](value, expressions.LtOrEq)
}

// GreaterThan
func Gt[T any](value T) Expression[T] {
	return NewCantBeNullValueExpression[T](value, expressions.Gt)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value T) Expression[T] {
	return NewCantBeNullValueExpression[T](value, expressions.GtOrEq)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](v1 T, v2 T) MultivalueExpression[T] {
	return NewMultivalueExpression(expressions.Between, "AND", "", "", v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1 T, v2 T) MultivalueExpression[T] {
	return NewMultivalueExpression(expressions.NotBetween, "AND", "", "", v1, v2)
}

func IsNull[T any]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NULL")
}

func IsNotNull[T any]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

// Not supported by: sqlserver
func IsTrue[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS TRUE")
}

// Not supported by: sqlserver
func IsNotTrue[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT TRUE")
}

// Not supported by: sqlserver
func IsFalse[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS FALSE")
}

// Not supported by: sqlserver
func IsNotFalse[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT FALSE")
}

// Not supported by: sqlserver, sqlite
func IsUnknown[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS UNKNOWN")
}

// Not supported by: sqlserver, sqlite
func IsNotUnknown[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT UNKNOWN")
}

// Not supported by: mysql
func IsDistinct[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, expressions.IsDistinct)
}

// Not supported by: mysql
func IsNotDistinct[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, expressions.IsNotDistinct)
}

// Pattern Matching

type LikeExpression[T string | sql.NullString] struct {
	ValueExpression[T]
}

func NewLikeExpression[T string | sql.NullString](pattern string, sqlExpression expressions.SQLExpression) LikeExpression[T] {
	return LikeExpression[T]{
		ValueExpression: NewValueExpression[T](pattern, sqlExpression),
	}
}

func (expr LikeExpression[T]) Escape(escape rune) ValueExpression[T] {
	return expr.AddSQLExpression(string(escape), expressions.Escape)
}

// Pattern in all databases:
//   - An underscore (_) in pattern stands for (matches) any single character.
//   - A percent sign (%) matches any sequence of zero or more characters.
//
// Additionally in SQLServer:
//   - Square brackets ([ ]) matches any single character within the specified range ([a-f]) or set ([abcdef]).
//   - [^] matches any single character not within the specified range ([^a-f]) or set ([^abcdef]).
//
// WARNINGS:
//   - SQLite: LIKE is case-insensitive unless case_sensitive_like pragma (https://www.sqlite.org/pragma.html#pragma_case_sensitive_like) is true.
//   - SQLServer, MySQL: the case-sensitivity depends on the collation used in compared column.
//   - PostgreSQL: LIKE is always case-sensitive, if you want case-insensitive use the ILIKE operator (implemented in psql.ILike)
//
// refs:
//   - mysql: https://dev.mysql.com/doc/refman/8.0/en/string-comparison-functions.html#operator_like
//   - postgresql: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-LIKE
//   - sqlserver: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/like-transact-sql?view=sql-server-ver16
//   - sqlite: https://www.sqlite.org/lang_expr.html#like
func Like[T string | sql.NullString](pattern string) LikeExpression[T] {
	return NewLikeExpression[T](pattern, expressions.Like)
}

// TODO Subquery Expressions
