package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ettle/strcase"

	"github.com/ditrit/badaas/tools/badctl/cmd/log"
)

const (
	// badorm/condition.go
	badORMCondition       = "Condition"
	badORMFieldCondition  = "FieldCondition"
	badORMWhereCondition  = "WhereCondition"
	badORMJoinCondition   = "JoinCondition"
	badORMFieldIdentifier = "FieldIdentifier"
	IDFieldID             = "IDFieldID"
	CreatedAtFieldID      = "CreatedAtFieldID"
	UpdatedAtFieldID      = "UpdatedAtFieldID"
	DeletedAtFieldID      = "DeletedAtFieldID"
	// badorm/expression.go
	badORMExpression = "Expression"
)

var constantFieldIdentifiers = map[string]*jen.Statement{
	"ID":        jen.Qual(badORMPath, IDFieldID),
	"CreatedAt": jen.Qual(badORMPath, CreatedAtFieldID),
	"UpdatedAt": jen.Qual(badORMPath, UpdatedAtFieldID),
	"DeletedAt": jen.Qual(badORMPath, DeletedAtFieldID),
}

type Condition struct {
	codes           []jen.Code
	param           *JenParam
	destPkg         string
	fieldIdentifier string
	preloadName     string
}

func NewCondition(destPkg string, objectType Type, field Field) *Condition {
	condition := &Condition{
		param:   NewJenParam(),
		destPkg: destPkg,
	}
	condition.generate(objectType, field)

	return condition
}

// Generate the condition between the object and the field
func (condition *Condition) generate(objectType Type, field Field) {
	switch fieldType := field.GetType().(type) {
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
		condition.generateForNamedType(
			objectType,
			field,
		)
	case *types.Pointer:
		// the field is a pointer
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
}

// Generate condition between the object and the field when the field is a slice
func (condition *Condition) generateForSlice(objectType Type, field Field) {
	switch elemType := field.GetType().(type) {
	case *types.Basic:
		// slice of basic types ([]string, []int, etc.)
		// the only one supported directly by gorm is []byte
		// but the others can be used after configuration in some dbs
		condition.generate(
			objectType,
			field,
		)
	case *types.Pointer:
		// slice of pointers, generate code for a slice of the pointed type
		condition.generateForSlice(
			objectType,
			field.ChangeType(elemType.Elem()),
		)
	default:
		log.Logger.Debugf("struct field list elem type not handled: %T", elemType)
	}
}

// Generate condition between object and field when the field is a named type (user defined struct)
func (condition *Condition) generateForNamedType(objectType Type, field Field) {
	_, err := field.Type.BadORMModelStruct()

	if err == nil {
		// field is a BaDORM Model
		condition.generateForBadormModel(objectType, field)
	} else {
		// field is not a BaDORM Model
		if field.Type.IsGormCustomType() || field.TypeString() == "time.Time" {
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
}

// Generate condition between object and field when the field is a BaDORM Model
func (condition *Condition) generateForBadormModel(objectType Type, field Field) {
	hasFK, _ := objectType.HasFK(field)
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
	}
}

// Generate a WhereCondition between object and field
func (condition *Condition) generateWhere(objectType Type, field Field) {
	objectTypeQual := jen.Qual(
		getRelativePackagePath(condition.destPkg, objectType),
		objectType.Name(),
	)

	fieldCondition := jen.Qual(
		badORMPath, badORMFieldCondition,
	).Types(
		objectTypeQual,
		condition.param.GenericType(),
	)

	whereCondition := jen.Qual(
		badORMPath, badORMWhereCondition,
	).Types(
		objectTypeQual,
	)

	conditionName := getConditionName(objectType, field)
	log.Logger.Debugf("Generated %q", conditionName)

	var fieldIdentifier *jen.Statement

	if constantFieldIdentifier, ok := constantFieldIdentifiers[field.Name]; ok {
		fieldIdentifier = constantFieldIdentifier
	} else {
		fieldIdentifier = condition.createFieldIdentifier(objectType, field, conditionName)
	}

	condition.codes = append(
		condition.codes,
		jen.Func().Id(
			conditionName,
		).Params(
			condition.param.Statement(),
		).Add(
			whereCondition,
		).Block(
			jen.Return(
				fieldCondition.Clone().Values(jen.Dict{
					jen.Id("Expression"):      jen.Id("expr"),
					jen.Id("FieldIdentifier"): fieldIdentifier,
				}),
			),
		),
	)
}

// create a variable containing the definition of the field identifier
// to use it in the where condition and in the preload condition
func (condition *Condition) createFieldIdentifier(objectType Type, field Field, conditionName string) *jen.Statement {
	fieldIdentifierValues := jen.Dict{}
	columnName := field.getColumnName()

	if columnName != "" {
		fieldIdentifierValues[jen.Id("Column")] = jen.Lit(columnName)
	} else {
		fieldIdentifierValues[jen.Id("Field")] = jen.Lit(field.Name)
	}

	columnPrefix := field.ColumnPrefix
	if columnPrefix != "" {
		fieldIdentifierValues[jen.Id("ColumnPrefix")] = jen.Lit(columnPrefix)
	}

	fieldIdentifierVar := jen.Qual(
		badORMPath, badORMFieldIdentifier,
	).Values(fieldIdentifierValues)

	fieldIdentifierName := strcase.ToCamel(conditionName) + "FieldID"

	condition.codes = append(
		condition.codes,
		jen.Var().Id(
			fieldIdentifierName,
		).Op("=").Add(fieldIdentifierVar),
	)

	condition.fieldIdentifier = fieldIdentifierName

	return jen.Qual("", fieldIdentifierName)
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

	conditionName := getConditionName(objectType, field)
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
					jen.Id("T1Field"):            jen.Lit(t1Field),
					jen.Id("T2Field"):            jen.Lit(t2Field),
					jen.Id("RelationField"):      jen.Lit(field.Name),
					jen.Id("Conditions"):         jen.Id("conditions"),
					jen.Id("T1PreloadCondition"): jen.Id(getPreloadAttributesName(objectType.Name())),
				}),
			),
		),
	)

	// preload for the relation
	preloadName := objectType.Name() + "Preload" + field.Name
	condition.codes = append(
		condition.codes,
		jen.Var().Id(
			preloadName,
		).Op("=").Add(jen.Id(conditionName)).Call(
			jen.Id(getPreloadAttributesName(field.TypeName())),
		),
	)
	condition.preloadName = preloadName
}

// Generate condition names
func getConditionName(typeV Type, field Field) string {
	return typeV.Name() + strcase.ToPascal(field.NamePrefix) + strcase.ToPascal(field.Name)
}

// Avoid importing the same package as the destination one
func getRelativePackagePath(destPkg string, typeV Type) string {
	srcPkg := typeV.Pkg()
	if srcPkg == nil || srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}
