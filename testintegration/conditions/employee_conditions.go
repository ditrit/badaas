// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	"time"
)

func EmployeeId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Employee] {
	return orm.FieldCondition[models.Employee, orm.UUID]{
		FieldIdentifier: orm.IDFieldID,
		Operator:        operator,
	}
}
func EmployeeCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Employee] {
	return orm.FieldCondition[models.Employee, time.Time]{
		FieldIdentifier: orm.CreatedAtFieldID,
		Operator:        operator,
	}
}
func EmployeeUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Employee] {
	return orm.FieldCondition[models.Employee, time.Time]{
		FieldIdentifier: orm.UpdatedAtFieldID,
		Operator:        operator,
	}
}
func EmployeeDeletedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Employee] {
	return orm.FieldCondition[models.Employee, time.Time]{
		FieldIdentifier: orm.DeletedAtFieldID,
		Operator:        operator,
	}
}

var employeeNameFieldID = orm.FieldIdentifier{Field: "Name"}

func EmployeeName(operator orm.Operator[string]) orm.WhereCondition[models.Employee] {
	return orm.FieldCondition[models.Employee, string]{
		FieldIdentifier: employeeNameFieldID,
		Operator:        operator,
	}
}
func EmployeeBoss(conditions ...orm.Condition[models.Employee]) orm.Condition[models.Employee] {
	return orm.JoinCondition[models.Employee, models.Employee]{
		Conditions:    conditions,
		RelationField: "Boss",
		T1Field:       "BossID",
		T2Field:       "ID",
	}
}

var EmployeePreloadBoss = EmployeeBoss(EmployeePreloadAttributes)
var employeeBossIdFieldID = orm.FieldIdentifier{Field: "BossID"}

func EmployeeBossId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Employee] {
	return orm.FieldCondition[models.Employee, orm.UUID]{
		FieldIdentifier: employeeBossIdFieldID,
		Operator:        operator,
	}
}

var EmployeePreloadAttributes = orm.NewPreloadCondition[models.Employee](employeeNameFieldID, employeeBossIdFieldID)
var EmployeePreloadRelations = []orm.Condition[models.Employee]{EmployeePreloadBoss}
