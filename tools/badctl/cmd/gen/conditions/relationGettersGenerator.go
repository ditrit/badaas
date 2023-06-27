package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/ditrit/badaas/tools/badctl/cmd/log"
)

const (
	badORMVerifyStructLoaded        = "VerifyStructLoaded"
	badORMVerifyPointerLoaded       = "VerifyPointerLoaded"
	badORMVerifyPointerWithIDLoaded = "VerifyPointerWithIDLoaded"
)

type RelationGettersGenerator struct {
	object     types.Object
	objectType Type
}

func NewRelationGettersGenerator(object types.Object) *RelationGettersGenerator {
	return &RelationGettersGenerator{
		object:     object,
		objectType: Type{object.Type()},
	}
}

// Add conditions for an object in the file
func (generator RelationGettersGenerator) Into(file *File) error {
	fields, err := getFields(generator.objectType)
	if err != nil {
		return err
	}

	log.Logger.Infof("Generating relation getters for type %q in %s", generator.object.Name(), file.name)

	file.Add(generator.ForEachField(file, fields)...)

	return nil
}

func (generator RelationGettersGenerator) ForEachField(file *File, fields []Field) []jen.Code {
	relationGetters := []jen.Code{}

	for _, field := range fields {
		if field.Embedded {
			relationGetters = append(
				relationGetters,
				generateForEmbeddedField[jen.Code](
					file,
					field,
					generator,
				)...,
			)
		} else {
			getterForField := generator.generateForField(field)
			if getterForField != nil {
				relationGetters = append(relationGetters, getterForField)
			}
		}
	}

	return relationGetters
}

func (generator RelationGettersGenerator) generateForField(field Field) jen.Code {
	switch fieldType := field.GetType().(type) {
	case *types.Named:
		// the field is a named type (user defined structs)
		_, err := field.Type.BadORMModelStruct()
		if err == nil {
			log.Logger.Debugf("Generating relation getter for type %q and field %s", generator.object.Name(), field.Name)
			// field is a BaDORM Model
			return generator.verifyStruct(field)
		}
	case *types.Pointer:
		// the field is a pointer
		return generator.generateForPointer(field.ChangeType(fieldType.Elem()))
	default:
		log.Logger.Debugf("struct field type not handled: %T", fieldType)
	}

	return nil
}

func (generator RelationGettersGenerator) generateForPointer(field Field) jen.Code {
	_, err := field.Type.BadORMModelStruct()
	if err == nil {
		// field is a pointer to BaDORM Model
		fk, err := generator.objectType.GetFK(field)
		if err != nil {
			log.Logger.Debugf("unhandled: field is a pointer and object not has the fk: %s", field.Type)
			return nil
		}

		log.Logger.Debugf("Generating relation getter for type %q and field %s", generator.object.Name(), field.Name)

		switch fk.GetType().(type) {
		case *types.Named:
			if fk.IsBadORMID() {
				return generator.verifyPointerWithID(field)
			}
		case *types.Pointer:
			// the fk is a pointer
			return generator.verifyPointer(field)
		}
	}

	return nil
}

func getGetterName(field Field) string {
	return "Get" + strcase.ToPascal(field.Name)
}

func (generator RelationGettersGenerator) verifyStruct(field Field) *jen.Statement {
	return generator.verifyCommon(
		field,
		badORMVerifyStructLoaded,
		jen.Op("&").Id("m").Op(".").Id(field.Name),
	)
}

func (generator RelationGettersGenerator) verifyPointer(field Field) *jen.Statement {
	return generator.verifyPointerCommon(field, badORMVerifyPointerLoaded)
}

func (generator RelationGettersGenerator) verifyPointerWithID(field Field) *jen.Statement {
	return generator.verifyPointerCommon(field, badORMVerifyPointerWithIDLoaded)
}

func (generator RelationGettersGenerator) verifyPointerCommon(field Field, verifyFunc string) *jen.Statement {
	return generator.verifyCommon(
		field,
		verifyFunc,
		jen.Id("m").Op(".").Id(field.Name+"ID"),
		jen.Id("m").Op(".").Id(field.Name),
	)
}

func (generator RelationGettersGenerator) verifyCommon(field Field, verifyFunc string, callParams ...jen.Code) *jen.Statement {
	return jen.Func().Parens(
		jen.Id("m").Id(generator.object.Name()),
	).Id(getGetterName(field)).Params().Add(
		jen.Parens(
			jen.List(
				jen.Op("*").Qual(
					getRelativePackagePath(
						generator.object.Pkg().Name(),
						field.Type,
					),
					field.TypeName(),
				),

				jen.Id("error"),
			),
		),
	).Block(
		jen.Return(
			jen.Qual(
				badORMPath,
				verifyFunc,
			).Types(
				jen.Id(field.TypeName()),
			).Call(
				callParams...,
			),
		),
	)
}
