package badorm

import (
	"errors"
	"fmt"
	"reflect"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldTypeDoesNotMatch      = errors.New("field type does not match expression type")
	ErrFieldModelNotConcerned     = errors.New("field's model is not concerned by the query so it can't be used in a expression")
	ErrExpressionTypeNotSupported = errors.New("expression type not supported")
)

// TODO ver el nombre
// TODO quizas field identifier deberia llamarse solo field
func NewDynamicExpression[T any](expression func(value T) Expression[T], field FieldIdentifier) Expression[T] {
	tValue := *new(T)

	if field.Type != reflect.TypeOf(tValue).Kind() {
		return NewInvalidExpression[T](ErrFieldTypeDoesNotMatch)
	}

	staticExpression := expression(
		// TODO esto podria no pasar alguna validacion
		tValue,
	)

	if _, isInvalid := staticExpression.(InvalidExpression[T]); isInvalid {
		return staticExpression
	}

	if valueExpression, isValue := staticExpression.(ValueExpression[T]); isValue {
		// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
		// creo que si hay mas de uno solo se tendria que aplicar en el primero,
		return DynamicExpression[T]{
			SQLExpressions: []LiteralSQLExpression{
				{
					SQL:   valueExpression.ExpressionsAndValues[0].SQLExpression,
					Field: field,
				},
			},
		}
	}

	// TODO soportar multivalue, no todos necesariamente dinamicos
	return NewInvalidExpression[T](ErrExpressionTypeNotSupported)
}

// TODO doc
type DynamicExpression[T any] struct {
	// TODO hacer el cambio de nombre en el anterior tambien?
	SQLExpressions []LiteralSQLExpression
}

type LiteralSQLExpression struct {
	SQL   string
	Field FieldIdentifier
}

func (expr DynamicExpression[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Expression[T]
}

// verificar que en las condiciones anteriores alguien usó el field con el que se intenta comparar
// obtener de ahi cual es el nombre de la table a usar con ese field.
// TODO doc a ingles
func (expr DynamicExpression[T]) ToSQL(query *query, columnName string) (string, []any, error) {
	exprString := columnName
	values := []any{}

	for _, sqlExpr := range expr.SQLExpressions {
		modelTable, isConcerned := query.concernedModels[sqlExpr.Field.ModelType]
		if !isConcerned {
			return "", nil, ErrFieldModelNotConcerned
		}

		exprString += fmt.Sprintf(
			" "+sqlExpr.SQL+" %s",
			sqlExpr.Field.ColumnSQL(query, modelTable),
		)
	}

	return exprString, values, nil
}
