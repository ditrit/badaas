package badorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

const UndefinedJoinNumber = -1

// Operator that compares the value of the column against a fixed value
// If SQLOperators has multiple entries, comparisons will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	Operations []Operation
	JoinNumber int
}

type Operation struct {
	SQLOperator badormSQL.Operator
	Value       any
}

func (operator ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (operator *ValueOperator[T]) SelectJoin(joinNumber uint) DynamicOperator[T] {
	operator.JoinNumber = int(joinNumber)
	return operator
}

func (operator *ValueOperator[T]) AddOperation(sqlOperator badormSQL.Operator, value any) ValueOperator[T] {
	operator.Operations = append(
		operator.Operations,
		Operation{
			Value:       value,
			SQLOperator: sqlOperator,
		},
	)

	return *operator
}

func (operator ValueOperator[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	operationString := columnName
	values := []any{}

	// add each operation to the sql
	for _, operation := range operator.Operations {
		field, isField := operation.Value.(iFieldIdentifier)
		if isField {
			// if the value of the operation is a field,
			// verify that this field is concerned by the query
			// (a join was performed with the model to which this field belongs)
			// and get the alias of the table of this model.
			modelTable, err := getModelTable(query, field, operator.JoinNumber, operation.SQLOperator)
			if err != nil {
				return "", nil, err
			}

			operationString += fmt.Sprintf(
				" %s %s",
				operation.SQLOperator,
				field.ColumnSQL(query, modelTable),
			)
		} else {
			operationString += " " + operation.SQLOperator.String() + " ?"
			values = append(values, operation.Value)
		}
	}

	return operationString, values, nil
}

func NewValueOperator[T any](sqlOperator badormSQL.Operator, value any) ValueOperator[T] {
	expr := &ValueOperator[T]{}

	return expr.AddOperation(sqlOperator, value)
}

var nullableKinds = []reflect.Kind{
	reflect.Chan, reflect.Func,
	reflect.Map, reflect.Pointer,
	reflect.UnsafePointer, reflect.Interface,
	reflect.Slice,
}

func NewCantBeNullValueOperator[T any](sqlOperator badormSQL.Operator, value any) Operator[T] {
	if value == nil || mapsToNull(value) {
		return NewInvalidOperator[T](
			OperatorError(ErrValueCantBeNull, sqlOperator),
		)
	}

	return NewValueOperator[T](sqlOperator, value)
}

func NewMustBePOSIXValueOperator[T string | sql.NullString](sqlOperator badormSQL.Operator, pattern string) Operator[T] {
	_, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return NewInvalidOperator[T](err)
	}

	return NewValueOperator[T](sqlOperator, pattern)
}
