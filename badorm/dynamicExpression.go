package badorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm/expressions"
)

// TODO ponerles mas informacion precisa a los errores
var (
	ErrFieldTypeDoesNotMatch      = errors.New("field type does not match expression type")
	ErrFieldModelNotConcerned     = errors.New("field's model is not concerned by the query so it can't be used in a expression")
	ErrExpressionTypeNotSupported = errors.New("expression type not supported")
)

var nullableTypes = []reflect.Type{
	reflect.TypeOf(sql.NullBool{}),
	reflect.TypeOf(sql.NullByte{}),
	reflect.TypeOf(sql.NullFloat64{}),
	reflect.TypeOf(sql.NullInt16{}),
	reflect.TypeOf(sql.NullInt32{}),
	reflect.TypeOf(sql.NullInt64{}),
	reflect.TypeOf(sql.NullString{}),
	reflect.TypeOf(sql.NullTime{}),
	reflect.TypeOf(gorm.DeletedAt{}),
}

func isNullable(fieldType reflect.Type) bool {
	return pie.Contains(nullableTypes, fieldType)
}

// TODO ver el nombre
// TODO quizas field identifier deberia llamarse solo field
func NewDynamicExpression[T any](expression expressions.SQLExpression, field FieldIdentifier) Expression[T] {
	expressionType := reflect.TypeOf(*new(T))
	fieldType := field.Type

	// TODO si el field lo paso a tipado esto podria andar tambien en tiempo de compilacion,
	// el problema es que no tengo como hacer esto de los nullables, a menos que cree otro metodo
	// o que a los fields de cosas nullables le ponga el T de no nullable
	// pero la expression me va a quedar de T no nullable y no va a andar
	// podria ser otro metodo que acepte un field de cualquier tipo y ahi ver las condiciones estas
	// el problema es que ese tipo despues no lo puedo obtener en tiempo de ejecucion creo
	// entonces la unica posibilidad seria una funcion para cada una de las posibilidades

	if fieldType != expressionType &&
		!((isNullable(fieldType) && fieldType.Field(0).Type == expressionType) ||
			(isNullable(expressionType) && expressionType.Field(0).Type == fieldType)) {
		return NewInvalidExpression[T](ErrFieldTypeDoesNotMatch)
	}

	// TODO que pasa con los que solo aceptan cierto tipo, ver las de like por ejemplo
	// TODO que pasa si hay mas de uno, no se si me gusta mucho esto
	// creo que si hay mas de uno solo se tendria que aplicar en el primero,
	return DynamicExpression[T]{
		SQLExpressions: []LiteralSQLExpression{
			{
				SQL:   expressions.ToSQL[expression],
				Field: field,
			},
		},
	}
	// TODO soportar multivalue, no todos necesariamente dinamicos
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
