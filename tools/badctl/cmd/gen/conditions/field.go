package conditions

import (
	"errors"
	"go/types"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/ettle/strcase"
)

type Field struct {
	Name     string
	Object   types.Object
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
	return pie.Last(strings.Split(field.Object.Type().String(), "."))
}

func (field Field) ChangeType(newType types.Type) Field {
	return Field{
		Name: field.Name,
		Object: types.NewTypeName(
			field.Object.Pos(),
			field.Object.Pkg(),
			field.Object.Name(),
			newType,
		),
		Tags: field.Tags,
	}
}

func getFields(structType *types.Struct, prefix string) ([]Field, error) {
	numFields := structType.NumFields()
	if numFields == 0 {
		return nil, errors.New("Struct has 0 fields")
	}

	fields := []Field{}

	// Iterate over struct fields
	for i := 0; i < numFields; i++ {
		fieldType := structType.Field(i)
		gormTags := getGormTags(structType.Tag(i))
		fields = append(fields, Field{
			// TODO el Name se podria sacar si meto este prefix adentro del tipo
			Name:     prefix + fieldType.Name(),
			Object:   fieldType,
			Embedded: fieldType.Embedded() || gormTags.hasEmbedded(),
			Tags:     gormTags,
		})
	}

	return fields, nil
}
