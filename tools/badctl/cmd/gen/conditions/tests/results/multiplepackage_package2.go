// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	package2 "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/multiplepackage/package2"
	gorm "gorm.io/gorm"
	"time"
)

func Package2Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[package2.Package2] {
	return badorm.FieldCondition[package2.Package2, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func Package2CreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[package2.Package2] {
	return badorm.FieldCondition[package2.Package2, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func Package2UpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[package2.Package2] {
	return badorm.FieldCondition[package2.Package2, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func Package2DeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[package2.Package2] {
	return badorm.FieldCondition[package2.Package2, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var package2Package1IdFieldID = badorm.FieldIdentifier{Field: "Package1ID"}

func Package2Package1Id(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[package2.Package2] {
	return badorm.FieldCondition[package2.Package2, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: package2Package1IdFieldID,
	}
}

var Package2PreloadAttributes = badorm.NewPreloadCondition[package2.Package2](package2Package1IdFieldID)
var Package2PreloadRelations = []badorm.Condition[package2.Package2]{}
