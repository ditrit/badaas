// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overrideforeignkeyinverse "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overrideforeignkeyinverse"
	gorm "gorm.io/gorm"
	"time"
)

func UserId(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[overrideforeignkeyinverse.User, badorm.UUID] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func UserCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[overrideforeignkeyinverse.User, time.Time] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func UserUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[overrideforeignkeyinverse.User, time.Time] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func UserDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[overrideforeignkeyinverse.User, gorm.DeletedAt] {
	return badorm.FieldCondition[overrideforeignkeyinverse.User, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func UserCreditCard(conditions ...badorm.Condition[overrideforeignkeyinverse.CreditCard]) badorm.Condition[overrideforeignkeyinverse.User] {
	return badorm.JoinCondition[overrideforeignkeyinverse.User, overrideforeignkeyinverse.CreditCard]{
		Conditions: conditions,
		T1Field:    "ID",
		T2Field:    "UserReference",
	}
}
func CreditCardUser(conditions ...badorm.Condition[overrideforeignkeyinverse.User]) badorm.Condition[overrideforeignkeyinverse.CreditCard] {
	return badorm.JoinCondition[overrideforeignkeyinverse.CreditCard, overrideforeignkeyinverse.User]{
		Conditions: conditions,
		T1Field:    "UserReference",
		T2Field:    "ID",
	}
}