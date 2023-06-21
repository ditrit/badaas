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
func ProductString(expr badorm.Expression[string]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, string]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Column: "string_something_else"},
	}
}
func ProductInt(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "Int"},
	}
}
func ProductIntPointer(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "IntPointer"},
	}
}
func ProductFloat(expr badorm.Expression[float64]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, float64]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "Float"},
	}
}
func ProductNullFloat(expr badorm.Expression[sql.NullFloat64]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, sql.NullFloat64]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "NullFloat"},
	}
}
func ProductBool(expr badorm.Expression[bool]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, bool]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "Bool"},
	}
}
func ProductNullBool(expr badorm.Expression[sql.NullBool]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, sql.NullBool]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "NullBool"},
	}
}
func ProductByteArray(expr badorm.Expression[[]uint8]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, []uint8]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "ByteArray"},
	}
}
func ProductMultiString(expr badorm.Expression[models.MultiString]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, models.MultiString]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "MultiString"},
	}
}
func ProductEmbeddedInt(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression:      expr,
		FieldIdentifier: badorm.FieldIdentifier{Field: "EmbeddedInt"},
	}
}
func ProductGormEmbeddedInt(expr badorm.Expression[int]) badorm.WhereCondition[models.Product] {
	return badorm.FieldCondition[models.Product, int]{
		Expression: expr,
		FieldIdentifier: badorm.FieldIdentifier{
			ColumnPrefix: "gorm_embedded_",
			Field:        "Int",
		},
	}
}
