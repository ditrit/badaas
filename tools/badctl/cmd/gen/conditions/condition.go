package conditions

import (
	"go/types"
	"log"
	"regexp"

	"github.com/dave/jennifer/jen"
	"github.com/elliotchance/pie/v2"
	"github.com/ettle/strcase"
)

type Condition struct {
	codes []jen.Code
	param *jen.Statement
}

func NewCondition(object types.Object, field Field) *Condition {
	condition := &Condition{
		param: jen.Id("v"),
	}
	condition.generateCode(object, field)
	return condition
}

func (condition *Condition) generateCode(object types.Object, field Field) {
	switch fieldType := field.Object.Type().(type) {
	case *types.Basic:
		condition.codes = []jen.Code{generateWhereCondition(
			object,
			field,
			typeKindToJenStatement[fieldType.Kind()](condition.param),
		)}
	case *types.Named:
		condition.codes = generateConditionsForNamedType(
			object,
			field, fieldType,
			condition.param,
		)
	case *types.Pointer:
		condition.param = condition.param.Op("*")
		condition.generateCode(
			object,
			field.ChangeType(fieldType.Elem()),
		)
	case *types.Slice:
		condition.param = condition.param.Index()
		condition.generateCodeForSlice(
			object,
			field, fieldType.Elem(),
		)
	default:
		log.Printf("struct field type not handled: %T", fieldType)
	}
}

func (condition *Condition) generateCodeForSlice(object types.Object, field Field, elemType types.Type) {
	switch elemTypeTyped := elemType.(type) {
	case *types.Basic:
		// una list de strings o algo asi,
		// por el momento solo anda con []byte porque el resto gorm no lo sabe encodear
		condition.generateCode(
			object,
			field.ChangeType(elemTypeTyped),
		)
	case *types.Named:
		elemObject := elemTypeTyped.Obj()
		// inverse relation condition
		_, err := getBadORMModelStruct(elemObject)
		if err == nil {
			log.Println(elemObject.Name())
			condition.codes = []jen.Code{
				generateOppositeJoinCondition(
					object,
					field,
					elemObject,
				),
			}
		}
	case *types.Pointer:
		condition.param = condition.param.Op("*")
		// slice de pointers, solo testeado temporalmente porque despues gorm no lo soporta
		condition.generateCodeForSlice(
			object,
			field, elemTypeTyped.Elem(),
		)
	default:
		log.Printf("struct field list elem type not handled: %T", elemTypeTyped)
	}
}

func generateConditionsForNamedType(object types.Object, field Field, fieldType *types.Named, param *jen.Statement) []jen.Code {
	// TODO quizas aca se puede eliminar el fieldType
	fieldObject := fieldType.Obj()
	// TODO esta linea de aca quedo rara
	_, err := getBadORMModelStruct(fieldObject)
	log.Println(err)

	if err == nil {
		objectStruct, err := getBadORMModelStruct(object)
		if err != nil {
			// TODO ver esto
			return []jen.Code{}
		}
		// TODO que pasa si esta en otro package? se importa solo?
		fields, err := getFields(
			objectStruct,
			// TODO testear esto si esta bien aca
			field.Tags.getEmbeddedPrefix(),
		)
		if err != nil {
			// TODO ver esto
			return []jen.Code{}
		}
		thisEntityHasTheFK := pie.Any(fields, func(otherField Field) bool {
			return otherField.Name == field.getJoinFromColumn()
		})

		log.Println(field.getJoinFromColumn())
		log.Println(thisEntityHasTheFK)

		if thisEntityHasTheFK {
			// belongsTo relation
			return []jen.Code{
				generateJoinCondition(
					object,
					field,
				),
			}
		}

		// hasOne or hasMany relation
		inverseJoinCondition := generateInverseJoinCondition(
			object,
			field, fieldObject,
		)

		return []jen.Code{
			inverseJoinCondition,
			generateOppositeJoinCondition(
				object,
				field,
				fieldObject,
			),
		}

		// TODO DeletedAt
	} else if (isGormCustomType(fieldType) || fieldType.String() == "time.Time") && fieldType.String() != "gorm.io/gorm.DeletedAt" {
		return []jen.Code{
			generateWhereCondition(
				object,
				field,
				param.Clone().Qual(
					getRelativePackagePath(fieldObject.Pkg()),
					fieldObject.Name(),
				),
			),
		}
	}

	log.Printf("struct field type not handled: %s", fieldType.String())
	return []jen.Code{}
}

var scanMethod = regexp.MustCompile(`func \(\*.*\)\.Scan\([a-zA-Z0-9_-]* interface\{\}\) error$`)
var valueMethod = regexp.MustCompile(`func \(.*\)\.Value\(\) \(database/sql/driver\.Value\, error\)$`)

