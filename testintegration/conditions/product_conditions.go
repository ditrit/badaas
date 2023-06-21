// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	"database/sql"
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func ProductId(expr badorm.Expression[badorm.UUID]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, badorm.UUID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func ProductCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func ProductUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func ProductDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var productStringFieldID = badorm.FieldIdentifier{Column: "string_something_else"}

func ProductString(expr badorm.Expression[string]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, string]{
		Expression:      expr,
		FieldIdentifier: productStringFieldID,
	}
}

var productIntFieldID = badorm.FieldIdentifier{Field: "Int"}

func ProductInt(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: productIntFieldID,
	}
}

var productIntPointerFieldID = badorm.FieldIdentifier{Field: "IntPointer"}

func ProductIntPointer(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: productIntPointerFieldID,
	}
}

var productFloatFieldID = badorm.FieldIdentifier{Field: "Float"}

func ProductFloat(expr badorm.Expression[float64]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, float64]{
		Expression:      expr,
		FieldIdentifier: productFloatFieldID,
	}
}

var productNullFloatFieldID = badorm.FieldIdentifier{Field: "NullFloat"}

func ProductNullFloat(expr badorm.Expression[sql.NullFloat64]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, sql.NullFloat64]{
		Expression:      expr,
		FieldIdentifier: productNullFloatFieldID,
	}
}

var productBoolFieldID = badorm.FieldIdentifier{Field: "Bool"}

func ProductBool(expr badorm.Expression[bool]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, bool]{
		Expression:      expr,
		FieldIdentifier: productBoolFieldID,
	}
}

var productNullBoolFieldID = badorm.FieldIdentifier{Field: "NullBool"}

func ProductNullBool(expr badorm.Expression[sql.NullBool]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, sql.NullBool]{
		Expression:      expr,
		FieldIdentifier: productNullBoolFieldID,
	}
}

var productByteArrayFieldID = badorm.FieldIdentifier{Field: "ByteArray"}

func ProductByteArray(expr badorm.Expression[[]uint8]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, []uint8]{
		Expression:      expr,
		FieldIdentifier: productByteArrayFieldID,
	}
}

var productMultiStringFieldID = badorm.FieldIdentifier{Field: "MultiString"}

func ProductMultiString(expr badorm.Expression[models.MultiString]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, models.MultiString]{
		Expression:      expr,
		FieldIdentifier: productMultiStringFieldID,
	}
}

var productToBeEmbeddedEmbeddedIntFieldID = badorm.FieldIdentifier{Field: "EmbeddedInt"}

func ProductToBeEmbeddedEmbeddedInt(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: productToBeEmbeddedEmbeddedIntFieldID,
	}
}

var productGormEmbeddedIntFieldID = badorm.FieldIdentifier{
	ColumnPrefix: "gorm_embedded_",
	Field:        "Int",
}

func ProductGormEmbeddedInt(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: productGormEmbeddedIntFieldID,
	}
}

var ProductPreload = badorm.NewPreloadCondition[models.Product](productStringFieldID, productIntFieldID, productIntPointerFieldID, productFloatFieldID, productNullFloatFieldID, productBoolFieldID, productNullBoolFieldID, productByteArrayFieldID, productMultiStringFieldID, productToBeEmbeddedEmbeddedIntFieldID, productGormEmbeddedIntFieldID)
