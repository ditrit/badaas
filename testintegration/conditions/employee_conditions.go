// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func EmployeeId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[models.Employee, badorm.UUID] {
	return badorm.FieldCondition[models.Employee, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func EmployeeCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[models.Employee, time.Time] {
	return badorm.FieldCondition[models.Employee, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func EmployeeUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[models.Employee, time.Time] {
	return badorm.FieldCondition[models.Employee, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func EmployeeDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[models.Employee, gorm.DeletedAt] {
	return badorm.FieldCondition[models.Employee, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func EmployeeName(exprs ...badorm.Expression[string]) badorm.FieldCondition[models.Employee, string] {
	return badorm.FieldCondition[models.Employee, string]{
		Expressions: exprs,
		Field:       "Name",
	}
}
func EmployeeBoss(conditions ...badorm.Condition[models.Employee]) badorm.Condition[models.Employee] {
	return badorm.JoinCondition[models.Employee, models.Employee]{
		Conditions: conditions,
		T1Field:    "BossID",
		T2Field:    "ID",
	}
}
func EmployeeBossId(exprs ...badorm.Expression[*badorm.UUID]) badorm.FieldCondition[models.Employee, *badorm.UUID] {
	return badorm.FieldCondition[models.Employee, *badorm.UUID]{
		Expressions: exprs,
		Field:       "BossID",
	}
}