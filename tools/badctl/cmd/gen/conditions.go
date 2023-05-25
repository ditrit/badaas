package gen

import (
	"errors"
	"fmt"
	"go/types"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/ditrit/badaas/tools/badctl/cmd/version"
	"github.com/ditrit/verdeter"
	"github.com/elliotchance/pie/v2"
	"github.com/ettle/strcase"
	"github.com/spf13/cobra"

	"golang.org/x/tools/go/packages"
)

var genConditionsCmd = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "conditions",
	Short: "Generate conditions to query your objects using BaDORM",
	Long:  `gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.`,
	Run:   generateConditions,
	Args:  cobra.MinimumNArgs(1),
})

// badorm/baseModels.go
var badORMModels = []string{"UUIDModel", "UIntModel"}

const (
	badORMPath = "github.com/ditrit/badaas/badorm"
	// badorm/condition.go
	badORMCondition      = "Condition"
	badORMWhereCondition = "WhereCondition"
	badORMJoinCondition  = "JoinCondition"
)

func generateConditions(cmd *cobra.Command, args []string) {
	// Inspect package and use type checker to infer imported types
	pkgs := loadPackages(args)

	// Get the package of the file with go:generate comment
	destPkg := os.Getenv("GOPACKAGE")
	if destPkg == "" {
		// TODO que tambien se pueda usar solo
		failErr(errors.New("this command should be called using go generate"))
	}
	log.Println(destPkg)

	for _, pkg := range pkgs {
		log.Println(pkg.Types.Path())
		log.Println(pkg.Types.Name())

		for _, name := range pkg.Types.Scope().Names() {
			object := getObject(pkg, name)
			if object != nil {
				log.Println(name)

				// Generate code using jennifer
				err := generateConditionsFile(
					destPkg,
					object,
				)
				if err != nil {
					failErr(err)
				}
			}
		}
	}
}

func getObject(pkg *packages.Package, name string) types.Object {
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			name, pkg))
	}

	// Generate only if it is a declared type
	object, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}

	return object
}

// TODO add logs

func generateConditionsFile(destPkg string, object types.Object) error {
	// Generate only when underlying type is a struct
	// (ignore const, var, func, etc.)
	structType := getBadORMModelStruct(object)
	if structType == nil {
		return nil
	}

	fields := getFields(structType, "")
	if len(fields) == 0 {
		return nil
	}

	// Start a new file in destination package
	f := jen.NewFile(destPkg)

	// Add a package comment, so IDEs detect files as generated
	f.PackageComment("Code generated by badctl v" + version.Version + ", DO NOT EDIT.")

	addConditionForEachField(f, fields, destPkg, object)

	// Write generated file
	return f.Save(strings.ToLower(object.Name()) + "_conditions.go")
}

func getBadORMModelStruct(object types.Object) *types.Struct {
	structType, ok := object.Type().Underlying().(*types.Struct)
	if !ok || !isBadORMModel(structType) {
		return nil
	}

	return structType
}

func isBadORMModel(structType *types.Struct) bool {
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		if field.Embedded() && pie.Contains(badORMModels, field.Name()) {
			return true
		}
	}

	return false
}

func addConditionForEachField(f *jen.File, fields []Field, destPkg string, object types.Object) {
	for _, field := range fields {
		log.Println(field.Name)
		if field.Embedded {
			addEmbeddedConditions(
				f,
				destPkg,
				object,
				field,
			)
		} else {
			addConditionsForField(
				f,
				destPkg,
				object,
				field,
			)
		}
	}
}

func addEmbeddedConditions(f *jen.File, destPkg string, object types.Object, field Field) {
	embededFieldType, ok := field.Type.Type().(*types.Named)
	if !ok {
		failErr(errors.New("unreacheable! embeddeds are allways of type Named"))
	}
	embededStructType, ok := embededFieldType.Underlying().(*types.Struct)
	if !ok {
		failErr(errors.New("unreacheable! embeddeds are allways structs"))
	}

	addConditionForEachField(
		f,
		getFields(embededStructType, field.Tags.getEmbeddedPrefix()),
		destPkg,
		object,
	)
}

var basicParam = jen.Id("v")

func addConditionsForField(f *jen.File, destPkg string, object types.Object, field Field) {
	conditions := generateConditionsForField(
		destPkg,
		object, field,
		basicParam.Clone(),
	)

	for _, condition := range conditions {
		f.Add(condition)
	}
}

func generateConditionsForField(destPkg string, object types.Object, field Field, param *jen.Statement) []jen.Code {
	switch fieldTypeTyped := field.Type.Type().(type) {
	case *types.Basic:
		return []jen.Code{
			generateWhereCondition(
				destPkg,
				object,
				field,
				typeKindToJenStatement[fieldTypeTyped.Kind()](param),
			),
		}
	case *types.Named:
		return generateConditionsForNamedType(
			destPkg, object,
			field, fieldTypeTyped,
			param,
		)
	case *types.Pointer:
		return generateConditionsForField(
			destPkg,
			object,
			field.ChangeType(fieldTypeTyped.Elem()),
			param.Clone().Op("*"),
		)
	case *types.Slice:
		return generateConditionForSlice(
			destPkg, object,
			field, fieldTypeTyped.Elem(),
			param.Clone().Index(),
		)
	default:
		log.Printf("struct field type not handled: %T", fieldTypeTyped)
	}

	// TODO ver este error
	return []jen.Code{}
}

