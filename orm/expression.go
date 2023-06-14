package orm

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/elliotchance/pie/v2"
)

type Expression[T any] interface {
	// TODO agregar el metodo de validacion de interface
	ToSQL(columnName string) (string, []any)
}

// TODO
// string, int, etc. uuid, cualquier custom, time, deletedAt, asi que es any al final
// aunque algunos como like y eso solo funcionan para string, el problema es que yo no se si
// uno custom va a ir a string o no
// podria igual mirar que condiciones les genero y cuales no segun el tipo
type ValueExpression[T any] struct {
	// TODO creo que como no uso T esto no va a verificar nada,
	// aca antes habia []T pero me limita para cosas que no necesariamente comparan contra T como el startsWith
	Value         any
	SQLExpression string
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

// TODO aca me gustaria que devuelva []T pero no me anda asi
func (expr ValueExpression[T]) ToSQL(columnName string) (string, []any) {
	// sino que para punteros no haya equal nil?
	// TODO este chequeo deberia ser solo cuando T es un puntero
	// TODO y aca que pasa con time, deletedAt, y otros nullables por valuer
	// TODO que pasa para los demas symbols, puede meterme un null en un lt?
	// TODO esto esta feo
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

func NewValueExpression[T any](value T, sqlExpression string) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		SQLExpression: sqlExpression,
	}
}

type PredicateExpression[T any] struct {
	// TODO creo que como no uso T esto no va a verificar nada,
	// aca antes habia []T pero me limita para cosas que no necesariamente comparan contra T como el startsWith
	SQLExpression string
}

// TODO aca me gustaria que devuelva []T pero no me anda asi
func (expr PredicateExpression[T]) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf("%s %s", columnName, expr.SQLExpression), []any{}
}

func NewPredicateExpression[T any](sqlExpression string) PredicateExpression[T] {
	return PredicateExpression[T]{
		SQLExpression: sqlExpression,
	}
}

// Comparison Operators
func Eq[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, "=")
}

func NotEq[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, "<>")
}

func Lt[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, "<")
}

func LtOrEq[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, "<=")
}

// TODO no existe en psql
func NotLt[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, "!<")
}

func Gt[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, ">")
}

func GtOrEq[T any](value T) ValueExpression[T] {
	return NewValueExpression(value, ">=")
}

// Comparison Predicates

// TODO BETWEEN, NOT BETWEEN

// TODO no deberia ser posible para todos, solo los que son nullables
// pero como puedo saberlo, los que son pointers?, pero tambien hay otros como deletedAt que pueden ser null por su valuer
func IsNull[T any]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NULL")
}

func IsNotNull[T any]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

// TODO que pasa con otros que mapean a bool por valuer?
func IsTrue[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS TRUE")
}

func IsNotTrue[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT TRUE")
}

func IsFalse[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS FALSE")
}

func IsNotFalse[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT FALSE")
}

func IsUnknown[T *bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS UNKNOWN")
}

func IsNotUnknown[T *bool | sql.NullBool]() PredicateExpression[T] {
	return NewPredicateExpression[T]("IS NOT UNKNOWN")
}

// Row and Array Comparisons

func ArrayIn[T any](values ...T) ValueExpression[T] {
	return NewValueExpression[T](values, "IN")
}

func ArrayNotIn[T any](values ...T) ValueExpression[T] {
	return NewValueExpression[T](values, "NOT IN")
}
