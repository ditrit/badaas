package conditions

import (
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
	"github.com/elliotchance/pie/v2"
	"github.com/ettle/strcase"
)

type Condition struct {
	codes []jen.Code
	param *jen.Statement
}

func NewCondition(objectType types.Type, field Field) *Condition {
	condition := &Condition{
		param: jen.Id("v"),
	}
	condition.generateCode(objectType, field)
	return condition
}

func (condition *Condition) generateCode(objectType types.Type, field Field) {
	switch fieldType := field.Type.(type) {
	case *types.Basic:
		condition.adaptParamByKind(fieldType)
		condition.generateWhereCondition(
			objectType,
			field,
		)
	case *types.Named:
		condition.generateCodeForNamedType(
			objectType,
			field,
		)
	case *types.Pointer:
		condition.param = condition.param.Op("*")
		condition.generateCode(
			objectType,
			field.ChangeType(fieldType.Elem()),
		)
	case *types.Slice:
		condition.param = condition.param.Index()
		condition.generateCodeForSlice(
			objectType,
			field.ChangeType(fieldType.Elem()),
		)
	default:
		log.Printf("struct field type not handled: %T", fieldType)
	}
}

func (condition *Condition) generateCodeForSlice(objectType types.Type, field Field) {
	switch elemType := field.Type.(type) {
	case *types.Basic:
		// una list de strings o algo asi,
		// por el momento solo anda con []byte porque el resto gorm no lo sabe encodear
		condition.generateCode(
			objectType,
			field,
		)
	case *types.Named:
		// inverse relation condition
		_, err := getBadORMModelStruct(field.Type)
		if err == nil {
			// slice of BadORM models
			log.Println(field.TypeName())
			condition.generateOppositeJoinCondition(
				objectType,
				field,
			)
		}
	case *types.Pointer:
		condition.param = condition.param.Op("*")
		// slice de pointers, solo testeado temporalmente porque despues gorm no lo soporta
		condition.generateCodeForSlice(
			objectType,
			field.ChangeType(elemType.Elem()),
		)
	default:
		log.Printf("struct field list elem type not handled: %T", elemType)
	}
}

func (condition *Condition) generateCodeForNamedType(objectType types.Type, field Field) {
	// TODO esta linea de aca quedo rara
	_, err := getBadORMModelStruct(field.Type)
	log.Println(err)

	if err == nil {
		objectStruct, err := getBadORMModelStruct(objectType)
		if err != nil {
			// TODO ver esto
			return
		}
		// TODO que pasa si esta en otro package? se importa solo?
		fields, err := getFields(
			objectStruct,
			// TODO testear esto si esta bien aca
			field.Tags.getEmbeddedPrefix(),
		)
		if err != nil {
			// TODO ver esto
			return
		}
		thisEntityHasTheFK := pie.Any(fields, func(otherField Field) bool {
			return otherField.Name == field.getJoinFromColumn()
		})

		log.Println(field.getJoinFromColumn())
		log.Println(thisEntityHasTheFK)

		if thisEntityHasTheFK {
			// belongsTo relation
			condition.generateJoinCondition(
				objectType,
				field,
			)
		} else {
			// hasOne or hasMany relation
			condition.generateInverseJoinCondition(
				objectType,
				field,
			)

			condition.generateOppositeJoinCondition(
				objectType,
				field,
			)
		}
	} else {
		if (field.IsGormCustomType() || field.TypeString() == "time.Time") && field.TypeString() != "gorm.io/gorm.DeletedAt" {
			// TODO DeletedAt
			condition.param = condition.param.Qual(
				getRelativePackagePath(field.TypePkg()),
				field.TypeName(),
			)
			condition.generateWhereCondition(
				objectType,
				field,
			)
		} else {
			log.Printf("struct field type not handled: %s", field.TypeString())
		}
	}
}

