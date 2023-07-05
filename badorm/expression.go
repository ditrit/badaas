package badorm

import (
	"fmt"
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
