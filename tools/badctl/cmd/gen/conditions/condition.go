package conditions

import (
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"
)

type Condition struct {
	codes   []jen.Code
	param   *JenParam
	destPkg string
}

func NewCondition(destPkg string, objectType types.Type, field Field) (*Condition, error) {
	condition := &Condition{
		param:   NewJenParam(),
		destPkg: destPkg,
	}
	err := condition.generate(objectType, field)
	if err != nil {
		return nil, err
	}

	return condition, nil
}

func (condition *Condition) generate(objectType types.Type, field Field) error {
	switch fieldType := field.Type.(type) {
	case *types.Basic:
		// the field is a basic type (string, int, etc)
		// adapt param to that type and generate a WhereCondition
		condition.param.ToBasicKind(fieldType)
		condition.generateWhere(
			objectType,
			field,
		)
	case *types.Named:
		// the field is a named type (user defined structs)
		return condition.generateForNamedType(
			objectType,
			field,
		)
	case *types.Pointer:
		// the field is a pointer
		// adapt param to pointer and generate code for the pointed type
		condition.param.ToPointer()
		condition.generate(
			objectType,
			field.ChangeType(fieldType.Elem()),
		)
	case *types.Slice:
		// the field is a slice
		// adapt param to slice and
		// generate code for the type of the elements of the slice
		condition.param.ToSlice()
		condition.generateForSlice(
			objectType,
			field.ChangeType(fieldType.Elem()),
		)
	default:
		log.Printf("struct field type not handled: %T", fieldType)
	}

	return nil
}

func (condition *Condition) generateForSlice(objectType types.Type, field Field) {
	switch elemType := field.Type.(type) {
	case *types.Basic:
		// slice of basic types ([]string, []int, etc.)
		// the only one supported directly by gorm is []byte
		// but the others can be used after configuration in some dbs
		condition.generate(
			objectType,
			field,
		)
	case *types.Named:
		// slice of named types (user defined types)
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
		// TODO pointer solo testeado temporalmente porque despues gorm no lo soporta
		// slice of pointers, generate code for a slice of the pointed type
		// after changing the param
		condition.param.ToPointer()
		condition.generateForSlice(
			objectType,
			field.ChangeType(elemType.Elem()),
		)
	default:
		log.Printf("struct field list elem type not handled: %T", elemType)
	}
}

func (condition *Condition) generateForNamedType(objectType types.Type, field Field) error {
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
			// field is a Gorm Custom type (implements Scanner and Valuer interfaces)
			// or a named type supported by gorm (time.Time, gorm.DeletedAt)
			condition.param.ToCustomType(condition.destPkg, field.Type)
			condition.generateWhere(
				objectType,
				field,
			)
		} else {
			log.Printf("struct field type not handled: %s", field.TypeString())
		}
	}

	return nil
}

func (condition *Condition) generateWhere(objectType types.Type, field Field) {
	whereCondition := jen.Qual(
		badORMPath, badORMWhereCondition,
	).Types(
		jen.Qual(
			getRelativePackagePath(condition.destPkg, objectType),
			getTypeName(objectType),
		),
	)

	condition.codes = append(
		condition.codes,
		jen.Func().Id(
			getConditionName(objectType, field.Name),
		).Params(
			condition.param.Statement(),
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
		getRelativePackagePath(condition.destPkg, objectType),
		getTypeName(objectType),
	)

	t2 := jen.Qual(
		getRelativePackagePath(condition.destPkg, field.Type),
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
// avoid importing the same package as the destination one
func getRelativePackagePath(destPkg string, typeV types.Type) string {
	srcPkg := getTypePkg(typeV)
	if srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}
