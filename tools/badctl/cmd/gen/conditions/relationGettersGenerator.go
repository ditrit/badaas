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

// TODO codigo repetido con el conditions generator
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
func (generator RelationGettersGenerator) GenerateInto(file *File) error {
	// TODO codigo repetido con file
	fields, _ := getFields(generator.objectType)
	relationGetters := []jen.Code{}

	for _, field := range fields {
		if field.Embedded {
			// TODO
		} else {
			// TODO codigo repetido con condition.go
			switch fieldType := field.GetType().(type) {
			case *types.Named:
				// the field is a named type (user defined structs)
				_, err := field.Type.BadORMModelStruct()

				if err == nil {
					// field is a BaDORM Model
					relationGetters = append(
						relationGetters,
						generator.verifyStruct(field),
					)
				}
			case *types.Pointer:
				// the field is a pointer
				_, err := field.ChangeType(fieldType.Elem()).Type.BadORMModelStruct()

				if err == nil {
					// field is a BaDORM Model
					fk, err := generator.objectType.GetFK(field)
					if err != nil {
						log.Logger.Debugf("unhandled: field is a pointer and object not has the fk: %T", fieldType)
						continue
					}

					switch fk.GetType().(type) {
					// TODO verificar que sea de los ids correctos?
					// TODO basics para strings y eso?
					case *types.Named:
						relationGetters = append(
							relationGetters,
							generator.verifyPointerWithID(field),
						)
					case *types.Pointer:
						relationGetters = append(
							relationGetters,
							generator.verifyPointer(field),
						)
					}
				}
			default:
				log.Logger.Debugf("struct field type not handled: %T", fieldType)
			}
		}
	}

	file.Add(relationGetters...)

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
