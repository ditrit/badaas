package gen

import (
	"go/types"

	"github.com/ettle/strcase"
)

type Field struct {
	Name     string
	Type     types.Type
	Embedded bool
	Tags     GormTags
}

func (field Field) getColumnName() string {
	columnTag, isPresent := field.Tags[columnTagName]
	if isPresent {
		return columnTag
	}

	return strcase.ToSnake(field.Name)
}

func (field Field) getJoinFromColumn() string {
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		return foreignKeyTag
	}

	return field.Name + "ID"
}

func (field Field) getJoinToColumn() string {
	referencesTag, isPresent := field.Tags[referencesTagName]
	// TODO testear cuando hay redefinicion en la inversa
	if isPresent {
		return referencesTag
	}

	return "ID"
}

// TODO
func (field Field) NoSePonerNombre(structName string) string {
	// TODO testear cuando hay redefinicion
	foreignKeyTag, isPresent := field.Tags[foreignKeyTagName]
	if isPresent {
		return foreignKeyTag
	}

	return structName + "ID"
}

func getFields(structType *types.Struct, prefix string) []Field {
	fields := []Field{}

	// Iterate over struct fields
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		gormTags := getGormTags(structType.Tag(i))
		fields = append(fields, Field{
			Name:     prefix + field.Name(),
			Type:     field.Type(),
			Embedded: field.Embedded() || gormTags.hasEmbedded(),
			Tags:     gormTags,
		})
	}

	return fields
}
