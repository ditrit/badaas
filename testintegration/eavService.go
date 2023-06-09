package testintegration

import (
	"fmt"

	"github.com/elliotchance/pie/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services"
	"github.com/ditrit/badaas/utils"
)

type EAVServiceIntTestSuite struct {
	suite.Suite
	db               *gorm.DB
	eavService       services.EAVService
	entityRepository *repository.EntityRepository
	entityType1      *models.EntityType
	entityType2      *models.EntityType
}

func NewEAVServiceIntTestSuite(
	db *gorm.DB,
	eavService services.EAVService,
	entityRepository *repository.EntityRepository,
) *EAVServiceIntTestSuite {
	return &EAVServiceIntTestSuite{
		db:               db,
		eavService:       eavService,
		entityRepository: entityRepository,
	}
}

func (ts *EAVServiceIntTestSuite) TearDownSuite() {
	ts.cleanDB()
}

func (ts *EAVServiceIntTestSuite) SetupTest() {
	ts.cleanDB()

	ts.entityType1 = ts.createEntityType("entityType1", nil)
	ts.entityType2 = ts.createEntityType("entityType2", ts.entityType1)
}

func (ts *EAVServiceIntTestSuite) cleanDB() {
	CleanDBTables(
		ts.db,
		[]any{
			models.Value{},
			models.Attribute{},
			models.Entity{},
			models.EntityType{},
		},
	)
}

// ------------------------- GetEntity --------------------------------

func (ts *EAVServiceIntTestSuite) TestGetEntityReturnsErrorIfEntityDoesNotExist() {
	_, err := ts.eavService.GetEntity(ts.entityType1.Name, badorm.NewUUID())
	ts.ErrorContains(err, "record not found")
}

func (ts *EAVServiceIntTestSuite) TestGetEntityReturnsErrorIfEntityTypeDoesNotMatch() {
	otherEntity1 := ts.createEntity(ts.entityType1, map[string]any{})
	_, err := ts.eavService.GetEntity(ts.entityType2.Name, otherEntity1.ID)
	ts.ErrorContains(err, "record not found")
}

func (ts *EAVServiceIntTestSuite) TestGetEntityWorksIfEntityTypeMatch() {
	entity1 := ts.createEntity(ts.entityType2, map[string]any{})

	entityReturned, err := ts.eavService.GetEntity(ts.entityType2.Name, entity1.ID)
	ts.Nil(err)
	EqualEntity(&ts.Suite, entity1, entityReturned)
}

// ------------------------- GetEntities --------------------------------

