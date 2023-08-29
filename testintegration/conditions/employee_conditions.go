// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	orm "github.com/ditrit/badaas/orm"
	condition "github.com/ditrit/badaas/orm/condition"
	model "github.com/ditrit/badaas/orm/model"
	query "github.com/ditrit/badaas/orm/query"
	models "github.com/ditrit/badaas/testintegration/models"
	"reflect"
	"time"
)

var employeeType = reflect.TypeOf(*new(models.Employee))

func (employeeConditions employeeConditions) IdIs() orm.FieldIs[models.Employee, model.UUID] {
	return orm.FieldIs[models.Employee, model.UUID]{FieldID: employeeConditions.ID}
}
func (employeeConditions employeeConditions) CreatedAtIs() orm.FieldIs[models.Employee, time.Time] {
	return orm.FieldIs[models.Employee, time.Time]{FieldID: employeeConditions.CreatedAt}
}
func (employeeConditions employeeConditions) UpdatedAtIs() orm.FieldIs[models.Employee, time.Time] {
	return orm.FieldIs[models.Employee, time.Time]{FieldID: employeeConditions.UpdatedAt}
}
func (employeeConditions employeeConditions) DeletedAtIs() orm.FieldIs[models.Employee, time.Time] {
	return orm.FieldIs[models.Employee, time.Time]{FieldID: employeeConditions.DeletedAt}
}
func (employeeConditions employeeConditions) NameIs() orm.StringFieldIs[models.Employee] {
	return orm.StringFieldIs[models.Employee]{FieldIs: orm.FieldIs[models.Employee, string]{FieldID: employeeConditions.Name}}
}
func (employeeConditions employeeConditions) Boss(conditions ...condition.Condition[models.Employee]) condition.JoinCondition[models.Employee] {
	return condition.NewJoinCondition[models.Employee, models.Employee](conditions, "Boss", "BossID", employeeConditions.Preload(), "ID")
}
func (employeeConditions employeeConditions) PreloadBoss() condition.JoinCondition[models.Employee] {
	return employeeConditions.Boss(Employee.Preload())
}
func (employeeConditions employeeConditions) BossIdIs() orm.FieldIs[models.Employee, model.UUID] {
	return orm.FieldIs[models.Employee, model.UUID]{FieldID: employeeConditions.BossID}
}

type employeeConditions struct {
	ID        query.FieldIdentifier[model.UUID]
	CreatedAt query.FieldIdentifier[time.Time]
	UpdatedAt query.FieldIdentifier[time.Time]
	DeletedAt query.FieldIdentifier[time.Time]
	Name      query.FieldIdentifier[string]
	BossID    query.FieldIdentifier[model.UUID]
}

var Employee = employeeConditions{
	BossID: query.FieldIdentifier[model.UUID]{
		Field:     "BossID",
		ModelType: employeeType,
	},
	CreatedAt: query.FieldIdentifier[time.Time]{
		Field:     "CreatedAt",
		ModelType: employeeType,
	},
	DeletedAt: query.FieldIdentifier[time.Time]{
		Field:     "DeletedAt",
		ModelType: employeeType,
	},
	ID: query.FieldIdentifier[model.UUID]{
		Field:     "ID",
		ModelType: employeeType,
	},
	Name: query.FieldIdentifier[string]{
		Field:     "Name",
		ModelType: employeeType,
	},
	UpdatedAt: query.FieldIdentifier[time.Time]{
		Field:     "UpdatedAt",
		ModelType: employeeType,
	},
}

// Preload allows preloading the Employee when doing a query
func (employeeConditions employeeConditions) Preload() condition.Condition[models.Employee] {
	return condition.NewPreloadCondition[models.Employee](employeeConditions.ID, employeeConditions.CreatedAt, employeeConditions.UpdatedAt, employeeConditions.DeletedAt, employeeConditions.Name, employeeConditions.BossID)
}

// PreloadRelations allows preloading all the Employee's relation when doing a query
func (employeeConditions employeeConditions) PreloadRelations() []condition.Condition[models.Employee] {
	return []condition.Condition[models.Employee]{employeeConditions.PreloadBoss()}
}
