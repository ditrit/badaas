package badorm

import "fmt"

// TODO
// string, int, etc. uuid, cualquier custom, time, deletedAt, asi que es any al final
// aunque algunos como like y eso solo funcionan para string, el problema es que yo no se si
// uno custom va a ir a string o no
// podria igual mirar que condiciones les genero y cuales no segun el tipo
type Expression[T any] struct {
	Value  T
	symbol string
}

func (expr Expression[T]) ToSQL(columnName string) (string, []T) {
	return fmt.Sprintf("%s %s ?", columnName, expr.symbol), []T{expr.Value}
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
