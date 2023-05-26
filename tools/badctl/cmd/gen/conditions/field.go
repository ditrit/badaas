package conditions

import (
	"errors"
	"go/types"
	"regexp"

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

func (field Field) TypeString() string {
	return field.Object.Type().String()
}

func (field Field) TypeName() string {
	return getObjectTypeName(field.Object)
}

func (field Field) TypePkg() *types.Package {
	fieldType := field.Object.Type()
	switch fieldTypeTyped := fieldType.(type) {
	case *types.Named:
		return fieldTypeTyped.Obj().Pkg()
	// TODO ver el resto si al hacerlo me simplificaria algo
	default:
		return nil
	}
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

var scanMethod = regexp.MustCompile(`func \(\*.*\)\.Scan\([a-zA-Z0-9_-]* interface\{\}\) error$`)
var valueMethod = regexp.MustCompile(`func \(.*\)\.Value\(\) \(database/sql/driver\.Value\, error\)$`)

func (field Field) IsGormCustomType() bool {
	typeNamed, isNamedType := field.Object.Type().(*types.Named)
	if !isNamedType {
		return false
	}

	hasScanMethod := false
	hasValueMethod := false
	for i := 0; i < typeNamed.NumMethods(); i++ {
		methodSignature := typeNamed.Method(i).String()

		if !hasScanMethod && scanMethod.MatchString(methodSignature) {
			hasScanMethod = true
		} else if !hasValueMethod && valueMethod.MatchString(methodSignature) {
			hasValueMethod = true
		}
	}

	return hasScanMethod && hasValueMethod
}

func getFields(structType *types.Struct, prefix string) ([]Field, error) {
	numFields := structType.NumFields()
	if numFields == 0 {
		return nil, errors.New("struct has 0 fields")
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
