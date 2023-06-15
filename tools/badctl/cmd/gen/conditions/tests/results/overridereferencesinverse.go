// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overridereferencesinverse "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overridereferencesinverse"
	gorm "gorm.io/gorm"
	"time"
)

func ComputerId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[overridereferencesinverse.Computer, badorm.UUID] {
	return badorm.FieldCondition[overridereferencesinverse.Computer, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func ComputerCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[overridereferencesinverse.Computer, time.Time] {
	return badorm.FieldCondition[overridereferencesinverse.Computer, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func ComputerUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[overridereferencesinverse.Computer, time.Time] {
	return badorm.FieldCondition[overridereferencesinverse.Computer, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func ComputerDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[overridereferencesinverse.Computer, gorm.DeletedAt] {
	return badorm.FieldCondition[overridereferencesinverse.Computer, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func ComputerName(exprs ...badorm.Expression[string]) badorm.FieldCondition[overridereferencesinverse.Computer, string] {
	return badorm.FieldCondition[overridereferencesinverse.Computer, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
func ComputerProcessor(conditions ...badorm.Condition[overridereferencesinverse.Processor]) badorm.Condition[overridereferencesinverse.Computer] {
	return badorm.JoinCondition[overridereferencesinverse.Computer, overridereferencesinverse.Processor]{
		Conditions: conditions,
		T1Field:    "Name",
		T2Field:    "ComputerName",
	}
}
func ProcessorComputer(conditions ...badorm.Condition[overridereferencesinverse.Computer]) badorm.Condition[overridereferencesinverse.Processor] {
	return badorm.JoinCondition[overridereferencesinverse.Processor, overridereferencesinverse.Computer]{
		Conditions: conditions,
		T1Field:    "ComputerName",
		T2Field:    "Name",
	}
}