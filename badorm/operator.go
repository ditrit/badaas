package badorm

import (
	"fmt"
)

type Operator[T any] interface {
	// Transform the Operator to a SQL string and a list of values to use in the query
	// columnName is used by the operator to determine which is the objective column.
	ToSQL(query *Query, columnName string) (string, []any, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T],
	// since if no method receives by parameter a type T,
	// any other Operator[T2] would also be considered a Operator[T].
	InterfaceVerificationMethod(T)
}

type DynamicOperator[T any] interface {
	Operator[T]

	// Allows to choose which number of join use
	// for the value in position "valueNumber"
	// when the value is a field and its model is joined more than once.
	// Does nothing if the valueNumber is bigger than the amount of values.
	SelectJoin(valueNumber, joinNumber uint) DynamicOperator[T]
}

// Operator that verifies a predicate
// Example: value IS TRUE
type PredicateOperator[T any] struct {
	SQLOperator string
}

func (expr PredicateOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (expr PredicateOperator[T]) ToSQL(_ *Query, columnName string) (string, []any, error) {
	return fmt.Sprintf("%s %s", columnName, expr.SQLOperator), []any{}, nil
}

func NewPredicateOperator[T any](sqlOperator string) PredicateOperator[T] {
	return PredicateOperator[T]{
		SQLOperator: sqlOperator,
	}
}

// Operator used to return an error
type InvalidOperator[T any] struct {
	Err error
}

func (expr InvalidOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

// InvalidOperator has SelectJoin to implement DynamicOperator
func (expr InvalidOperator[T]) SelectJoin(_, _ uint) DynamicOperator[T] {
	return expr
}

func (expr InvalidOperator[T]) ToSQL(_ *Query, _ string) (string, []any, error) {
	return "", nil, expr.Err
}

func NewInvalidOperator[T any](err error) InvalidOperator[T] {
	return InvalidOperator[T]{
		Err: err,
	}
}
