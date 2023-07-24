package badorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"

	badormSQL "github.com/ditrit/badaas/badorm/sql"
)

const undefinedJoinNumber = -1

// Operator that compares the value of the column against a fixed value
// If SQLOperators has multiple entries, comparisons will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	operations []operation
}

type operation struct {
	SQLOperator badormSQL.Operator
	Value       any
	JoinNumber  int
}

func (operator ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

// Allows to choose which number of join use
// for the operation in position "operationNumber"
// when the value is a field and its model is joined more than once.
// Does nothing if the operationNumber is bigger than the amount of operations.
func (operator *ValueOperator[T]) SelectJoin(operationNumber, joinNumber uint) DynamicOperator[T] {
	if operationNumber >= uint(len(operator.operations)) {
		return operator
	}

	operationSaved := operator.operations[operationNumber]
	operationSaved.JoinNumber = int(joinNumber)
	operator.operations[operationNumber] = operationSaved

	return operator
}

func (operator *ValueOperator[T]) AddOperation(sqlOperator badormSQL.Operator, value any) *ValueOperator[T] {
	operator.operations = append(
		operator.operations,
		operation{
			Value:       value,
			SQLOperator: sqlOperator,
			JoinNumber:  undefinedJoinNumber,
		},
	)

	return operator
}

func (operator ValueOperator[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	operationString := columnName
	values := []any{}

	// add each operation to the sql
	for _, operation := range operator.operations {
		field, isField := operation.Value.(iFieldIdentifier)
		if isField {
			// if the value of the operation is a field,
			// verify that this field is concerned by the query
			// (a join was performed with the model to which this field belongs)
			// and get the alias of the table of this model.
			modelTable, err := getModelTable(query, field, operation.JoinNumber, operation.SQLOperator)
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
	return *new(ValueOperator[T]).AddOperation(sqlOperator, value)
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
		return NewInvalidOperator[T](OperatorError(err, sqlOperator))
	}

	return NewValueOperator[T](sqlOperator, pattern)
}
