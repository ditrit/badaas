package conditions

import (
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"
)

type Condition struct {
	codes []jen.Code
	param *jen.Statement
}

func NewCondition(objectType types.Type, field Field) (*Condition, error) {
	condition := &Condition{
		param: jen.Id("v"),
	}
	err := condition.generateCode(objectType, field)
	if err != nil {
		return nil, err
	}

	return condition, nil
}

func (condition *Condition) generateCode(objectType types.Type, field Field) error {
	switch fieldType := field.Type.(type) {
	case *types.Basic:
		condition.adaptParamByKind(fieldType)
		condition.generateWhereCondition(
			objectType,
			field,
		)
	case *types.Named:
		return condition.generateCodeForNamedType(
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

	return nil
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
		_, err := getBadORMModelStruct(field.Type)
		if err == nil {
			// slice of BadORM models -> hasMany relation
			log.Println(field.TypeName())
			condition.generateInverseJoin(
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

func (condition *Condition) generateCodeForNamedType(objectType types.Type, field Field) error {
	_, err := getBadORMModelStruct(field.Type)

	if err == nil {
		// field is a BaDORM Model
		// TODO que pasa si esta en otro package? se importa solo?

		hasFK, err := hasFK(objectType, field)
		if err != nil {
			return err
		}

		if hasFK {
			// belongsTo relation
			condition.generateJoinWithFK(
				objectType,
				field,
			)
		} else {
			// hasOne relation
			condition.generateJoinWithoutFK(
				objectType,
				field,
			)

			condition.generateInverseJoin(
				objectType,
				field,
			)
		}
	} else {
		// field is not a BaDORM Model
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

	return nil
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

func (condition *Condition) generateInverseJoin(objectType types.Type, field Field) {
	condition.generateJoinWithFK(
		field.Type,
		// TODO testear los Override Foreign Key
		Field{
			Name: getTypeName(objectType),
			Type: objectType,
			Tags: field.Tags,
		},
	)
}

func (condition *Condition) generateJoinWithFK(objectType types.Type, field Field) {
	condition.generateJoin(
		objectType,
		field,
		field.getFKAttribute(),
		field.getFKReferencesAttribute(),
	)
}

func (condition *Condition) generateJoinWithoutFK(objectType types.Type, field Field) {
	condition.generateJoin(
		objectType,
		field,
		field.getFKReferencesAttribute(),
		field.getRelatedTypeFKAttribute(getTypeName(objectType)),
	)
}

func (condition *Condition) generateJoin(objectType types.Type, field Field, t1Field, t2Field string) {
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
					jen.Id("T1Field"):    jen.Lit(strcase.ToSnake(t1Field)),
					jen.Id("T2Field"):    jen.Lit(strcase.ToSnake(t2Field)),
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