func (condition *Condition) adaptParamByKind(basicType *types.Basic) {
	switch basicType.Kind() {
	case types.Bool:
		condition.param = condition.param.Bool()
	case types.Int:
		condition.param = condition.param.Int()
	case types.Int8:
		condition.param = condition.param.Int8()
	case types.Int16:
		condition.param = condition.param.Int16()
	case types.Int32:
		condition.param = condition.param.Int32()
	case types.Int64:
		condition.param = condition.param.Int64()
	case types.Uint:
		condition.param = condition.param.Uint()
	case types.Uint8:
		condition.param = condition.param.Uint8()
	case types.Uint16:
		condition.param = condition.param.Uint16()
	case types.Uint32:
		condition.param = condition.param.Uint32()
	case types.Uint64:
		condition.param = condition.param.Uint64()
	case types.Uintptr:
		condition.param = condition.param.Uintptr()
	case types.Float32:
		condition.param = condition.param.Float32()
	case types.Float64:
		condition.param = condition.param.Float64()
	case types.Complex64:
		condition.param = condition.param.Complex64()
	case types.Complex128:
		condition.param = condition.param.Complex128()
	case types.String:
		condition.param = condition.param.String()
	}
}

// TODO sacar condition del nombre
func (condition *Condition) generateWhereCondition(objectType types.Type, field Field) {
	whereCondition := jen.Qual(
		badORMPath, badORMWhereCondition,
	).Types(
		jen.Qual(
			getRelativePackagePath(getTypePkg(objectType)),
			getTypeName(objectType),
		),
	)

	condition.codes = append(
		condition.codes,
		jen.Func().Id(
			getConditionName(objectType, field.Name),
		).Params(
			condition.param,
		).Add(
			whereCondition.Clone(),
		).Block(
			jen.Return(
				whereCondition.Clone().Values(jen.Dict{
					jen.Id("Field"): jen.Lit(field.getColumnName()),
					jen.Id("Value"): jen.Id("v"),
				}),
			),
		),
	)
}

func (condition *Condition) generateOppositeJoinCondition(objectType types.Type, field Field) {
	condition.generateJoinCondition(
		field.Type,
		// TODO testear los Override Foreign Key
		Field{
			Name: getTypeName(objectType),
			Type: objectType,
			Tags: field.Tags,
		},
	)
}

func (condition *Condition) generateJoinCondition(objectType types.Type, field Field) {
	condition.generateJoinFromAndTo(
		objectType,
		field,
		field.getJoinFromColumn(),
		field.getJoinToColumn(),
	)
}

func (condition *Condition) generateInverseJoinCondition(objectType types.Type, field Field) {
	condition.generateJoinFromAndTo(
		objectType,
		field,
		field.getJoinToColumn(),
		field.NoSePonerNombre(getTypeName(objectType)),
	)
}

func (condition *Condition) generateJoinFromAndTo(objectType types.Type, field Field, from, to string) {
	log.Println(field.Name)

	t1 := jen.Qual(
		getRelativePackagePath(getTypePkg(objectType)),
		getTypeName(objectType),
	)

	t2 := jen.Qual(
		getRelativePackagePath(field.TypePkg()),
		field.TypeName(),
	)

	badormT1Condition := jen.Qual(
		badORMPath, badORMCondition,
	).Types(t1)
	badormT2Condition := jen.Qual(
		badORMPath, badORMCondition,
	).Types(t2)
	badormJoinCondition := jen.Qual(
		badORMPath, badORMJoinCondition,
	).Types(
		t1, t2,
	)

	condition.codes = append(
		condition.codes,
		jen.Func().Id(
			getConditionName(objectType, field.Name),
		).Params(
			jen.Id("conditions").Op("...").Add(badormT2Condition),
		).Add(
			badormT1Condition,
		).Block(
			jen.Return(
				badormJoinCondition.Values(jen.Dict{
					jen.Id("T1Field"):    jen.Lit(strcase.ToSnake(from)),
					jen.Id("T2Field"):    jen.Lit(strcase.ToSnake(to)),
					jen.Id("Conditions"): jen.Id("conditions"),
				}),
			),
		),
	)
}

func getConditionName(typeV types.Type, fieldName string) string {
	return getTypeName(typeV) + strcase.ToPascal(fieldName) + badORMCondition
}

// TODO testear esto
func getRelativePackagePath(srcPkg *types.Package) string {
	if srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}
