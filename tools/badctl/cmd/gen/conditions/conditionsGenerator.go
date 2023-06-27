package conditions

import (
	"errors"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/tools/badctl/cmd/cmderrors"
	"github.com/ditrit/badaas/tools/badctl/cmd/log"
)

const badORMNewPreloadCondition = "NewPreloadCondition"

type ConditionsGenerator struct {
	object     types.Object
	objectType Type
}

func NewConditionsGenerator(object types.Object) *ConditionsGenerator {
	return &ConditionsGenerator{
		object:     object,
		objectType: Type{object.Type()},
	}
}

// Add conditions for an object in the file
func (cg ConditionsGenerator) AddConditions(file *File) error {
	fields, err := getFields(cg.objectType)
	if err != nil {
		return err
	}

	log.Logger.Infof("Generating conditions for type %q in %s", cg.object.Name(), file.name)

	cg.addConditionsForEachField(file, fields)

	return nil
}

// Add one condition for each field of the object
func (cg ConditionsGenerator) addConditionsForEachField(file *File, fields []Field) {
	conditions := cg.generateConditionsForEachField(file, fields)

	objectName := cg.object.Name()
	objectQual := jen.Qual(
		getRelativePackagePath(file.destPkg, cg.objectType),
		cg.objectType.Name(),
	)

	preloadAttributesCondition := jen.Var().Id(
		getPreloadAttributesName(objectName),
	).Op("=").Add(jen.Qual(
		badORMPath, badORMNewPreloadCondition,
	)).Types(
		objectQual,
	)
	fieldIdentifiers := []jen.Code{}

	preloadRelationsCondition := jen.Var().Id(
		objectName + "PreloadRelations",
	).Op("=").Index().Add(jen.Qual(
		badORMPath, badORMCondition,
	)).Types(
		objectQual,
	)
	relationPreloads := []jen.Code{}

	for _, condition := range conditions {
		file.Add(condition.codes...)

		// add all field names to the list of fields of the preload condition
		if condition.fieldIdentifier != "" {
			fieldIdentifiers = append(
				fieldIdentifiers,
				jen.Qual("", condition.fieldIdentifier),
			)
		}

		// add the preload to the list of all possible preloads
		if condition.preloadName != "" {
			relationPreloads = append(
				relationPreloads,
				jen.Qual("", condition.preloadName),
			)
		}
	}

	file.Add(preloadAttributesCondition.Call(fieldIdentifiers...))
	file.Add(preloadRelationsCondition.Values(relationPreloads...))
}

func getPreloadAttributesName(objectName string) string {
	return objectName + "PreloadAttributes"
}

// Generate the conditions for each of the object's fields
func (cg ConditionsGenerator) generateConditionsForEachField(file *File, fields []Field) []*Condition {
	conditions := []*Condition{}

	for _, field := range fields {
		log.Logger.Debugf("Generating condition for field %q", field.Name)

		if field.Embedded {
			conditions = append(conditions, cg.generateEmbeddedConditions(
				file,
				field,
			)...)
		} else {
			conditions = append(conditions, NewCondition(
				file.destPkg, cg.objectType, field,
			))
		}
	}

	return conditions
}

// Generate conditions for a embedded field
// it will generate a condition for each of the field of the embedded field's type
func (cg ConditionsGenerator) generateEmbeddedConditions(file *File, field Field) []*Condition {
	embeddedStructType, ok := field.Type.Underlying().(*types.Struct)
	if !ok {
		cmderrors.FailErr(errors.New("unreachable! embedded objects are always structs"))
	}

	fields, err := getStructFields(embeddedStructType)
	if err != nil {
		// embedded field's type has not fields
		return []*Condition{}
	}

	if !isBadORMBaseModel(field.TypeString()) {
		fields = pie.Map(fields, func(embeddedField Field) Field {
			embeddedField.ColumnPrefix = field.Tags.getEmbeddedPrefix()
			embeddedField.NamePrefix = field.Name

			return embeddedField
		})
	}

	return cg.generateConditionsForEachField(file, fields)
}
