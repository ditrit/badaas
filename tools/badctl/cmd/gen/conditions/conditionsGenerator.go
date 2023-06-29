package conditions

import (
	"go/types"

	"github.com/dave/jennifer/jen"

	"github.com/ditrit/badaas/tools/badctl/cmd/log"
)

const badORMNewPreloadCondition = "NewPreloadCondition"

type CodeConditionsGenerator struct {
	object     types.Object
	objectType Type
}

func NewConditionsGenerator(object types.Object) *CodeConditionsGenerator {
	return &CodeConditionsGenerator{
		object:     object,
		objectType: Type{object.Type()},
	}
}

// Add conditions for an object in the file
func (cg CodeConditionsGenerator) Into(file *File) error {
	fields, err := getFields(cg.objectType)
	if err != nil {
		return err
	}

	log.Logger.Infof("Generating conditions for type %q in %s", cg.object.Name(), file.name)

	cg.addConditionsForEachField(file, fields)

	return nil
}

// Add one condition for each field of the object
func (cg CodeConditionsGenerator) addConditionsForEachField(file *File, fields []Field) {
	conditions := cg.ForEachField(file, fields)

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

	if len(relationPreloads) > 0 {
		file.Add(preloadRelationsCondition.Values(relationPreloads...))
	}
}

func getPreloadAttributesName(objectName string) string {
	return objectName + "PreloadAttributes"
}

// Generate the conditions for each of the object's fields
func (cg CodeConditionsGenerator) ForEachField(file *File, fields []Field) []Condition {
	conditions := []Condition{}

	for _, field := range fields {
		log.Logger.Debugf("Generating condition for field %q", field.Name)

		if field.Embedded {
			conditions = append(
				conditions,
				generateForEmbeddedField[Condition](
					file,
					field,
					cg,
				)...,
			)
		} else {
			conditions = append(conditions, *NewCondition(
				file.destPkg, cg.objectType, field,
			))
		}
	}

	return conditions
}