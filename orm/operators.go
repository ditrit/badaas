package orm

import (
	"database/sql"
)

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
// IsNotDistinct must be used in cases where value can be NULL
func Eq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, "=")
}

// NotEqualTo
// IsDistinct must be used in cases where value can be NULL
func NotEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, "<>")
}

// LessThan
func Lt[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, "<")
}

// LessThanOrEqualTo
func LtOrEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, "<=")
}

// GreaterThan
func Gt[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, ">")
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value T) Operator[T] {
	return NewCantBeNullValueOperator[T](value, ">=")
}

// Comparison Predicates
// refs: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

// Equivalent to v1 < value < v2
func Between[T any](v1 T, v2 T) MultivalueOperator[T] {
	return NewMultivalueOperator("BETWEEN", "AND", "", "", v1, v2)
}

// Equivalent to NOT (v1 < value < v2)
func NotBetween[T any](v1 T, v2 T) MultivalueOperator[T] {
	return NewMultivalueOperator("NOT BETWEEN", "AND", "", "", v1, v2)
}

func IsNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

func IsTrue[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS TRUE")
}

func IsNotTrue[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT TRUE")
}

func IsFalse[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS FALSE")
}

func IsNotFalse[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT FALSE")
}

func IsUnknown[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS UNKNOWN")
}

func IsNotUnknown[T bool | sql.NullBool]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT UNKNOWN")
}

func IsDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T](value, "IS DISTINCT FROM")
}

func IsNotDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T](value, "IS NOT DISTINCT FROM")
}

// Row and Array Comparisons

func ArrayIn[T any](values ...T) ValueOperator[T] {
	return NewValueOperator[T](values, "IN")
}
