package badorm

import (
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/badorm/expressions"
)

// Expression that compares the value of the column against multiple values
// Example: value IN (v1, v2, v3, ..., vN)
type MultivalueExpression[T any] struct {
	// TODO hacer el cambio de nombre en el anterior tambien?
	// TODO con esto podria reemplazar el SQLExpressionAndValue para que todos sean por dentro dynamics
	Values        []any  // the values to compare with
	SQLExpression string // the expression used to compare, example: IN
	SQLConnector  string // the connector between values, example: ', '
	SQLPrefix     string // something to put before the values, example: (
	SQLSuffix     string // something to put after the values, example: )
	JoinNumber    int
}

func (expr MultivalueExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

func (expr *MultivalueExpression[T]) SelectJoin(joinNumber uint) DynamicExpression[T] {
	expr.JoinNumber = int(joinNumber)
	return expr
}

func (expr MultivalueExpression[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
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

	placeholders := strings.Join(placeholderList, " "+expr.SQLConnector+" ")

	return fmt.Sprintf(
		"%s %s %s"+placeholders+"%s",
		columnName,
		expr.SQLExpression,
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

func NewMultivalueExpression[T any](sqlExpression expressions.SQLExpression, sqlConnector, sqlPrefix, sqlSuffix string, values ...T) Expression[T] {
	valuesAny := pie.Map(values, func(value T) any {
		return value
	})

	return &MultivalueExpression[T]{
		Values:        valuesAny,
		SQLExpression: expressions.ToSQL[sqlExpression],
		SQLConnector:  sqlConnector,
		SQLPrefix:     sqlPrefix,
		SQLSuffix:     sqlSuffix,
	}
}
