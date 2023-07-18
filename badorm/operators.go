package badorm

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/elliotchance/pie/v2"

	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

// Comparison Operators
// refs:
// - MySQL: https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// - PostgreSQL: https://www.postgresql.org/docs/current/functions-comparison.html
// - SQLServer: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comparison-operators-transact-sql?view=sql-server-ver16
// - SQLite: https://www.sqlite.org/lang_expr.html

// EqualTo
// EqOrIsNull must be used in cases where value can be NULL
func Eq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](badormSQL.Eq, value)
}

// if value is not NULL returns a Eq operator
// but if value is NULL returns a IsNull operator
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
func EqOrIsNull[T any](value any) Operator[T] {
	return operatorFromValueOrNil[T](value, "EqOrIsNull", Eq[T], IsNull[T]())
}

// NotEqualTo
// NotEqOrNotIsNull must be used in cases where value can be NULL
func NotEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](badormSQL.NotEq, value)
}

// if value is not NULL returns a NotEq operator
// but if value is NULL returns a IsNotNull operator
// this must be used because ANSI SQL-92 standard defines:
// NULL = NULL evaluates to unknown, which is later considered a false
//
// this behavior can be also avoided in other ways:
//   - in SQLServer you can:
//     ** set ansi_nulls setting to off and use sqlserver.NotEqNullable
//     ** use the IS DISTINCT operator (implemented in IsDistinct)
//   - in PostgreSQL you can use the IS DISTINCT operator (implemented in IsDistinct)
//   - in SQLite you can use the IS DISTINCT operator (implemented in IsDistinct)
func NotEqOrIsNotNull[T any](value any) Operator[T] {
	return operatorFromValueOrNil[T](value, "NotEqOrIsNotNull", NotEq[T], IsNotNull[T]())
}

func operatorFromValueOrNil[T any](
	value any, operatorName string,
	notNilFunc func(T) Operator[T], nilOperator Operator[T],
) Operator[T] {
	if value == nil {
		return nilOperator
	}

	valueTPointer, isTPointer := value.(*T)
	if isTPointer {
		if valueTPointer == nil {
			return nilOperator
		}

		return notNilFunc(*valueTPointer)
	}

	valueT, isT := value.(T)
	if isT {
		if mapsToNull(value) {
			return nilOperator
		}

		return notNilFunc(valueT)
	}

	return NewInvalidOperator[T](
		notRelatedError[T](value, operatorName),
	)
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
func Lt[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](badormSQL.Lt, value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](badormSQL.LtOrEq, value)
}

// GreaterThan
func Gt[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](badormSQL.Gt, value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](badormSQL.GtOrEq, value)
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

func IsNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT NULL")
}

// Equivalent to v1 < attribute < v2
func Between[T any](v1 T, v2 T) Operator[T] {
	return NewMultivalueOperator(badormSQL.Between, badormSQL.And, "", "", v1, v2)
}

// Equivalent to NOT (v1 < attribute < v2)
func NotBetween[T any](v1 T, v2 T) Operator[T] {
	return NewMultivalueOperator(badormSQL.NotBetween, badormSQL.And, "", "", v1, v2)
}

// Boolean Comparison Predicates

// Not supported by: sqlserver
func IsTrue[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS TRUE")
}

// Not supported by: sqlserver
func IsNotTrue[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT TRUE")
}

// Not supported by: sqlserver
func IsFalse[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS FALSE")
}

// Not supported by: sqlserver
func IsNotFalse[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT FALSE")
}

// Not supported by: sqlserver, sqlite
func IsUnknown[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS UNKNOWN")
}

// Not supported by: sqlserver, sqlite
func IsNotUnknown[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT UNKNOWN")
}

// Not supported by: mysql
func IsDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T](badormSQL.IsDistinct, value)
}

// Not supported by: mysql
func IsNotDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T](badormSQL.IsNotDistinct, value)
}

// Pattern Matching

type LikeOperator[T string | sql.NullString] struct {
	ValueOperator[T]
}

func NewLikeOperator[T string | sql.NullString](pattern string, sqlOperator badormSQL.Operator) LikeOperator[T] {
	return LikeOperator[T]{
		ValueOperator: NewValueOperator[T](sqlOperator, pattern),
	}
}

func (expr LikeOperator[T]) Escape(escape rune) ValueOperator[T] {
	return *expr.AddOperation(badormSQL.Escape, string(escape))
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
func Like[T string | sql.NullString](pattern string) LikeOperator[T] {
	return NewLikeOperator[T](pattern, badormSQL.Like)
}

// Row and Array Comparisons

// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html#operator_in
func ArrayIn[T any](values ...T) Operator[T] {
	return NewMultivalueOperator(badormSQL.ArrayIn, badormSQL.Comma, "(", ")", values...)
}

func ArrayNotIn[T any](values ...T) Operator[T] {
	return NewMultivalueOperator(badormSQL.ArrayNotIn, badormSQL.Comma, "(", ")", values...)
}

// TODO Subquery Operators