func generateConditionsForNamedType(destPkg string, object types.Object, field Field, fieldType *types.Named, param *jen.Statement) []jen.Code {
	// TODO quizas aca se puede eliminar el fieldType
	fieldObject := fieldType.Obj()
	fieldModelStruct := getBadORMModelStruct(fieldObject)
	if fieldModelStruct != nil {
		// TODO que pasa si esta en otro package? se importa solo?
		fields := getFields(
			getBadORMModelStruct(object),
			// TODO testear esto
			field.Tags.getEmbeddedPrefix(),
		)
		thisEntityHasTheFK := pie.Any(fields, func(otherField Field) bool {
			return otherField.Name == field.getJoinFromColumn()
		})

		log.Println(field.getJoinFromColumn())
		log.Println(thisEntityHasTheFK)

		if thisEntityHasTheFK {
			// belongsTo relation
			return []jen.Code{
				generateJoinCondition(
					destPkg,
					object,
					field,
				),
			}
		}

		// hasOne or hasMany relation
		inverseJoinCondition := generateInverseJoinCondition(
			destPkg,
			object,
			field, fieldObject,
		)

		return []jen.Code{
			inverseJoinCondition,
			generateOppositeJoinCondition(
				destPkg,
				object,
				field,
				fieldObject,
			),
		}

		// TODO DeletedAt
	} else if (isGormCustomType(fieldType) || fieldType.String() == "time.Time") && fieldType.String() != "gorm.io/gorm.DeletedAt" {
		return []jen.Code{
			generateWhereCondition(
				destPkg,
				object,
				field,
				param.Clone().Qual(
					getRelativePackagePath(fieldObject.Pkg(), destPkg),
					fieldObject.Name(),
				),
			),
		}
	}

	log.Printf("struct field type not handled: %s", fieldType.String())
	return []jen.Code{}
}

func generateConditionForSlice(destPkg string, object types.Object, field Field, elemType types.Type, param *jen.Statement) []jen.Code {
	switch elemTypeTyped := elemType.(type) {
	case *types.Basic:
		// una list de strings o algo asi,
		// por el momento solo anda con []byte porque el resto gorm no lo sabe encodear
		return generateConditionsForField(
			destPkg,
			object,
			field.ChangeType(elemTypeTyped),
			param,
		)
	case *types.Named:
		elemTypeName := elemTypeTyped.Obj()
		// inverse relation condition
		if getBadORMModelStruct(elemTypeName) != nil {
			log.Println(elemTypeName.Name())
			return []jen.Code{
				generateOppositeJoinCondition(
					destPkg,
					object,
					field,
					elemTypeName,
				),
			}
		}
	case *types.Pointer:
		// slice de pointers, solo testeado temporalmente porque despues gorm no lo soporta
		return generateConditionForSlice(
			destPkg, object,
			field, elemTypeTyped.Elem(),
			param.Op("*"),
		)
	default:
		log.Printf("struct field list elem type not handled: %T", elemTypeTyped)
	}

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

func generateWhereCondition(destPkg string, object types.Object, field Field, param *jen.Statement) *jen.Statement {
	whereCondition := jen.Qual(
		badORMPath, badORMWhereCondition,
	).Types(
		jen.Qual(
			getRelativePackagePath(object.Pkg(), destPkg),
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

func generateOppositeJoinCondition(destPkg string, object types.Object, field Field, fieldObject types.Object) *jen.Statement {
	return generateJoinCondition(
		destPkg,
		fieldObject,
		// TODO testear los Override Foreign Key
		Field{
			Name: object.Name(),
			Type: object,
			Tags: field.Tags,
		},
	)
}

func generateJoinCondition(destPkg string, object types.Object, field Field) *jen.Statement {
	log.Println(field.Type.Name())

	t1 := jen.Qual(
		getRelativePackagePath(object.Pkg(), destPkg),
		object.Name(),
	)

	// TODO field.Type.Name me da lo mismo que field.Name
	t2 := jen.Qual(
		getRelativePackagePath(field.Type.Pkg(), destPkg),
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
func generateInverseJoinCondition(destPkg string, object types.Object, field Field, fieldObject types.Object) *jen.Statement {
	log.Println(fieldObject.String())

	t1 := jen.Qual(
		getRelativePackagePath(object.Pkg(), destPkg),
		object.Name(),
	)

	t2 := jen.Qual(
		getRelativePackagePath(fieldObject.Pkg(), destPkg),
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
func getRelativePackagePath(srcPkg *types.Package, destPkg string) string {
	if srcPkg.Name() == destPkg {
		return ""
	}

	return srcPkg.Path()
}

func loadPackages(paths []string) []*packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes}
	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		failErr(fmt.Errorf("loading packages for inspection: %v", err))
	}

	// print compilation errors of source packages
	packages.PrintErrors(pkgs)

	return pkgs
}

func failErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
