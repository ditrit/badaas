// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	gormembedded "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/gormembedded"
	gorm "gorm.io/gorm"
	"time"
)

func GormEmbeddedId(v uint) badorm.WhereCondition[gormembedded.GormEmbedded] {
	return badorm.WhereCondition[gormembedded.GormEmbedded]{
		Field: "id",
		Value: v,
	}
}
func GormEmbeddedCreatedAt(v time.Time) badorm.WhereCondition[gormembedded.GormEmbedded] {
	return badorm.WhereCondition[gormembedded.GormEmbedded]{
		Field: "created_at",
		Value: v,
	}
}
func GormEmbeddedUpdatedAt(v time.Time) badorm.WhereCondition[gormembedded.GormEmbedded] {
	return badorm.WhereCondition[gormembedded.GormEmbedded]{
		Field: "updated_at",
		Value: v,
	}
}
func GormEmbeddedDeletedAt(v gorm.DeletedAt) badorm.WhereCondition[gormembedded.GormEmbedded] {
	return badorm.WhereCondition[gormembedded.GormEmbedded]{
		Field: "deleted_at",
		Value: v,
	}
}
func GormEmbeddedGormEmbeddedInt(v int) badorm.WhereCondition[gormembedded.GormEmbedded] {
	return badorm.WhereCondition[gormembedded.GormEmbedded]{
		Field: "gorm_embedded_int",
		Value: v,
	}
}
