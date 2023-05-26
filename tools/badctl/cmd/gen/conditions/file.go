package conditions

import (
	"errors"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
	"github.com/ditrit/badaas/tools/badctl/cmd/log"
	"github.com/ditrit/badaas/tools/badctl/cmd/version"
)

type File struct {
	destPkg string
	jenFile *jen.File
	name    string
}

func NewConditionsFile(destPkg string, name string) *File {
	// Start a new file in destination package
	f := jen.NewFile(destPkg)

	// Add a package comment, so IDEs detect files as generated
	f.PackageComment("Code generated by badctl v" + version.Version + ", DO NOT EDIT.")

	return &File{
		destPkg: destPkg,
		name:    name,
		jenFile: f,
	}
}

func (file File) AddConditionsFor(object types.Object) error {
	fields, err := getFields(Type{object.Type()}, "")
	if err != nil {
		return err
	}

	log.Logger.Infof("Generating conditions for type %q in %s", object.Name(), file.name)

	file.addConditionsForEachField(object, fields)
	return nil
}

func (file File) addConditionsForEachField(object types.Object, fields []Field) {
	conditions := file.generateConditionsForEachField(object, fields)

	for _, condition := range conditions {
		// TODO esto no me gusta mucho que este aca
		for _, code := range condition.codes {
			file.jenFile.Add(code)
		}
	}
}

// Write generated file
func (file File) Save() error {
	return file.jenFile.Save(file.name)
}

// badorm/baseModels.go
var badORMModels = []string{"UUIDModel", "UIntModel"}

const (
	badORMPath = "github.com/ditrit/badaas/badorm"
	// badorm/condition.go
	badORMCondition      = "Condition"
	badORMWhereCondition = "WhereCondition"
	badORMJoinCondition  = "JoinCondition"
)

func (file File) generateConditionsForEachField(object types.Object, fields []Field) []*Condition {
	conditions := []*Condition{}
	for _, field := range fields {
		log.Logger.Debugf("Generating condition for field %q", field.Name)
		if field.Embedded {
			conditions = append(conditions, file.generateEmbeddedConditions(
				object,
				field,
			)...)
		} else {
			newCondition, err := NewCondition(
				file.destPkg, Type{object.Type()}, field,
			)
			if err != nil {
				cmderrors.FailErr(err)
			}

			conditions = append(conditions, newCondition)
		}
	}

	return conditions
}

// TODO quizas esto no deberia estar aca
func (file File) generateEmbeddedConditions(object types.Object, field Field) []*Condition {
	embeddedStructType, ok := field.Type.Underlying().(*types.Struct)
	if !ok {
		cmderrors.FailErr(errors.New("unreachable! embedded objects are always structs"))
	}

	fields, err := getStructFields(embeddedStructType, field.Tags.getEmbeddedPrefix())
	if err != nil {
		// TODO ver esto
		return []*Condition{}
	}

	return file.generateConditionsForEachField(object, fields)
}
