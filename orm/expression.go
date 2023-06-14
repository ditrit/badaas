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
	// TODO creo que como no uso T esto no va a verificar nada, aca antes habia []T pero me limita para cosas que no necesariamente comparan contra T como el startsWith
	Value         any
	sqlExpression string
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
	if expr.sqlExpression == "=" {
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

	return fmt.Sprintf("%s %s ?", columnName, expr.sqlExpression), []any{expr.Value}
}

type PredicateExpression[T any] struct {
	// TODO creo que como no uso T esto no va a verificar nada, aca antes habia []T pero me limita para cosas que no necesariamente comparan contra T como el startsWith
	sqlExpression string
}

// TODO aca me gustaria que devuelva []T pero no me anda asi
func (expr PredicateExpression[T]) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf("%s %s", columnName, expr.sqlExpression), []any{}
}

// Comparison Operators
// TODO aca hay codigo repetido entre los constructores
func Eq[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: "=",
	}
}

func NotEq[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: "<>",
	}
}

func Lt[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: "<",
	}
}

func LtOrEq[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: "<=",
	}
}

func Gt[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: ">",
	}
}

func GtOrEq[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: ">=",
	}
}

// Comparison Predicates

// TODO BETWEEN, NOT BETWEEN

func IsDistinct[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: "IS DISTINCT FROM",
	}
}

func IsNotDistinct[T any](value T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         value,
		sqlExpression: "IS NOT DISTINCT FROM",
	}
}

// TODO no deberia ser posible para todos, solo los que son nullables
// pero como puedo saberlo, los que son pointers?, pero tambien hay otros como deletedAt que pueden ser null por su valuer
func IsNull[T any]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS NULL",
	}
}

func IsNotNull[T any]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS NOT NULL",
	}
}

// Boolean Comparison Predicates

// TODO que pasa con otros que mapean a bool por valuer?
func IsTrue[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS TRUE",
	}
}

func IsNotTrue[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS NOT TRUE",
	}
}

func IsFalse[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS FALSE",
	}
}

func IsNotFalse[T bool | *bool | sql.NullBool]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS NOT FALSE",
	}
}

func IsUnknown[T *bool | sql.NullBool]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS UNKNOWN",
	}
}

func IsNotUnknown[T *bool | sql.NullBool]() PredicateExpression[T] {
	return PredicateExpression[T]{
		sqlExpression: "IS NOT UNKNOWN",
	}
}

// Row and Array Comparisons

// Row and Array Comparisons

func ArrayIn[T any](values ...T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         values,
		sqlExpression: "IN",
	}
}

func ArrayNotIn[T any](values ...T) ValueExpression[T] {
	return ValueExpression[T]{
		Value:         values,
		sqlExpression: "NOT IN",
	}
}
