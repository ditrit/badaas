// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overrideforeignkeyinverse "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overrideforeignkeyinverse"
	gorm "gorm.io/gorm"
	"time"
)

func CreditCardId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[overrideforeignkeyinverse.CreditCard] {
	return badorm.FieldCondition[overrideforeignkeyinverse.CreditCard, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func CreditCardCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[overrideforeignkeyinverse.CreditCard] {
	return badorm.FieldCondition[overrideforeignkeyinverse.CreditCard, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func CreditCardUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[overrideforeignkeyinverse.CreditCard] {
	return badorm.FieldCondition[overrideforeignkeyinverse.CreditCard, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func CreditCardDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[overrideforeignkeyinverse.CreditCard] {
	return badorm.FieldCondition[overrideforeignkeyinverse.CreditCard, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var creditCardUserReferenceFieldID = badorm.FieldIdentifier{Field: "UserReference"}

func CreditCardUserReference(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[overrideforeignkeyinverse.CreditCard] {
	return badorm.FieldCondition[overrideforeignkeyinverse.CreditCard, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: creditCardUserReferenceFieldID,
	}
}

var CreditCardPreloadAttributes = badorm.NewPreloadCondition[overrideforeignkeyinverse.CreditCard](creditCardUserReferenceFieldID)