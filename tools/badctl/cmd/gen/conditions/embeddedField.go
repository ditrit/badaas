package conditions

import (
	"errors"
	"go/types"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
)

// Generate conditions for a embedded field
// it will generate a condition for each of the field of the embedded field's type
func generateForEmbeddedField[T any](file *File, field Field, generator CodeGenerator[T]) []T {
	embeddedStructType, ok := field.Type.Underlying().(*types.Struct)
	if !ok {
		cmderrors.FailErr(errors.New("unreachable! embedded objects are always structs"))
	}

	fields, err := getStructFields(embeddedStructType)
	if err != nil {
		// embedded field's type has not fields
		return []T{}
	}

	if !isBadORMBaseModel(field.TypeString()) {
		fields = pie.Map(fields, func(embeddedField Field) Field {
			embeddedField.ColumnPrefix = field.Tags.getEmbeddedPrefix()
			embeddedField.NamePrefix = field.Name

			return embeddedField
		})
	}

	return generator.ForEachField(file, fields)
}
