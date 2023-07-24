package badorm

import (
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm/sql"
)

// Operator that compares the value of the column against multiple values
// Example: value IN (v1, v2, v3, ..., vN)
type MultivalueOperator[T any] struct {
	Values       []any         // the values to compare with
	SQLOperator  sql.Operator  // the operator used to compare, example: IN
	SQLConnector sql.Connector // the connector between values, example: ', '
	SQLPrefix    string        // something to put before the values, example: (
	SQLSuffix    string        // something to put after the values, example: )
	JoinNumbers  map[uint]int  // join number to use in each value
}

func (operator MultivalueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

// Allows to choose which number of join use
// for the value in position "valueNumber"
// when the value is a field and its model is joined more than once.
// Does nothing if the valueNumber is bigger than the amount of values.
func (operator *MultivalueOperator[T]) SelectJoin(valueNumber, joinNumber uint) DynamicOperator[T] {
	joinNumbers := operator.JoinNumbers
	if joinNumbers == nil {
		joinNumbers = map[uint]int{}
	}

	joinNumbers[valueNumber] = int(joinNumber)
	operator.JoinNumbers = joinNumbers

	return operator
}

func (operator MultivalueOperator[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	placeholderList := []string{}
	values := []any{}

	for i, value := range operator.Values {
		field, isField := value.(iFieldIdentifier)
		if isField {
			joinNumber, isPresent := operator.JoinNumbers[uint(i)]
			if !isPresent {
				joinNumber = undefinedJoinNumber
			}

			// if it is a field, add the field column to the query
			modelTable, err := getModelTable(query, field, joinNumber, operator.SQLOperator)
			if err != nil {
				return "", nil, err
			}

			placeholderList = append(placeholderList, field.ColumnSQL(query, modelTable))
		} else {
			// if it is not a field, it a value, ass the placeholder ? and the value to the list
			placeholderList = append(placeholderList, "?")
			values = append(values, value)
		}
	}

	placeholders := strings.Join(placeholderList, " "+operator.SQLConnector.String()+" ")

	return fmt.Sprintf(
		"%s %s %s"+placeholders+"%s",
		columnName,
		operator.SQLOperator,
		operator.SQLPrefix,
		operator.SQLSuffix,
	), values, nil
}

func getModelTable(query *Query, field iFieldIdentifier, joinNumber int, sqlOperator sql.Operator) (Table, error) {
	modelTables := query.GetTables(field.GetModelType())
	if modelTables == nil {
		return Table{}, fieldModelNotConcernedError(field, sqlOperator)
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	if joinNumber == undefinedJoinNumber {
		return Table{}, joinMustBeSelectedError(field, sqlOperator)
	}

	return modelTables[joinNumber], nil
}

func NewMultivalueOperator[T any](
	sqlOperator sql.Operator,
	sqlConnector sql.Connector,
	sqlPrefix, sqlSuffix string,
	values ...T,
) Operator[T] {
	valuesAny := pie.Map(values, func(value T) any {
		return value
	})

	return &MultivalueOperator[T]{
		Values:       valuesAny,
		SQLOperator:  sqlOperator,
		SQLConnector: sqlConnector,
		SQLPrefix:    sqlPrefix,
		SQLSuffix:    sqlSuffix,
	}
}
