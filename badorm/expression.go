package badorm

import (
	"fmt"
	"reflect"

	"github.com/elliotchance/pie/v2"
)

// TODO
// string, int, etc. uuid, cualquier custom, time, deletedAt, asi que es any al final
// aunque algunos como like y eso solo funcionan para string, el problema es que yo no se si
// uno custom va a ir a string o no
// podria igual mirar que condiciones les genero y cuales no segun el tipo
type Expression[T any] struct {
	Value  T
	symbol string
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

// TODO aca me gustaria que devuelva []T pero no me anda asi
func (expr Expression[T]) ToSQL(columnName string) (string, []any) {
	// TODO esto es valido solo para el equal, para los demas no se deberia hacer esto
	// sino que para punteros no haya equal nil?
	// TODO este chequeo deberia ser solo cuando T es un puntero
	reflectVal := reflect.ValueOf(expr.Value)
	isNullableKind := pie.Contains(nullableKinds, reflectVal.Kind())
	// avoid nil is not nil behavior of go
	if isNullableKind && reflectVal.IsNil() {
		return fmt.Sprintf(
			"%s IS NULL",
			columnName,
		), []any{}
	}

	return fmt.Sprintf("%s %s ?", columnName, expr.symbol), []any{expr.Value}
}

// Comparison Operators

func Eq[T any](value T) Expression[T] {
	return Expression[T]{
		Value:  value,
		symbol: "=",
	}
}

func NotEq[T any](value T) Expression[T] {
	return Expression[T]{
		Value:  value,
		symbol: "<>",
	}
}

func Lt[T any](value T) Expression[T] {
	return Expression[T]{
		Value:  value,
		symbol: "<",
	}
}

func LtOrEq[T any](value T) Expression[T] {
	return Expression[T]{
		Value:  value,
		symbol: "<=",
	}
}

func Gt[T any](value T) Expression[T] {
	return Expression[T]{
		Value:  value,
		symbol: ">",
	}
}

func GtOrEq[T any](value T) Expression[T] {
	return Expression[T]{
		Value:  value,
		symbol: ">=",
	}
}

// Comparison Predicates

// // TODO BETWEEN, NOT BETWEEN
// func IsDistinct[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value,
// 		symbol: "IS DISTINCT FROM",
// 	}
// }

// func IsNotDistinct[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value,
// 		symbol: "IS NOT DISTINCT FROM",
// 	}
// }

// // TODO no deberia ser posible para todos
// func IsNull[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NULL",
// 	}
// }

// func IsNotNull[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT NULL",
// 	}
// }

// // TODO no deberia ser posible para todos
// func IsTrue[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS TRUE",
// 	}
// }

// func IsNotTrue[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT TRUE",
// 	}
// }

// // TODO no deberia ser posible para todos
// func IsFalse[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS FALSE",
// 	}
// }

// func IsNotFalse[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT FALSE",
// 	}
// }

// func IsUnknown[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS UNKNOWN",
// 	}
// }

// func IsNotUnknown[T any](value T) Expression[T] {
// 	return Expression[T]{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT UNKNOWN",
// 	}
// }

// // TODO no se a que grupo pertenece

// func In[T []any](values T) Expression[T] {
// 	return Expression[T]{
// 		Value:  values,
// 		symbol: "IN",
// 	}
// }

// func NotIn[T []any](values T) Expression[T] {
// 	return Expression[T]{
// 		Value:  values,
// 		symbol: "NOT IN",
// 	}
// }