func isGormCustomType(typeNamed *types.Named) bool {
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

var typeKindToJenStatement = map[types.BasicKind]func(*jen.Statement) *jen.Statement{
	types.Bool:       func(param *jen.Statement) *jen.Statement { return param.Bool() },
	types.Int:        func(param *jen.Statement) *jen.Statement { return param.Int() },
	types.Int8:       func(param *jen.Statement) *jen.Statement { return param.Int8() },
	types.Int16:      func(param *jen.Statement) *jen.Statement { return param.Int16() },
	types.Int32:      func(param *jen.Statement) *jen.Statement { return param.Int32() },
	types.Int64:      func(param *jen.Statement) *jen.Statement { return param.Int64() },
	types.Uint:       func(param *jen.Statement) *jen.Statement { return param.Uint() },
	types.Uint8:      func(param *jen.Statement) *jen.Statement { return param.Uint8() },
	types.Uint16:     func(param *jen.Statement) *jen.Statement { return param.Uint16() },
	types.Uint32:     func(param *jen.Statement) *jen.Statement { return param.Uint32() },
	types.Uint64:     func(param *jen.Statement) *jen.Statement { return param.Uint64() },
	types.Uintptr:    func(param *jen.Statement) *jen.Statement { return param.Uintptr() },
	types.Float32:    func(param *jen.Statement) *jen.Statement { return param.Float32() },
	types.Float64:    func(param *jen.Statement) *jen.Statement { return param.Float64() },
	types.Complex64:  func(param *jen.Statement) *jen.Statement { return param.Complex64() },
	types.Complex128: func(param *jen.Statement) *jen.Statement { return param.Complex128() },
	types.String:     func(param *jen.Statement) *jen.Statement { return param.String() },
}

func generateWhereCondition(object types.Object, field Field, param *jen.Statement) *jen.Statement {
	whereCondition := jen.Qual(
		badORMPath, badORMWhereCondition,
	).Types(
		jen.Qual(
			getRelativePackagePath(object.Pkg()),
			object.Name(),
		),
	)

	return jen.Func().Id(
		getConditionName(object, field.Name),
	).Params(
		param,
	).Add(
		whereCondition.Clone(),
	).Block(
		jen.Return(
			whereCondition.Clone().Values(jen.Dict{
				jen.Id("Field"): jen.Lit(field.getColumnName()),
				jen.Id("Value"): jen.Id("v"),
			}),
		),
	)
}

func generateOppositeJoinCondition(object types.Object, field Field, fieldObject types.Object) *jen.Statement {
	return generateJoinCondition(
		fieldObject,
		// TODO testear los Override Foreign Key
		Field{
			Name:   object.Name(),
			Object: object,
			Tags:   field.Tags,
		},
	)
}

func generateJoinCondition(object types.Object, field Field) *jen.Statement {
	log.Println(field.Object.Name())

	t1 := jen.Qual(
		getRelativePackagePath(object.Pkg()),
		object.Name(),
	)

	// TODO field.Type.Name me da lo mismo que field.Name
	t2 := jen.Qual(
		getRelativePackagePath(field.Object.Pkg()),
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

	return jen.Func().Id(
		getConditionName(object, field.Name),
	).Params(
		jen.Id("conditions").Op("...").Add(badormT2Condition),
	).Add(
		badormT1Condition,
	).Block(
		jen.Return(
			badormJoinCondition.Values(jen.Dict{
				jen.Id("T1Field"):    jen.Lit(strcase.ToSnake(field.getJoinFromColumn())),
				jen.Id("T2Field"):    jen.Lit(strcase.ToSnake(field.getJoinToColumn())),
				jen.Id("Conditions"): jen.Id("conditions"),
			}),
		),
	)
}

// TODO codigo duplicado
// TODO probablemente se puede hacer con el mismo metodo pero con el orden inverso
func generateInverseJoinCondition(object types.Object, field Field, fieldObject types.Object) *jen.Statement {
	log.Println(fieldObject.String())

	t1 := jen.Qual(
		getRelativePackagePath(object.Pkg()),
		object.Name(),
	)

	t2 := jen.Qual(
		getRelativePackagePath(fieldObject.Pkg()),
		fieldObject.Name(),
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

	return jen.Func().Id(
		getConditionName(object, field.Name),
	).Params(
		jen.Id("conditions").Op("...").Add(badormT2Condition),
	).Add(
		badormT1Condition,
	).Block(
		jen.Return(
			badormJoinCondition.Values(jen.Dict{
				jen.Id("T1Field"):    jen.Lit(strcase.ToSnake(field.getJoinToColumn())),
				jen.Id("T2Field"):    jen.Lit(strcase.ToSnake(field.NoSePonerNombre(object.Name()))),
				jen.Id("Conditions"): jen.Id("conditions"),
			}),
		),
	)
}

func getConditionName(object types.Object, fieldName string) string {
	return strcase.ToPascal(object.Name()) + strcase.ToPascal(fieldName) + badORMCondition
}

// TODO testear esto
func getRelativePackagePath(srcPkg *types.Package) string {
	if srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}
