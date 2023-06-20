// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package conditions

import (
	"database/sql"
	orm "github.com/ditrit/badaas/orm"
	models "github.com/ditrit/badaas/testintegration/models"
	gorm "gorm.io/gorm"
	"time"
)

func ProductId(operator orm.Operator[orm.UUID]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, orm.UUID]{
		Operator:      operator,
		FieldIdentifier: orm.IDFieldID,
	}
}
func ProductCreatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, time.Time]{
		Operator:      operator,
		FieldIdentifier: orm.CreatedAtFieldID,
	}
}
func ProductUpdatedAt(operator orm.Operator[time.Time]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, time.Time]{
		Operator:      operator,
		FieldIdentifier: orm.UpdatedAtFieldID,
	}
}
func ProductDeletedAt(operator orm.Operator[gorm.DeletedAt]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, gorm.DeletedAt]{
		Operator:      operator,
		FieldIdentifier: orm.DeletedAtFieldID,
	}
}

var productStringFieldID = orm.FieldIdentifier{Column: "string_something_else"}

func ProductString(operator orm.Operator[string]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, string]{
		Operator:      operator,
		FieldIdentifier: productStringFieldID,
	}
}

var productIntFieldID = orm.FieldIdentifier{Field: "Int"}

func ProductInt(operator orm.Operator[int]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, int]{
		Operator:      operator,
		FieldIdentifier: productIntFieldID,
	}
}

var productIntPointerFieldID = orm.FieldIdentifier{Field: "IntPointer"}

func ProductIntPointer(operator orm.Operator[int]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, int]{
		Operator:      operator,
		FieldIdentifier: productIntPointerFieldID,
	}
}

var productFloatFieldID = orm.FieldIdentifier{Field: "Float"}

func ProductFloat(operator orm.Operator[float64]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, float64]{
		Operator:      operator,
		FieldIdentifier: productFloatFieldID,
	}
}

var productNullFloatFieldID = orm.FieldIdentifier{Field: "NullFloat"}

func ProductNullFloat(operator orm.Operator[sql.NullFloat64]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, sql.NullFloat64]{
		Operator:      operator,
		FieldIdentifier: productNullFloatFieldID,
	}
}

var productBoolFieldID = orm.FieldIdentifier{Field: "Bool"}

func ProductBool(operator orm.Operator[bool]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, bool]{
		Operator:      operator,
		FieldIdentifier: productBoolFieldID,
	}
}

var productNullBoolFieldID = orm.FieldIdentifier{Field: "NullBool"}

func ProductNullBool(operator orm.Operator[sql.NullBool]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, sql.NullBool]{
		Operator:      operator,
		FieldIdentifier: productNullBoolFieldID,
	}
}

var productByteArrayFieldID = orm.FieldIdentifier{Field: "ByteArray"}

func ProductByteArray(operator orm.Operator[[]uint8]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, []uint8]{
		Operator:      operator,
		FieldIdentifier: productByteArrayFieldID,
	}
}

var productMultiStringFieldID = orm.FieldIdentifier{Field: "MultiString"}

func ProductMultiString(operator orm.Operator[models.MultiString]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, models.MultiString]{
		Operator:      operator,
		FieldIdentifier: productMultiStringFieldID,
	}
}

var productToBeEmbeddedEmbeddedIntFieldID = orm.FieldIdentifier{Field: "EmbeddedInt"}

func ProductToBeEmbeddedEmbeddedInt(operator orm.Operator[int]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, int]{
		Operator:      operator,
		FieldIdentifier: productToBeEmbeddedEmbeddedIntFieldID,
	}
}

var productGormEmbeddedIntFieldID = orm.FieldIdentifier{
	ColumnPrefix: "gorm_embedded_",
	Field:        "Int",
}

func ProductGormEmbeddedInt(operator orm.Operator[int]) orm.WhereCondition[models.Product] {
	return orm.FieldCondition[models.Product, int]{
		Operator:      operator,
		FieldIdentifier: productGormEmbeddedIntFieldID,
	}
}

var ProductPreloadAttributes = orm.NewPreloadCondition[models.Product](productStringFieldID, productIntFieldID, productIntPointerFieldID, productFloatFieldID, productNullFloatFieldID, productBoolFieldID, productNullBoolFieldID, productByteArrayFieldID, productMultiStringFieldID, productToBeEmbeddedEmbeddedIntFieldID, productGormEmbeddedIntFieldID)
