// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func PersonId(v badorm.UUID) badorm.WhereCondition[models.Person] {
	return badorm.WhereCondition[models.Person]{
		Field: "id",
		Value: v,
	}
}
func PersonCreatedAt(v time.Time) badorm.WhereCondition[models.Person] {
	return badorm.WhereCondition[models.Person]{
		Field: "created_at",
		Value: v,
	}
}
func PersonUpdatedAt(v time.Time) badorm.WhereCondition[models.Person] {
	return badorm.WhereCondition[models.Person]{
		Field: "updated_at",
		Value: v,
	}
}
func PersonDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[models.Person] {
	return badorm.WhereCondition[models.Person]{
		Field: "deleted_at",
		Value: v,
	}
}
func PersonName(v string) badorm.WhereCondition[models.Person] {
	return badorm.WhereCondition[models.Person]{
		Field: "name",
		Value: v,
	}
}
