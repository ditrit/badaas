// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func PersonId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Person] {
	return badorm.FieldCondition[models.Person, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func PersonCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Person] {
	return badorm.FieldCondition[models.Person, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func PersonUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Person] {
	return badorm.FieldCondition[models.Person, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func PersonDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Person] {
	return badorm.FieldCondition[models.Person, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var personNameFieldID = badorm.FieldIdentifier{Field: "Name"}

func PersonName(expr badorm.Expression[string]) badorm.WhereCondition[models.Person] {
	return badorm.FieldCondition[models.Person, string]{
		Expression:      expr,
		FieldIdentifier: personNameFieldID,
	}
}

var PersonPreloadAttributes = badorm.NewPreloadCondition[models.Person](personNameFieldID)
