package unsafe

import (
	"fmt"

	"github.com/ditrit/badaas/badorm"
)

// Condition that can be used to express conditions that are not supported (yet?) by BaDORM
// Example: table1.columnX = table2.columnY
type Condition[T badorm.Model] struct {
	SQLCondition string
	Values       []any
}

//nolint:unused // see inside
func (condition Condition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition Condition[T]) ApplyTo(query *badorm.Query, table badorm.Table) error {
	return badorm.ApplyWhereCondition[T](condition, query, table)
}

func (condition Condition[T]) GetSQL(_ *badorm.Query, table badorm.Table) (string, []any, error) {
	return fmt.Sprintf(
		condition.SQLCondition,
		table.Alias,
	), condition.Values, nil
}

//nolint:unused // is used
func (condition Condition[T]) AffectsDeletedAt() bool {
	return false
}

// Condition that can be used to express conditions that are not supported (yet?) by BaDORM
// Example: table1.columnX = table2.columnY
func NewCondition[T badorm.Model](condition string, values ...any) Condition[T] {
	return Condition[T]{
		SQLCondition: condition,
		Values:       values,
	}
}
