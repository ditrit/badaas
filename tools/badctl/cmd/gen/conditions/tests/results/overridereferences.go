// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	overridereferences "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/overridereferences"
	gorm "gorm.io/gorm"
	"time"
)

func PhoneId(v badorm.UUID) badorm.WhereCondition[overridereferences.Phone] {
	return badorm.WhereCondition[overridereferences.Phone]{
		Field: "id",
		Value: v,
	}
}
func PhoneCreatedAt(v time.Time) badorm.WhereCondition[overridereferences.Phone] {
	return badorm.WhereCondition[overridereferences.Phone]{
		Field: "created_at",
		Value: v,
	}
}
func PhoneUpdatedAt(v time.Time) badorm.WhereCondition[overridereferences.Phone] {
	return badorm.WhereCondition[overridereferences.Phone]{
		Field: "updated_at",
		Value: v,
	}
}
func PhoneDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[overridereferences.Phone] {
	return badorm.WhereCondition[overridereferences.Phone]{
		Field: "deleted_at",
		Value: v,
	}
}
func PhoneBrand(conditions ...badorm.Condition[overridereferences.Brand]) badorm.Condition[overridereferences.Phone] {
	return badorm.JoinCondition[overridereferences.Phone, overridereferences.Brand]{
		Conditions: conditions,
		T1Field:    "brand_name",
		T2Field:    "name",
	}
}
func PhoneBrandName(v string) badorm.WhereCondition[overridereferences.Phone] {
	return badorm.WhereCondition[overridereferences.Phone]{
		Field: "brand_name",
		Value: v,
	}
}
