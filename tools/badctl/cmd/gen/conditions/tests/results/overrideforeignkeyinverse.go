// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overrideforeignkeyinverse "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overrideforeignkeyinverse"
	gorm "gorm.io/gorm"
	"time"
)

func UserId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[overrideforeignkeyinverse.User] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func UserCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[overrideforeignkeyinverse.User] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func UserUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[overrideforeignkeyinverse.User] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func UserDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[overrideforeignkeyinverse.User] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func UserCreditCard(conditions ...badorm.Condition[overrideforeignkeyinverse.CreditCard]) badorm.Condition[overrideforeignkeyinverse.User] {
	return badorm.JoinCondition[overrideforeignkeyinverse.User, overrideforeignkeyinverse.CreditCard]{
		Conditions:    conditions,
		RelationField: "CreditCard",
		T1Field:       "ID",
		T2Field:       "UserReference",
	}
}
func CreditCardUser(conditions ...badorm.Condition[overrideforeignkeyinverse.User]) badorm.Condition[overrideforeignkeyinverse.CreditCard] {
	return badorm.JoinCondition[overrideforeignkeyinverse.CreditCard, overrideforeignkeyinverse.User]{
		Conditions:    conditions,
		RelationField: "User",
		T1Field:       "UserReference",
		T2Field:       "ID",
	}
}

var UserPreload = badorm.NewPreloadCondition[overrideforeignkeyinverse.User]()
