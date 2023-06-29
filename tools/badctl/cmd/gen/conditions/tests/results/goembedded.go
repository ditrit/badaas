// Code generated by badctl v0.0.0, DO NOT EDIT.
package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	goembedded "github.com/ditrit/badaas/tools/badctl/cmd/gen/conditions/tests/goembedded"
	gorm "gorm.io/gorm"
	"time"
)

func GoEmbeddedId(expr badorm.Expression[badorm.UIntID]) badorm.WhereCondition[goembedded.GoEmbedded] {
	return badorm.FieldCondition[goembedded.GoEmbedded, badorm.UIntID]{
		Expression:      expr,
		FieldIdentifier: badorm.IDFieldID,
	}
}
func GoEmbeddedCreatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[goembedded.GoEmbedded] {
	return badorm.FieldCondition[goembedded.GoEmbedded, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.CreatedAtFieldID,
	}
}
func GoEmbeddedUpdatedAt(expr badorm.Expression[time.Time]) badorm.WhereCondition[goembedded.GoEmbedded] {
	return badorm.FieldCondition[goembedded.GoEmbedded, time.Time]{
		Expression:      expr,
		FieldIdentifier: badorm.UpdatedAtFieldID,
	}
}
func GoEmbeddedDeletedAt(expr badorm.Expression[gorm.DeletedAt]) badorm.WhereCondition[goembedded.GoEmbedded] {
	return badorm.FieldCondition[goembedded.GoEmbedded, gorm.DeletedAt]{
		Expression:      expr,
		FieldIdentifier: badorm.DeletedAtFieldID,
	}
}

var goEmbeddedIntFieldID = badorm.FieldIdentifier{Field: "Int"}

func GoEmbeddedInt(expr badorm.Expression[int]) badorm.WhereCondition[goembedded.GoEmbedded] {
	return badorm.FieldCondition[goembedded.GoEmbedded, int]{
		Expression:      expr,
		FieldIdentifier: goEmbeddedIntFieldID,
	}
}

var goEmbeddedToBeEmbeddedIntFieldID = badorm.FieldIdentifier{Field: "Int"}

func GoEmbeddedToBeEmbeddedInt(expr badorm.Expression[int]) badorm.WhereCondition[goembedded.GoEmbedded] {
	return badorm.FieldCondition[goembedded.GoEmbedded, int]{
		Expression:      expr,
		FieldIdentifier: goEmbeddedToBeEmbeddedIntFieldID,
	}
}

var GoEmbeddedPreloadAttributes = badorm.NewPreloadCondition[goembedded.GoEmbedded](goEmbeddedIntFieldID, goEmbeddedToBeEmbeddedIntFieldID)
