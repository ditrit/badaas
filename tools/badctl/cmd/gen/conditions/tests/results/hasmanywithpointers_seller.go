// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	hasmanywithpointers "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/hasmanywithpointers"
	gorm "gorm.io/gorm"
	"time"
)

func SellerInPointersId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasmanywithpointers.SellerInPointers] {
	return badorm.FieldCondition[hasmanywithpointers.SellerInPointers, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func SellerInPointersCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasmanywithpointers.SellerInPointers] {
	return badorm.FieldCondition[hasmanywithpointers.SellerInPointers, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func SellerInPointersUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[hasmanywithpointers.SellerInPointers] {
	return badorm.FieldCondition[hasmanywithpointers.SellerInPointers, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func SellerInPointersDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[hasmanywithpointers.SellerInPointers] {
	return badorm.FieldCondition[hasmanywithpointers.SellerInPointers, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}
func SellerInPointersCompany(conditions ...badorm.Condition[hasmanywithpointers.CompanyWithPointers]) badorm.IJoinCondition[hasmanywithpointers.SellerInPointers] {
	return badorm.JoinCondition[hasmanywithpointers.SellerInPointers, hasmanywithpointers.CompanyWithPointers]{
		Conditions:         conditions,
		RelationField:      "Company",
		T1Field:            "CompanyID",
		T1PreloadCondition: SellerInPointersPreloadAttributes,
		T2Field:            "ID",
	}
}

var SellerInPointersPreloadCompany = SellerInPointersCompany(CompanyWithPointersPreloadAttributes)
var sellerInPointersCompanyIdFieldID = badorm.FieldIdentifier{Field: "CompanyID"}

func SellerInPointersCompanyId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[hasmanywithpointers.SellerInPointers] {
	return badorm.FieldCondition[hasmanywithpointers.SellerInPointers, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: sellerInPointersCompanyIdFieldID,
	}
}

var SellerInPointersPreloadAttributes = badorm.NewPreloadCondition[hasmanywithpointers.SellerInPointers](sellerInPointersCompanyIdFieldID)
var SellerInPointersPreloadRelations = []badorm.Condition[hasmanywithpointers.SellerInPointers]{SellerInPointersPreloadCompany}