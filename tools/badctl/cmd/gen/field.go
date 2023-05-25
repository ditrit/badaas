package gen

import (
	"go/types"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/ettle/strcase"
)

type Field struct {
	Name     string
	Type     types.Object
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

func (field Field) TypeName() string {
	return pie.Last(strings.Split(field.Type.Type().String(), "."))
}

func (field Field) ChangeType(newType types.Type) Field {
	return Field{
		Name: field.Name,
		Type: types.NewTypeName(
			field.Type.Pos(),
			field.Type.Pkg(),
			field.Type.Name(),
			newType,
		),
		Tags: field.Tags,
	}
}

func getFields(structType *types.Struct, prefix string) []Field {
	fields := []Field{}

	// Iterate over struct fields
	for i := 0; i < structType.NumFields(); i++ {
		fieldType := structType.Field(i)
		gormTags := getGormTags(structType.Tag(i))
		fields = append(fields, Field{
			// TODO el Name se podria sacar si meto este prefix adentro del tipo
			Name:     prefix + fieldType.Name(),
			Type:     fieldType,
			Embedded: fieldType.Embedded() || gormTags.hasEmbedded(),
			Tags:     gormTags,
		})
	}

	return fields
}
