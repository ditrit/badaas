// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	columndefinition "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/columndefinition"
	gorm "gorm.io/gorm"
	"time"
)

func ColumnDefinitionId(v badorm.UUID) badorm.WhereCondition[columndefinition.ColumnDefinition] {
	return badorm.WhereCondition[columndefinition.ColumnDefinition]{
		Field: "ID",
		Value: v,
	}
}
func ColumnDefinitionCreatedAt(v time.Time) badorm.WhereCondition[columndefinition.ColumnDefinition] {
	return badorm.WhereCondition[columndefinition.ColumnDefinition]{
		Field: "CreatedAt",
		Value: v,
	}
}
func ColumnDefinitionUpdatedAt(v time.Time) badorm.WhereCondition[columndefinition.ColumnDefinition] {
	return badorm.WhereCondition[columndefinition.ColumnDefinition]{
		Field: "UpdatedAt",
		Value: v,
	}
}
func ColumnDefinitionDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[columndefinition.ColumnDefinition] {
	return badorm.WhereCondition[columndefinition.ColumnDefinition]{
		Field: "DeletedAt",
		Value: v,
	}
}
func ColumnDefinitionString(v string) badorm.WhereCondition[columndefinition.ColumnDefinition] {
	return badorm.WhereCondition[columndefinition.ColumnDefinition]{
		Column: "string_something_else",
		Value:  v,
	}
}
