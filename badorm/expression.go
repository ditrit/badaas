package badorm

import "fmt"

// TODO
// string, int, etc. uuid, cualquier custom, time, deletedAt, asi que es any al final
// aunque algunos como like y eso solo funcionan para string, el problema es que yo no se si
// uno custom va a ir a string o no
// podria igual mirar que condiciones les genero y cuales no segun el tipo
type Expression struct {
	Value  any
	symbol string
}

func (expr Expression) ToSQL(columnName string) (string, []any) {
	return fmt.Sprintf("%s %s ?", columnName, expr.symbol), []any{expr.Value}
}

// Comparison Operators

func Eq(value any) Expression {
	return Expression{
		Value:  value,
		symbol: "=",
	}
}

func NotEq(value any) Expression {
	return Expression{
		Value:  value,
		symbol: "<>",
	}
}

func Lt(value any) Expression {
	return Expression{
		Value:  value,
		symbol: "<",
	}
}

func LtOrEq(value any) Expression {
	return Expression{
		Value:  value,
		symbol: "<=",
	}
}

func Gt(value any) Expression {
	return Expression{
		Value:  value,
		symbol: ">",
	}
}

func GtOrEq(value any) Expression {
	return Expression{
		Value:  value,
		symbol: ">=",
	}
}

// Comparison Predicates

// // TODO BETWEEN, NOT BETWEEN
// func IsDistinct(value any) Expression {
// 	return Expression{
// 		Value:  value,
// 		symbol: "IS DISTINCT FROM",
// 	}
// }

// func IsNotDistinct(value any) Expression {
// 	return Expression{
// 		Value:  value,
// 		symbol: "IS NOT DISTINCT FROM",
// 	}
// }

// // TODO no deberia ser posible para todos
// func IsNull(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NULL",
// 	}
// }

// func IsNotNull(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT NULL",
// 	}
// }

// // TODO no deberia ser posible para todos
// func IsTrue(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS TRUE",
// 	}
// }

// func IsNotTrue(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT TRUE",
// 	}
// }

// // TODO no deberia ser posible para todos
// func IsFalse(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS FALSE",
// 	}
// }

// func IsNotFalse(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT FALSE",
// 	}
// }

// func IsUnknown(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS UNKNOWN",
// 	}
// }

// func IsNotUnknown(value any) Expression {
// 	return Expression{
// 		Value:  value, // TODO ver aca que hago
// 		symbol: "IS NOT UNKNOWN",
// 	}
// }

// // TODO no se a que grupo pertenece

// func In[T []any](values T) Expression {
// 	return Expression{
// 		Value:  values,
// 		symbol: "IN",
// 	}
// }

// func NotIn[T []any](values T) Expression {
// 	return Expression{
// 		Value:  values,
// 		symbol: "NOT IN",
// 	}
// }
