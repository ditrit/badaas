package orm

import (
	"database/sql/driver"
	"errors"
)

var ErrValueCantBeNull = errors.New("value to compare can't be null")

type Operator[T any] interface {
	// Transform the Operator to a SQL string and a list of values to use in the query
	// columnName is used by the operator to determine which is the objective column.
	ToSQL(columnName string) (string, []any, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T],
	// since if no method receives by parameter a type T,
	// any other Operator[T2] would also be considered a Operator[T].
	InterfaceVerificationMethod(T)
}

// Operator that compares the value of the column against a fixed value
// If Operations has multiple entries, operations will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	Operations []Operation
}

type Operation struct {
	SQLOperator string
	Value       any
}

func (expr ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (expr ValueOperator[T]) ToSQL(columnName string) (string, []any, error) {
	operatorString := columnName
	values := []any{}

	for _, operation := range expr.Operations {
		operatorString += " " + operation.SQLOperator + " ?"
		values = append(values, operation.Value)
	}

	return operatorString, values, nil
}

func NewValueOperator[T any](value any, sqlOperator string) ValueOperator[T] {
	expr := ValueOperator[T]{}

	return expr.AddOperation(value, sqlOperator)
}

func NewCantBeNullValueOperator[T any](value any, sqlOperator string) Operator[T] {
	if value == nil || mapsToNull(value) {
		return NewInvalidOperator[T](ErrValueCantBeNull)
	}

	return NewValueOperator[T](value, sqlOperator)
}

func mapsToNull(value any) bool {
	valuer, isValuer := value.(driver.Valuer)
	if isValuer {
		valuerValue, err := valuer.Value()
		if err == nil && valuerValue == nil {
			return true
		}
	}

	return false
}

func (expr *ValueOperator[T]) AddOperation(value any, sqlOperator string) ValueOperator[T] {
	expr.Operations = append(
		expr.Operations,
		Operation{
			Value:       value,
			SQLOperator: sqlOperator,
		},
	)

	return *expr
}

// Operator used to return an error
type InvalidOperator[T any] struct {
	Err error
}

func (expr InvalidOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (expr InvalidOperator[T]) ToSQL(_ string) (string, []any, error) {
	return "", nil, expr.Err
}

func NewInvalidOperator[T any](err error) InvalidOperator[T] {
	return InvalidOperator[T]{
		Err: err,
	}
}