func (ts *EAVServiceIntTestSuite) TestGetEntitiesOfNotExistentTypeReturnsError() {
	_, err := ts.eavService.GetEntities("not-exists", map[string]any{})
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, map[string]any{})
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createEntity(ts.entityType2, map[string]any{})

	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, make(map[string]any))
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheListWhenMultipleCreated() {
	match1 := ts.createEntity(ts.entityType2, map[string]any{})
	match2 := ts.createEntity(ts.entityType2, map[string]any{})
	match3 := ts.createEntity(ts.entityType2, map[string]any{})

	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, make(map[string]any))
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match1, match2, match3}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNotEntitiesCreated() {
	params := map[string]any{
		"string": "not_created",
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNothingMatch() {
	ts.createEntity(ts.entityType2, map[string]any{
		"string": "something_else",
	})

	params := map[string]any{
		"string": "not_match",
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "something_else",
	})

	params := map[string]any{
		"string": "match",
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
	})
	match2 := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "something_else",
	})

	params := map[string]any{
		"string": "match",
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match1, match2}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionThatDoesNotExistReturnsEmpty() {
	ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
	})
	ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
	})
	ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
	})

	params := map[string]any{
		"not_exists": "not_exists",
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionOfIntType() {
	match := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
		"int":    1,
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "not_match",
		"int":    2,
	})

	params := map[string]any{
		"int": 1.0,
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionOfIncorrectTypeReturnsEmptyList() {
	ts.createEntity(ts.entityType2, map[string]any{
		"string": "not_match",
		"int":    1,
	})

	params := map[string]any{
		"int": "not_an_int",
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)
	ts.Len(entities, 0)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionOfFloatType() {
	match := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
		"float":  1.1,
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "not_match",
		"float":  2.0,
	})

	params := map[string]any{
		"float": 1.1,
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionOfBoolType() {
	match := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
		"bool":   true,
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "not_match",
		"bool":   false,
	})

	params := map[string]any{
		"bool": true,
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionOfRelationType() {
	otherEntity1 := ts.createEntity(ts.entityType1, map[string]any{})
	otherEntity2 := ts.createEntity(ts.entityType1, map[string]any{})

	match := ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity1.ID.String(),
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity2.ID.String(),
	})

	params := map[string]any{
		"relation": otherEntity1.ID.String(),
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionFilterByNull() {
	match := ts.createEntity(ts.entityType2, map[string]any{
		"string": nil,
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "something",
	})

	params := map[string]any{
		"string": nil,
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
		"int":    1,
		"bool":   true,
	})
	match2 := ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
		"int":    1,
		"bool":   true,
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"string": "not_match",
		"int":    1,
		"bool":   true,
	})
	ts.createEntity(ts.entityType2, map[string]any{
		"string": "match",
		"int":    2,
		"bool":   true,
	})

	params := map[string]any{
		"string": "match",
		"int":    1.0,
		"bool":   true,
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match1, match2}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionThatJoins() {
	otherEntity1 := ts.createEntity(ts.entityType1, map[string]any{
		"int": 1,
	})
	otherEntity2 := ts.createEntity(ts.entityType1, map[string]any{
		"int": 2,
	})

	match := ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity1.ID.String(),
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity2.ID.String(),
	})

	params := map[string]any{
		"relation": map[string]any{
			"int": 1.0,
		},
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnDifferentAttributes() {
	otherEntity1 := ts.createEntity(ts.entityType1, map[string]any{
		"int":    1,
		"string": "match",
	})
	otherEntity2 := ts.createEntity(ts.entityType1, map[string]any{
		"int":    2,
		"string": "match",
	})

	match := ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity1.ID.String(),
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity2.ID.String(),
	})

	params := map[string]any{
		"relation": map[string]any{
			"int":    1.0,
			"string": "match",
		},
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsDifferentEntities() {
	entityType3 := ts.createEntityType("entityType3", nil)

	otherEntity11 := ts.createEntity(ts.entityType1, map[string]any{
		"int": 1,
	})
	otherEntity12 := ts.createEntity(ts.entityType1, map[string]any{
		"int": 2,
	})

	otherEntity31 := ts.createEntity(entityType3, map[string]any{
		"int": 3,
	})
	otherEntity32 := ts.createEntity(entityType3, map[string]any{
		"int": 4,
	})

	relation3Attr := models.NewRelationAttribute(
		ts.entityType2, "relation2",
		false, false, entityType3,
	)

	ts.addAttributeToEntityType(ts.entityType2, relation3Attr)

	match := ts.createEntity(ts.entityType2, map[string]any{
		"relation":  otherEntity11.ID.String(),
		"relation2": otherEntity31.ID.String(),
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"relation":  otherEntity12.ID.String(),
		"relation2": otherEntity32.ID.String(),
	})

	params := map[string]any{
		"relation": map[string]any{
			"int": 1.0,
		},
		"relation2": map[string]any{
			"int": 3.0,
		},
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsMultipleTimes() {
	entityType3 := ts.createEntityType("entityType3", nil)

	ts.addAttributeToEntityType(ts.entityType1, models.NewRelationAttribute(
		ts.entityType1, "relation",
		false, false, entityType3,
	))

	otherEntity31 := ts.createEntity(entityType3, map[string]any{
		"int": 3,
	})
	otherEntity32 := ts.createEntity(entityType3, map[string]any{
		"int": 4,
	})

	otherEntity11 := ts.createEntity(ts.entityType1, map[string]any{
		"int":      1,
		"relation": otherEntity31.ID.String(),
	})
	otherEntity12 := ts.createEntity(ts.entityType1, map[string]any{
		"int":      2,
		"relation": otherEntity32.ID.String(),
	})

	match := ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity11.ID.String(),
	})

	ts.createEntity(ts.entityType2, map[string]any{
		"relation": otherEntity12.ID.String(),
	})

	params := map[string]any{
		"relation": map[string]any{
			"int": 1.0,
			"relation": map[string]any{
				"int": 3.0,
			},
		},
	}
	entities, err := ts.eavService.GetEntities(ts.entityType2.Name, params)
	ts.Nil(err)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

// ------------------------- CreateEntity --------------------------------

func (ts *EAVServiceIntTestSuite) TestCreateMultipleEntitiesDoesNotGenerateGormBug() {
	stringAttr := utils.FindFirst(ts.entityType1.Attributes, func(attr *models.Attribute) bool {
		return attr.ValueType == models.StringValueType
	})

	initialDisplayNameID := (*stringAttr).ID

	stringAttr2 := &models.Attribute{
		Name:      "string2",
		ValueType: models.StringValueType,
	}
	ts.addAttributeToEntityType(ts.entityType1, stringAttr2)

	initialDescriptionID := stringAttr2.ID

	for i := 0; i < 10; i++ {
		params := map[string]any{
			"string":      fmt.Sprintf("displayName%d", i),
			"description": fmt.Sprintf("description%d", i),
		}
		entity, err := ts.eavService.CreateEntity(ts.entityType1.Name, params)
		ts.Nil(err)

		for _, value := range entity.Fields {
			if value.Attribute.Name == "string" {
				ts.Equal(initialDisplayNameID, value.AttributeID)
			} else if value.Attribute.Name == "string2" {
				ts.Equal(initialDescriptionID, value.AttributeID)
			}
		}
	}
}

func (ts *EAVServiceIntTestSuite) TestCreateReturnsErrorIfEntityTypeDoesNotExist() {
	_, err := ts.eavService.CreateEntity("not-exists", map[string]any{})
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestCreateReturnsErrorIfTheTypeOfAValueIsUnsupported() {
	params := map[string]any{
		"string": []string{"salut", "bonjour"},
	}
	entity, err := ts.eavService.CreateEntity(ts.entityType1.Name, params)
	ts.Nil(entity)
	ts.ErrorContains(err, "unsupported type")
}

func (ts *EAVServiceIntTestSuite) TestCreateReturnsErrorIfUUIDCantBeParsed() {
	params := map[string]any{
		"relation": "not-a-uuid",
	}
	entity, err := ts.eavService.CreateEntity(ts.entityType2.Name, params)
	ts.Nil(entity)
	ts.ErrorIs(err, services.ErrCantParseUUID)
}

func (ts *EAVServiceIntTestSuite) TestCreateReturnsErrorIfRelationAttributePointsToNotExistentType() {
	ts.addAttributeToEntityType(ts.entityType2, &models.Attribute{
		Name:                       "relation2",
		ValueType:                  models.RelationValueType,
		RelationTargetEntityTypeID: badorm.NewUUID(),
	})

	params := map[string]any{
		"relation2": badorm.NewUUID().String(),
	}
	entity, err := ts.eavService.CreateEntity(ts.entityType2.Name, params)
	ts.Nil(entity)
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestCreatesDefaultAttributes() {
	ts.addAttributeToEntityType(ts.entityType2, &models.Attribute{
		Name:       "default",
		ValueType:  models.IntValueType,
		Default:    true,
		DefaultInt: 1,
	})

	entity, err := ts.eavService.CreateEntity(ts.entityType2.Name, map[string]any{})
	ts.Nil(err)
	ts.Len(entity.Fields, 6)
	notNull := pie.Filter(entity.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 1)
	ts.Equal(1, notNull[0].Value())
}

func (ts *EAVServiceIntTestSuite) TestCreatesWithoutRequiredValueRespondsError() {
	ts.addAttributeToEntityType(ts.entityType2, &models.Attribute{
		Name:      "required",
		Required:  true,
		ValueType: models.StringValueType,
	})

	entity, err := ts.eavService.CreateEntity(ts.entityType2.Name, map[string]any{})
	ts.Nil(entity)
	ts.ErrorContains(err, "field required is missing and is required")
}

func (ts *EAVServiceIntTestSuite) TestCreatesIntAttributeEvenIfItIsInFloatFormat() {
	entity, err := ts.eavService.CreateEntity(ts.entityType2.Name, map[string]any{
		"int": 2.0,
	})
	ts.Nil(err)
	ts.Len(entity.Fields, 5)
	notNull := pie.Filter(entity.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 1)
	ts.Equal(2, notNull[0].Value())
}

// // ------------------------- UpdateEntity --------------------------------

func (ts *EAVServiceIntTestSuite) TestUpdateEntityMultipleTimesDoesNotGenerateGormBug() {
	stringAttr := utils.FindFirst(ts.entityType1.Attributes, func(attr *models.Attribute) bool {
		return attr.ValueType == models.StringValueType
	})

	initialDisplayNameID := (*stringAttr).ID

	stringAttr2 := &models.Attribute{
		Name:      "string2",
		ValueType: models.StringValueType,
	}
	ts.addAttributeToEntityType(ts.entityType1, stringAttr2)
	initialDescriptionID := stringAttr2.ID

	entity := ts.createEntity(ts.entityType1, map[string]any{
		"string":  "displayName",
		"string2": "description",
	})

	for i := 0; i < 10; i++ {
		params := map[string]any{
			"string":  fmt.Sprintf("displayName%d", i),
			"string2": fmt.Sprintf("description%d", i),
		}
		entity, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, params)
		ts.Nil(err)

		for _, value := range entity.Fields {
			if value.Attribute.Name == "string" {
				ts.Equal(initialDisplayNameID, value.AttributeID)
			} else if value.Attribute.Name == "string2" {
				ts.Equal(initialDescriptionID, value.AttributeID)
			}
		}
	}
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfEntityDoesNotExist() {
	_, err := ts.eavService.UpdateEntity(ts.entityType1.Name, badorm.NewUUID(), map[string]any{})
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityWorksForAllTheTypes() {
	otherEntity1 := ts.createEntity(ts.entityType1, map[string]any{})
	entity := ts.createEntity(ts.entityType2, map[string]any{
		"string": "displayName",
	})

	paramsUpdate := map[string]any{
		"string":   nil,
		"int":      1,
		"float":    1.1,
		"bool":     true,
		"relation": otherEntity1.ID.String(),
	}
	entity, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.Nil(err)
	ts.Len(entity.Fields, 5)
	notNull := pie.Filter(entity.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 4)
	values := pie.Map(notNull, func(v *models.Value) any {
		return v.Value()
	})
	ts.Contains(values, 1)
	ts.Contains(values, 1.1)
	ts.Contains(values, true)
	ts.Contains(values, otherEntity1.ID)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfStringForIntType() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"int": "1",
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfStringForFloatType() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"float": "1",
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfStringForBoolType() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"bool": "1",
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfIntForStringType() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"string": 1,
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfUUIDCantBeParsedForRelationType() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"relation": "not-uuid",
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, services.ErrCantParseUUID)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfUUIDDoesNotExists() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"relation": badorm.NewUUID().String(),
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfUUIDDoesNotCorrespondsToTheRelationEntityType() {
	otherEntityType2 := &models.EntityType{
		Name: "other2",
	}

	err := ts.db.Create(otherEntityType2).Error
	ts.Nil(err)

	entity := ts.createEntity(ts.entityType2, map[string]any{})
	entityOther2 := ts.createEntity(otherEntityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"relation": entityOther2.ID.String(),
	}
	_, err = ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityDoesNotUpdateAValueIfOtherFails() {
	entity := ts.createEntity(ts.entityType2, map[string]any{})

	paramsUpdate := map[string]any{
		"string": "something",
		"int":    "1",
	}
	_, err := ts.eavService.UpdateEntity(entity.EntityType.Name, entity.ID, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)

	entityReturned, err := ts.eavService.GetEntity(ts.entityType2.Name, entity.ID)
	ts.Nil(err)

	notNull := pie.Filter(entityReturned.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 0)
}

// ------------------------- DeleteEntity -------------------------

func (ts *EAVServiceIntTestSuite) TestDeleteEntityReturnsErrorIfEntityDoesNotExist() {
	err := ts.eavService.DeleteEntity(ts.entityType2.Name, badorm.NewUUID())
	ts.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (ts *EAVServiceIntTestSuite) TestDeleteEntityReturnsErrorIfEntityTypeDoesNotMatch() {
	entity1 := ts.createEntity(ts.entityType2, map[string]any{
		"string": "displayName",
	})

	err := ts.eavService.DeleteEntity(ts.entityType1.Name, entity1.ID)
	ts.ErrorIs(err, gorm.ErrRecordNotFound)

	var values []models.Value
	err = ts.db.Find(&values).Error
	ts.Nil(err)
	ts.Len(values, 5)
}

func (ts *EAVServiceIntTestSuite) TestDeleteWorks() {
	entity1 := ts.createEntity(ts.entityType2, map[string]any{
		"string": "displayName",
	})

	err := ts.eavService.DeleteEntity(entity1.EntityType.Name, entity1.ID)
	ts.Nil(err)

	var profiles []models.Entity
	err = ts.db.Find(&profiles).Error
	ts.Nil(err)
	ts.Len(profiles, 0)

	var values []models.Value
	err = ts.db.Find(&values).Error
	ts.Nil(err)
	ts.Len(values, 0)
}

// ------------------------- utils -------------------------

func (ts *EAVServiceIntTestSuite) createEntityType(name string, relationEntityType *models.EntityType) *models.EntityType {
	entityType := &models.EntityType{
		Name: name,
	}

	entityType.Attributes = []*models.Attribute{
		{
			EntityTypeID: entityType.ID,
			Name:         "int",
			ValueType:    models.IntValueType,
		},
		{
			EntityTypeID: entityType.ID,
			Name:         "string",
			ValueType:    models.StringValueType,
		},
		{
			EntityTypeID: entityType.ID,
			Name:         "bool",
			ValueType:    models.BooleanValueType,
		},
		{
			EntityTypeID: entityType.ID,
			Name:         "float",
			ValueType:    models.FloatValueType,
		},
	}

	if relationEntityType != nil {
		entityType.Attributes = append(entityType.Attributes, models.NewRelationAttribute(
			entityType, "relation",
			false, false, relationEntityType,
		))
	}

	err := ts.db.Create(&entityType).Error
	ts.Nil(err)

	return entityType
}

func (ts *EAVServiceIntTestSuite) addAttributeToEntityType(entityType *models.EntityType, attribute *models.Attribute) {
	attribute.EntityTypeID = entityType.ID
	entityType.Attributes = append(entityType.Attributes, attribute)

	err := ts.db.Save(&entityType).Error
	ts.Nil(err)
}

func (ts *EAVServiceIntTestSuite) createEntity(entityType *models.EntityType, values map[string]any) *models.Entity {
	entity, err := ts.eavService.CreateEntity(entityType.Name, values)
	ts.Nil(err)

	return entity
}
