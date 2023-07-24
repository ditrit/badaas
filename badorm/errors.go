package badorm

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/badorm/sql"
)

// operators

var (
	ErrValueCantBeNull        = errors.New("value to be compared can't be null")
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query (not joined)")
	ErrJoinMustBeSelected     = errors.New("field's model is joined more than once, select which one you want to use with SelectJoin")
	ErrNotRelated             = errors.New("value type not related with T")
)

func OperatorError(err error, sqlOperator sql.Operator) error {
	return fmt.Errorf("%w; operator: %s", err, sqlOperator.Name())
}

func operatorNameError(err error, operatorName string) error {
	return fmt.Errorf("%w; operator: %s", err, operatorName)
}

func notRelatedError[T any](value any, operatorName string) error {
	return operatorNameError(fmt.Errorf("%w; type: %T, T: %T",
		ErrNotRelated,
		value,
		*new(T),
	), operatorName)
}

func fieldModelNotConcernedError(field iFieldIdentifier, sqlOperator sql.Operator) error {
	return OperatorError(fmt.Errorf("%w; not concerned model: %s",
		ErrFieldModelNotConcerned,
		field.GetModelType(),
	), sqlOperator)
}

func joinMustBeSelectedError(field iFieldIdentifier, sqlOperator sql.Operator) error {
	return OperatorError(fmt.Errorf("%w; joined multiple times model: %s",
		ErrJoinMustBeSelected,
		field.GetModelType(),
	), sqlOperator)
}

// conditions

var (
	ErrEmptyConditions     = errors.New("condition must have at least one inner condition")
	ErrOnlyPreloadsAllowed = errors.New("only conditions that do a preload are allowed")
)

func conditionOperatorError[TObject Model, TAtribute any](operatorErr error, condition FieldCondition[TObject, TAtribute]) error {
	return fmt.Errorf(
		"%w; model: %s, field: %s",
		operatorErr,
		condition.FieldIdentifier.ModelType.String(),
		condition.FieldIdentifier.Field,
	)
}

func emptyConditionsError[T Model](connector sql.Connector) error {
	return fmt.Errorf(
		"%w; connector: %s; model: %T",
		ErrEmptyConditions,
		connector.Name(),
		*new(T),
	)
}

func onlyPreloadsAllowedError[T Model](fieldName string) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		ErrOnlyPreloadsAllowed,
		*new(T),
		fieldName,
	)
}
