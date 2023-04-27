package repository

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/elliotchance/pie/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EntityRepository struct {
	logger          *zap.Logger
	valueRepository *ValueRepository
}

func NewEntityRepository(
	logger *zap.Logger,
	valueRepository *ValueRepository,
) *EntityRepository {
	return &EntityRepository{
		logger:          logger,
		valueRepository: valueRepository,
	}
}

// Get the Entity of type with name "entityTypeName" that has the "id"
func (r *EntityRepository) Get(tx *gorm.DB, entityTypeName string, id uuid.UUID) (*models.Entity, error) {
	var entity models.Entity

	query := tx.Preload("Fields").Preload("Fields.Attribute").Preload("EntityType")
	query = query.Joins(
		`JOIN entity_types ON
			entity_types.id = entities.entity_type_id`,
	)
	err := query.Where(
		map[string]any{"entities.id": id, "entity_types.name": entityTypeName},
	).First(&entity).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Creates an entity and its values in the database
// must be used in place of gorm's db.Save(entity) because of the bug
// when using gorm with cockroachDB. For more info refer to:
// https://github.com/FrancoLiberali/cockroachdb_gorm_bug
func (r *EntityRepository) Create(tx *gorm.DB, entity *models.Entity) error {
	now := time.Now()

	query, values, err := sq.Insert("entities").
		Columns("created_at", "updated_at", "entity_type_id").
		Values(now, now, entity.EntityType.ID).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return err
	}

	var result string
	err = tx.Raw(query, values...).Scan(&result).Error
	if err != nil {
		return err
	}

	uuid, err := uuid.Parse(result)
	if err != nil {
		return err
	}

	pie.Each(entity.Fields, func(value *models.Value) {
		value.EntityID = uuid
	})

	if len(entity.Fields) > 0 {
		err = r.valueRepository.Create(tx, entity.Fields)
		if err != nil {
			return err
		}
	}

	entity.ID = uuid
	return nil
}

// Adds to the "query" the verification that the value for "attribute" is "expectedValue"
func (r *EntityRepository) AddValueCheckToQuery(query *gorm.DB, attributeName string, expectedValue any) error {
	return r.addValueCheckToQueryInternal(query, attributeName, expectedValue, "")
}

// Adds to the "query" the verification that the value for "attribute" is "expectedValue"
func (r *EntityRepository) addValueCheckToQueryInternal(query *gorm.DB, attributeName string, expectedValue any, entitiesTableSuffix string) error {
	attributesSuffix := entitiesTableSuffix + "_" + attributeName
	stringQuery := fmt.Sprintf(
		`JOIN attributes attributes%[1]s ON
			attributes%[1]s.entity_type_id = entities%[2]s.entity_type_id
			AND attributes%[1]s.name = ?
		JOIN values values%[1]s ON
			values%[1]s.attribute_id = attributes%[1]s.id
			AND values%[1]s.entity_id = entities%[2]s.id
		`,
		attributesSuffix,
		entitiesTableSuffix,
	)
	switch expectedValueTyped := expectedValue.(type) {
	case float64:
		stringQuery += fmt.Sprintf(
			"AND ((%s) OR (%s))",
			getQueryCheckValueOfType(attributesSuffix, models.IntValueType),
			getQueryCheckValueOfType(attributesSuffix, models.FloatValueType),
		)
	case bool:
		stringQuery += "AND " + getQueryCheckValueOfType(attributesSuffix, models.BooleanValueType)
	case string:
		_, err := uuid.Parse(expectedValueTyped)
		if err == nil {
			stringQuery += "AND " + getQueryCheckValueOfType(attributesSuffix, models.RelationValueType)
		} else {
			stringQuery += "AND " + getQueryCheckValueOfType(attributesSuffix, models.StringValueType)
		}
	case nil:
		stringQuery += fmt.Sprintf(
			"AND values%s.is_null = true",
			attributesSuffix,
		)
	case map[string]any:
		return r.addJoinToQuery(
			query, attributeName, expectedValueTyped,
			attributesSuffix, stringQuery,
		)
	default:
		return fmt.Errorf("unsupported type")
	}

	query.Joins(stringQuery, attributeName, expectedValue, expectedValue)

	return nil
}

// Returns query string to check that the attribute is of type "valueType" and that the related value
// is the expected one
func getQueryCheckValueOfType(attributesSuffix string, valueType models.ValueTypeT) string {
	return fmt.Sprintf(
		"attributes%[1]s.value_type = '%[2]s' AND values%[1]s.%[2]s_val = ?",
		attributesSuffix, string(valueType),
	)
}

// Adds a join to the "query" by the "attributeName" that must be relation type
// then, adds the verification that the values for the joined entity are "expectedValues"

// "expectedValues" is in {"attributeName": expectedValue} format
func (r *EntityRepository) addJoinToQuery(
	query *gorm.DB, attributeName string, expectedValues map[string]any,
	attributesSuffix, stringQuery string,
) error {
	stringQuery += fmt.Sprintf(`
				AND attributes%[1]s.value_type = 'relation'
			JOIN entities entities%[1]s ON
				entities%[1]s.id = values%[1]s.relation_val
				AND entities%[1]s.deleted_at IS NULL
			`,
		attributesSuffix,
	)

	query.Joins(stringQuery, attributeName)

	var err error
	for joinEntityAttribute, joinEntityValue := range expectedValues {
		err = r.addValueCheckToQueryInternal(query, joinEntityAttribute, joinEntityValue, attributesSuffix)
		if err != nil {
			return err
		}
	}

	return nil
}
