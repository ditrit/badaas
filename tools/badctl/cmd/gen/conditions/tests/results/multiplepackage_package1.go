// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	package1 "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/multiplepackage/package1"
	package2 "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/multiplepackage/package2"
	gorm "gorm.io/gorm"
	"time"
)

func Package1Id(exprs ...badorm.Expression[badorm.UUID]) badorm.FieldCondition[package1.Package1, badorm.UUID] {
	return badorm.FieldCondition[package1.Package1, badorm.UUID]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func Package1CreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[package1.Package1, time.Time] {
	return badorm.FieldCondition[package1.Package1, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func Package1UpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[package1.Package1, time.Time] {
	return badorm.FieldCondition[package1.Package1, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func Package1DeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[package1.Package1, gorm.DeletedAt] {
	return badorm.FieldCondition[package1.Package1, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func Package1Package2(conditions ...badorm.Condition[package2.Package2]) badorm.Condition[package1.Package1] {
	return badorm.JoinCondition[package1.Package1, package2.Package2]{
		Conditions: conditions,
		T1Field:    "ID",
		T2Field:    "Package1ID",
	}
}
func Package2Package1(conditions ...badorm.Condition[package1.Package1]) badorm.Condition[package2.Package2] {
	return badorm.JoinCondition[package2.Package2, package1.Package1]{
		Conditions: conditions,
		T1Field:    "Package1ID",
		T2Field:    "ID",
	}
}
