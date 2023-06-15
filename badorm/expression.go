package badorm

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/elliotchance/pie/v2"
)

var ErrNotRelated = errors.New("value type not related with T")

type Expression[T any] interface {
	ToSQL(columnName string) (string, []any)

	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T],
	// since if no method receives by parameter a type T,
	// any other Expression[T2] would also be considered a Expression[T].
	InterfaceVerificationMethod(T)
}

type ValueExpression[T any] struct {
	Value         any
	SQLExpression string
}

//nolint:unused // see inside
func (expr ValueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

// TODO aca me gustaria que devuelva []T pero no me anda asi
func (expr ValueExpression[T]) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf("%s %s ?", columnName, expr.SQLExpression), []any{expr.Value}
}

func NewValueExpression[T any](value any, sqlExpression string) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		SQLExpression: sqlExpression,
	}
}

type MultivalueExpression[T any] struct {
	Values        []T
	SQLExpression string
	SQLConnector  string
	SQLPrefix     string
	SQLSuffix     string
}

//nolint:unused // see inside
func (expr MultivalueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr MultivalueExpression[T]) ToSQL(columnName string) (string, []any) {
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
	), values
}

func NewMultivalueExpression[T any](sqlExpression, sqlConnector, sqlPrefix, sqlSuffix string, values ...T) MultivalueExpression[T] {
	return MultivalueExpression[T]{
		Values:        values,
		SQLExpression: sqlExpression,
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
	}
}

type PredicateExpression[T any] struct {
	SQLExpression string
}

//nolint:unused // see inside
func (expr PredicateExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr PredicateExpression[T]) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf("%s %s", columnName, expr.SQLExpression), []any{}
}

func NewPredicateExpression[T any](sqlExpression string) PredicateExpression[T] {
	return PredicateExpression[T]{
		SQLExpression: sqlExpression,
	}
}

// Comparison Operators
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html

// EqOrIsNull must be used in cases where value can be NULL
func Eq[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, "=")
}

// if value is not NULL returns a Eq expression
// but if value is NULL returns a IsNull expression
// this must be used because ANSI SQL-92 standard defines:
// NULL = NULL evaluates to unknown, which is later considered a false
//
// this behavior can be also avoided in other ways as:
// * in SQLServer you can set ansi_nulls setting to off
// * in MySQL you can use equal_to operator (implemented in mysql.IsEqual)
// * in PostgreSQL you can use the IS NOT DISTINCT operator (implemented in psql.IsNotDistinct)
func EqOrIsNull[T any](value any) (Expression[T], error) {
	return expressionFromValueOrNil[T](value, Eq[T], IsNull[T]())
}

// NotEqOrNotIsNull must be used in cases where value can be NULL
func NotEq[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, "<>")
}

// if value is not NULL returns a NotEq expression
// but if value is NULL returns a IsNotNull expression
// this must be used because ANSI SQL-92 standard defines:
// NULL = NULL evaluates to unknown, which is later considered a false
//
// this behavior can be also avoided in other ways as:
// * in SQLServer you can set ansi_nulls setting to off
// * in PostgreSQL you can use the IS DISTINCT operator (implemented in psql.IsDistinct)
func NotEqOrIsNotNull[T any](value any) (Expression[T], error) {
	return expressionFromValueOrNil[T](value, NotEq[T], IsNotNull[T]())
}

func expressionFromValueOrNil[T any](value any, eqFunc func(T) ValueExpression[T], nilExpression Expression[T]) (Expression[T], error) {
	if value == nil {
		return nilExpression, nil
	}

	valueTPointer, isTPointer := value.(*T)
	if isTPointer {
		if valueTPointer == nil {
			return nilExpression, nil
		}

		return eqFunc(*valueTPointer), nil
	}

	valueT, isT := value.(T)
	if isT {
		reflectVal := reflect.ValueOf(value)
		isNullableKind := pie.Contains(nullableKinds, reflectVal.Kind())
		// avoid nil is not nil behavior of go
		if isNullableKind && reflectVal.IsNil() {
			return nilExpression, nil
		}

		// TODO creo que esto lo voy a tener que mover afuera si quiero que los nullables se comparen contra los no nullables
		valuer, isValuer := value.(driver.Valuer)
		if isValuer {
			valuerValue, err := valuer.Value()
			if err == nil && valuerValue == nil {
				return nilExpression, nil
			}
		}

		return eqFunc(valueT), nil
	}

	return nil, ErrNotRelated
}

func Lt[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, "<")
}

func LtOrEq[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, "<=")
}

func Gt[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, ">")
}

func GtOrEq[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, ">=")
}

// Comparison Predicates
// refs:
// https://dev.mysql.com/doc/refman/8.0/en/comparison-operators.html
// https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

func Between[T any](v1 T, v2 T) MultivalueExpression[T] {
	return NewMultivalueExpression("BETWEEN", "AND", "", "", v1, v2)
}

func NotBetween[T any](v1 T, v2 T) MultivalueExpression[T] {
	return NewMultivalueExpression("NOT BETWEEN", "AND", "", "", v1, v2)
}

func IsNull[T any]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NULL")
}

func IsNotNull[T any]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

// TODO que pasa con otros que mapean a bool por valuer?
func IsTrue[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS TRUE")
}

func IsNotTrue[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT TRUE")
}

func IsFalse[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS FALSE")
}

func IsNotFalse[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT FALSE")
}

func IsUnknown[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS UNKNOWN")
}

func IsNotUnknown[T bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT UNKNOWN")
}
