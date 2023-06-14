// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	goembedded "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/goembedded"
	gorm "gorm.io/gorm"
	"time"
)

func GoEmbeddedId(exprs ...badorm.Expression[uint]) badorm.FieldCondition[goembedded.GoEmbedded, uint] {
	return badorm.FieldCondition[goembedded.GoEmbedded, uint]{
		Expressions: exprs,
		Field:       "ID",
	}
}
func GoEmbeddedCreatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[goembedded.GoEmbedded, time.Time] {
	return badorm.FieldCondition[goembedded.GoEmbedded, time.Time]{
		Expressions: exprs,
		Field:       "CreatedAt",
	}
}
func GoEmbeddedUpdatedAt(exprs ...badorm.Expression[time.Time]) badorm.FieldCondition[goembedded.GoEmbedded, time.Time] {
	return badorm.FieldCondition[goembedded.GoEmbedded, time.Time]{
		Expressions: exprs,
		Field:       "UpdatedAt",
	}
}
func GoEmbeddedDeletedAt(exprs ...badorm.Expression[gorm.DeletedAt]) badorm.FieldCondition[goembedded.GoEmbedded, gorm.DeletedAt] {
	return badorm.FieldCondition[goembedded.GoEmbedded, gorm.DeletedAt]{
		Expressions: exprs,
		Field:       "DeletedAt",
	}
}
func GoEmbeddedEmbeddedInt(exprs ...badorm.Expression[int]) badorm.FieldCondition[goembedded.GoEmbedded, int] {
	return badorm.FieldCondition[goembedded.GoEmbedded, int]{
		Expressions: exprs,
		Field:       "EmbeddedInt",
	}
}
