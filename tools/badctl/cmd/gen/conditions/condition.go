package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ditrit/badaas/tools/badctl/cmd/log"
	"github.com/ettle/strcase"
)

type Condition struct {
	codes   []jen.Code
	param   *JenParam
	destPkg string
}

func NewCondition(destPkg string, objectType Type, field Field) (*Condition, error) {
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

// Generate the condition between the object and the field
func (condition *Condition) generate(objectType Type, field Field) error {
	switch fieldType := field.Type.Type.(type) {
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
		log.Logger.Debugf("struct field type not handled: %T", fieldType)
	}

	return nil
}

// Generate condition between the object and the field when the field is a slice
func (condition *Condition) generateForSlice(objectType Type, field Field) {
	switch elemType := field.Type.Type.(type) {
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
		_, err := field.Type.BadORMModelStruct()
		if err == nil {
			// slice of BadORM models -> hasMany relation
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
		log.Logger.Debugf("struct field list elem type not handled: %T", elemType)
	}
}

// Generate condition between object and field when the field is a named type (user defined struct)
func (condition *Condition) generateForNamedType(objectType Type, field Field) error {
	_, err := field.Type.BadORMModelStruct()

	if err == nil {
		// field is a BaDORM Model
		// TODO que pasa si esta en otro package? se importa solo?
		hasFK, err := objectType.HasFK(field)
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
		if (field.Type.IsGormCustomType() || field.TypeString() == "time.Time") && field.TypeString() != "gorm.io/gorm.DeletedAt" {
			// TODO DeletedAt
			// field is a Gorm Custom type (implements Scanner and Valuer interfaces)
			// or a named type supported by gorm (time.Time, gorm.DeletedAt)
			condition.param.ToCustomType(condition.destPkg, field.Type)
			condition.generateWhere(
				objectType,
				field,
			)
		} else {
			log.Logger.Debugf("struct field type not handled: %s", field.TypeString())
		}
	}

	return nil
}

// Generate a WhereCondition between object and field
func (condition *Condition) generateWhere(objectType Type, field Field) {
	whereCondition := jen.Qual(
		badORMPath, badORMWhereCondition,
	).Types(
		jen.Qual(
			getRelativePackagePath(condition.destPkg, objectType),
			objectType.Name(),
		),
	)

	conditionName := getConditionName(objectType, field.Name)
	log.Logger.Debugf("Generated %q", conditionName)

	condition.codes = append(
		condition.codes,
		jen.Func().Id(
			conditionName,
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

// Generate a inverse JoinCondition, so from the field's object to the object
func (condition *Condition) generateInverseJoin(objectType Type, field Field) {
	condition.generateJoinWithFK(
		field.Type,
		// TODO testear los Override Foreign Key
		Field{
			Name: objectType.Name(),
			Type: objectType,
			Tags: field.Tags,
		},
	)
}

// Generate a JoinCondition between the object and field's object
// when object has a foreign key to the field's object
func (condition *Condition) generateJoinWithFK(objectType Type, field Field) {
	condition.generateJoin(
		objectType,
		field,
		field.getFKAttribute(),
		field.getFKReferencesAttribute(),
	)
}

// Generate a JoinCondition between the object and field's object
// when object has not a foreign key to the field's object
// (so the field's object has it)
func (condition *Condition) generateJoinWithoutFK(objectType Type, field Field) {
	condition.generateJoin(
		objectType,
		field,
		field.getFKReferencesAttribute(),
		field.getRelatedTypeFKAttribute(objectType.Name()),
	)
}

// Generate a JoinCondition
func (condition *Condition) generateJoin(objectType Type, field Field, t1Field, t2Field string) {
	t1 := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	t2 := jen.Qual(
		getRelativePackagePath(condition.destPkg, field.Type),
		field.TypeName(),
	)

	conditionName := getConditionName(objectType, field.Name)
	log.Logger.Debugf("Generated %q", conditionName)

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
			conditionName,
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

// Generate condition names
func getConditionName(typeV Type, fieldName string) string {
	return typeV.Name() + strcase.ToPascal(fieldName)
}

// TODO testear esto
// Avoid importing the same package as the destination one
func getRelativePackagePath(destPkg string, typeV Type) string {
	srcPkg := typeV.Pkg()
	if srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}
