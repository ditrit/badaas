// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func EmployeeId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Employee] {
	return badorm.FieldCondition[models.Employee, badorm.UUID]{
		Expression: expr,
		Field:      "ID",
	}
}
func EmployeeCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Employee] {
	return badorm.FieldCondition[models.Employee, time.Time]{
		Expression: expr,
		Field:      "CreatedAt",
	}
}
func EmployeeUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Employee] {
	return badorm.FieldCondition[models.Employee, time.Time]{
		Expression: expr,
		Field:      "UpdatedAt",
	}
}
func EmployeeDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Employee] {
	return badorm.FieldCondition[models.Employee, gorm.DeletedAt]{
		Expression: expr,
		Field:      "DeletedAt",
	}
}
func EmployeeName(expr badorm.Expression[string]) badorm.WhereCondition[models.Employee] {
	return badorm.FieldCondition[models.Employee, string]{
		Expression: expr,
		Field:      "Name",
	}
}
func EmployeeBoss(conditions ...badorm.Condition[models.Employee]) badorm.Condition[models.Employee] {
	return badorm.JoinCondition[models.Employee, models.Employee]{
		Conditions: conditions,
		T1Field:    "BossID",
		T2Field:    "ID",
	}
}
func EmployeeBossId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Employee] {
	return badorm.FieldCondition[models.Employee, badorm.UUID]{
		Expression: expr,
		Field:      "BossID",
	}
}
