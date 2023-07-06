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
	JoinNumber   int
}

func (expr MultivalueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (expr *MultivalueOperator[T]) SelectJoin(joinNumber uint) DynamicOperator[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

func (expr MultivalueOperator[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	placeholderList := []string{}
	values := []any{}

	for _, value := range expr.Values {
		field, isField := value.(IFieldIdentifier)
		if isField {
			// if it is a field, add the field column to the query
			modelTable, err := getModelTable(query, field, expr.JoinNumber)
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

	placeholders := strings.Join(placeholderList, " "+expr.SQLConnector.String()+" ")

	return fmt.Sprintf(
		"%s %s %s"+placeholders+"%s",
		columnName,
		expr.SQLOperator,
		expr.SQLPrefix,
		expr.SQLSuffix,
	), values, nil
}

func getModelTable(query *Query, field IFieldIdentifier, joinNumber int) (Table, error) {
	modelTables := query.GetTables(field.GetModelType())
	if modelTables == nil {
		return Table{}, ErrFieldModelNotConcerned
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	if joinNumber == UndefinedJoinNumber {
		return Table{}, ErrJoinMustBeSelected
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
