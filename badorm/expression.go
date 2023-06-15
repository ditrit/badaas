package badorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/elliotchance/pie/v2"
)

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
	// TODO este chequeo deberia ser solo cuando T es un puntero
	// y que pasa con time, deletedAt, y otros nullables por valuer
	// TODO que pasa para los demas symbols, puede meterme un null en un lt?
	// TODO esto esta feo
	// TODO tambien lo que hace la libreria esa es transformarlo en in si es un array
	// TODO ahora solo es util para los arrays y eso, para pointers ya no existe esta posibilidad
	if expr.SQLExpression == "=" {
		reflectVal := reflect.ValueOf(expr.Value)
		isNullableKind := pie.Contains(nullableKinds, reflectVal.Kind())
		// avoid nil is not nil behavior of go
		if isNullableKind && reflectVal.IsNil() {
			return fmt.Sprintf(
				"%s IS NULL",
				columnName,
			), []any{}
		}
	}

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

func Eq[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, "=")
}

func NotEq[T any](value T) ValueExpression[T] {
	return NewValueExpression[T](value, "<>")
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
