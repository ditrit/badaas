package integrationtests

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/eavservice"
	"github.com/elliotchance/pie/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gotest.tools/assert"
)

type EAVServiceIntTestSuite struct {
	suite.Suite
	logger           *zap.Logger
	db               *gorm.DB
	eavService       eavservice.EAVService
	entityRepository *repository.EntityRepository
	profileType      *models.EntityType
	displayNameAttr  *models.Attribute
	descriptionAttr  *models.Attribute
}

func NewEAVServiceIntTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
	eavService eavservice.EAVService,
	entityRepository *repository.EntityRepository,
) *EAVServiceIntTestSuite {
	return &EAVServiceIntTestSuite{
		logger:           logger,
		db:               db,
		eavService:       eavService,
		entityRepository: entityRepository,
	}
}

func (ts *EAVServiceIntTestSuite) SetupTest() {
	SetupDB(ts.db)

	// CREATION OF THE PROFILE TYPE AND ASSOCIATED ATTRIBUTES
	ts.profileType = &models.EntityType{
		Name: "profile",
	}
	ts.displayNameAttr = &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "displayName",
		ValueType:    models.StringValueType,
		Required:     false,
	}
	ts.descriptionAttr = &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "description",
		ValueType:    models.StringValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		ts.displayNameAttr,
		ts.descriptionAttr,
	)

	err := ts.db.Create(&ts.profileType).Error
	ts.Nil(err)
}

// ------------------------- GetEntity --------------------------------

func (ts *EAVServiceIntTestSuite) TestGetEntityReturnsErrorIfEntityDoesNotExist() {
	_, err := ts.eavService.GetEntity(ts.profileType, uuid.New())
	ts.ErrorContains(err, "record not found")
}

func (ts *EAVServiceIntTestSuite) TestGetEntityReturnsErrorIfEntityTypeDoesNotMatch() {
	otherEntityType := &models.EntityType{
		Name: "other",
	}

	err := ts.db.Create(otherEntityType).Error
	ts.Nil(err)

	otherEntity1, err := ts.eavService.CreateEntity(otherEntityType, map[string]any{})
	ts.Nil(err)

	_, err = ts.eavService.GetEntity(ts.profileType, otherEntity1.ID)
	ts.ErrorIs(err, eavservice.ErrIDDontMatchEntityType)
}

func (ts *EAVServiceIntTestSuite) TestGetEntityWorksIfEntityTypeMatch() {
	entity1, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	entityReturned, err := ts.eavService.GetEntity(ts.profileType, entity1.ID)
	ts.Nil(err)
	EqualEntity(&ts.Suite, entity1, entityReturned)
}

// ------------------------- GetEntitiesWithParams --------------------------------

func (ts *EAVServiceIntTestSuite) TestWithoutParamsReturnsEmptyIfNotEntitiesCreated() {
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, make(map[string]string))

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithoutParamsReturnsTheOnlyOneIfOneEntityCreated() {
	profile := ts.createProfile(ts.profileType, "profile")

	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, make(map[string]string))

	EqualEntityList(&ts.Suite, []*models.Entity{profile}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithoutParamsReturnsTheListWhenMultipleCreated() {
	profile1 := ts.createProfile(ts.profileType, "profile1")
	profile2 := ts.createProfile(ts.profileType, "profile2")
	profile3 := ts.createProfile(ts.profileType, "profile3")

	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, make(map[string]string))

	EqualEntityList(&ts.Suite, []*models.Entity{profile1, profile2, profile3}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsEmptyIfNotEntitiesCreated() {
	params := map[string]string{
		"displayName": "not_created",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsEmptyIfNothingMatch() {
	ts.createProfile(ts.profileType, "profile")

	params := map[string]string{
		"displayName": "not_match",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsOneIfOnlyOneMatch() {
	matchProfile := ts.createProfile(ts.profileType, "match")
	ts.createProfile(ts.profileType, "something_else")

	params := map[string]string{
		"displayName": "match",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{matchProfile}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProfile(ts.profileType, "match")
	match2 := ts.createProfile(ts.profileType, "match")
	ts.createProfile(ts.profileType, "something_else")

	params := map[string]string{
		"displayName": "match",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match1, match2}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamThatDoesNotExistReturnsAllEntities() {
	match1 := ts.createProfile(ts.profileType, "match")
	match2 := ts.createProfile(ts.profileType, "match")
	match3 := ts.createProfile(ts.profileType, "match")

	params := map[string]string{
		"not_exists": "not_exists",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match1, match2, match3}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamOfIntType() {
	intAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "int",
		ValueType:    models.IntValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		intAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	match, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "match",
		"int":         1,
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "not_match",
		"int":         2,
	})
	ts.Nil(err)

	params := map[string]string{
		"int": "1",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamOfIntTypeThatIsNotAnIntReturnEmpty() {
	intAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "int",
		ValueType:    models.IntValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		intAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "match",
		"int":         1,
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "not_match",
		"int":         2,
	})
	ts.Nil(err)

	params := map[string]string{
		"int": "not_an_int",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamOfFloatType() {
	floatAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "float",
		ValueType:    models.FloatValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		floatAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	match, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "match",
		"float":       1.1,
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "not_match",
		"float":       2.0,
	})
	ts.Nil(err)

	params := map[string]string{
		"float": "1.1",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamOfBoolType() {
	boolAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "bool",
		ValueType:    models.BooleanValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		boolAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	match, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "match",
		"bool":        true,
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "not_match",
		"bool":        false,
	})
	ts.Nil(err)

	params := map[string]string{
		"bool": "true",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamOfRelationType() {
	otherEntityType := &models.EntityType{
		Name: "other",
	}

	err := ts.db.Create(otherEntityType).Error
	ts.Nil(err)

	otherEntity1, err := ts.eavService.CreateEntity(otherEntityType, map[string]any{})
	ts.Nil(err)

	otherEntity2, err := ts.eavService.CreateEntity(otherEntityType, map[string]any{})
	ts.Nil(err)

	relationAttr := models.NewRelationAttribute(
		ts.profileType, "relation",
		false, false, otherEntityType,
	)

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		relationAttr,
	)

	err = ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	match, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"relation": otherEntity1.ID.String(),
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"relation": otherEntity2.ID.String(),
	})
	ts.Nil(err)

	params := map[string]string{
		"relation": otherEntity1.ID.String(),
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamFilterByNull() {
	match, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": nil,
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "something",
	})
	ts.Nil(err)

	params := map[string]string{
		"displayName": "null",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	EqualEntityList(&ts.Suite, []*models.Entity{match}, entities)
}

// ------------------------- CreateEntity --------------------------------

func (ts *EAVServiceIntTestSuite) TestCreateMultipleEntitiesDoesNotGenerateGormBug() {
	initialDisplayNameID := ts.displayNameAttr.ID
	initialDescriptionID := ts.descriptionAttr.ID

	for i := 0; i < 10; i++ {
		params := map[string]any{
			"displayName": fmt.Sprintf("displayName%d", i),
			"description": fmt.Sprintf("description%d", i),
		}
		entity, err := ts.eavService.CreateEntity(ts.profileType, params)
		ts.Nil(err)

		for _, value := range entity.Fields {
			if value.Attribute.Name == "displayName" {
				ts.Equal(initialDisplayNameID, value.AttributeID)
			} else if value.Attribute.Name == "description" {
				ts.Equal(initialDescriptionID, value.AttributeID)
			}
		}
	}
}

func (ts *EAVServiceIntTestSuite) TestCreateReturnsErrorIfUUIDCantBeParsed() {
	otherType := &models.EntityType{
		Name: "other",
	}
	err := ts.db.Create(&otherType).Error
	ts.Nil(err)

	relationAttr := models.NewRelationAttribute(ts.profileType, "relation", false, false, otherType)
	ts.profileType.Attributes = append(ts.profileType.Attributes,
		relationAttr,
	)

	err = ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	params := map[string]any{
		"displayName": "displayName",
		"relation":    "not-a-uuid",
	}
	entity, err := ts.eavService.CreateEntity(ts.profileType, params)
	ts.Nil(entity)
	ts.ErrorIs(err, eavservice.ErrCantParseUUID)
}

func (ts *EAVServiceIntTestSuite) TestCreatesDefaultAttributes() {
	intAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "int",
		ValueType:    models.IntValueType,
		Required:     false,
		Default:      true,
		DefaultInt:   1,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		intAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)
	ts.Len(entity.Fields, 3)
	notNull := pie.Filter(entity.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 1)
	ts.Equal(1, notNull[0].Value())
}

func (ts *EAVServiceIntTestSuite) TestCreatesWithoutRequiredValueRespondsError() {
	requiredAttr := &models.Attribute{
		Name:      "required",
		Required:  true,
		ValueType: models.StringValueType,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		requiredAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(entity)
	ts.ErrorContains(err, "field required is missing and is required")
}

func (ts *EAVServiceIntTestSuite) TestCreatesIntAttributeEvenIfItIsInFloatFormat() {
	intAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "int",
		ValueType:    models.IntValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		intAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"int": 2.0,
	})
	ts.Nil(err)
	ts.Len(entity.Fields, 3)
	notNull := pie.Filter(entity.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 1)
	ts.Equal(2, notNull[0].Value())
}

// ------------------------- UpdateEntity --------------------------------

func (ts *EAVServiceIntTestSuite) TestUpdateEntityMultipleTimesDoesNotGenerateGormBug() {
	initialDisplayNameID := ts.displayNameAttr.ID
	initialDescriptionID := ts.descriptionAttr.ID

	params := map[string]any{
		"displayName": "displayName",
		"description": "description",
	}
	entity, err := ts.eavService.CreateEntity(ts.profileType, params)
	ts.Nil(err)

	params2 := map[string]any{
		"displayName": "other",
		"description": "other",
	}
	_, err = ts.eavService.CreateEntity(ts.profileType, params2)
	ts.Nil(err)

	for i := 0; i < 10; i++ {
		params := map[string]any{
			"displayName": fmt.Sprintf("displayName%d", i),
			"description": fmt.Sprintf("description%d", i),
		}
		err := ts.eavService.UpdateEntity(entity, params)
		ts.Nil(err)

		for _, value := range entity.Fields {
			if value.Attribute.Name == "displayName" {
				ts.Equal(initialDisplayNameID, value.AttributeID)
			} else if value.Attribute.Name == "description" {
				ts.Equal(initialDescriptionID, value.AttributeID)
			}
		}
	}
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityWorksForAllTheTypes() {
	intAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "int",
		ValueType:    models.IntValueType,
		Required:     false,
	}

	floatAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "float",
		ValueType:    models.FloatValueType,
		Required:     false,
	}

	boolAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "bool",
		ValueType:    models.BooleanValueType,
		Required:     false,
	}

	otherEntityType := &models.EntityType{
		Name: "other",
	}

	err := ts.db.Create(otherEntityType).Error
	ts.Nil(err)

	otherEntity1, err := ts.eavService.CreateEntity(otherEntityType, map[string]any{})
	ts.Nil(err)

	relationAttr := models.NewRelationAttribute(
		ts.profileType, "relation",
		false, false, otherEntityType,
	)

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		intAttr,
		floatAttr,
		boolAttr,
		relationAttr,
	)

	err = ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	params := map[string]any{
		"displayName": "displayName",
	}
	entity, err := ts.eavService.CreateEntity(ts.profileType, params)
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"displayName": nil,
		"int":         1,
		"float":       1.0,
		"bool":        true,
		"relation":    otherEntity1.ID.String(),
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.Nil(err)
	ts.Len(entity.Fields, 6)
	notNull := pie.Filter(entity.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 4)
	values := pie.Map(notNull, func(v *models.Value) any {
		return v.Value()
	})
	ts.Contains(values, 1)
	ts.Contains(values, 1.0)
	ts.Contains(values, true)
	ts.Contains(values, otherEntity1.ID)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfStringForIntType() {
	intAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "int",
		ValueType:    models.IntValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		intAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"int": "1",
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfStringForFloatType() {
	floatAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "float",
		ValueType:    models.FloatValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		floatAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"float": "1",
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfStringForBoolType() {
	boolAttr := &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "bool",
		ValueType:    models.BooleanValueType,
		Required:     false,
	}

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		boolAttr,
	)

	err := ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"bool": "1",
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfIntForStringType() {
	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"displayName": 1,
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityReturnsErrorIfUUIDCantBeParsedForRelationType() {
	otherEntityType := &models.EntityType{
		Name: "other",
	}

	err := ts.db.Create(otherEntityType).Error
	ts.Nil(err)

	relationAttr := models.NewRelationAttribute(
		ts.profileType, "relation",
		false, false, otherEntityType,
	)

	ts.profileType.Attributes = append(ts.profileType.Attributes,
		relationAttr,
	)

	err = ts.db.Save(&ts.profileType).Error
	ts.Nil(err)

	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"relation": "not-uuid",
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.ErrorIs(err, eavservice.ErrCantParseUUID)
}

func (ts *EAVServiceIntTestSuite) TestUpdateEntityDoesNotUpdateAValueIfOtherFails() {
	entity, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{})
	ts.Nil(err)

	paramsUpdate := map[string]any{
		"displayName": "something",
		"description": 1,
	}
	err = ts.eavService.UpdateEntity(entity, paramsUpdate)
	ts.ErrorIs(err, models.ErrAskingForWrongType)

	entityReturned, err := ts.eavService.GetEntity(ts.profileType, entity.ID)
	ts.Nil(err)
	notNull := pie.Filter(entityReturned.Fields, func(value *models.Value) bool {
		return !value.IsNull
	})
	ts.Len(notNull, 0)
}

// ------------------------- GetEntityTypeByName -------------------------

func (ts *EAVServiceIntTestSuite) TestGetEntityTypeReturnsErrorIfItDoesNotExist() {
	_, err := ts.eavService.GetEntityTypeByName("not_exists")
	ts.ErrorContains(err, "EntityType named \"not_exists\" not found")
}

func (ts *EAVServiceIntTestSuite) TestGetEntityTypeWorksIfItExists() {
	ett, err := ts.eavService.GetEntityTypeByName("profile")
	ts.Nil(err)
	assert.DeepEqual(ts.T(), ts.profileType, ett)
}

// ------------------------- DeleteEntity -------------------------

func (ts *EAVServiceIntTestSuite) TestDeleteWorks() {
	profile1, err := ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "displayName",
	})
	ts.Nil(err)

	err = ts.eavService.DeleteEntity(profile1)
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

func (ts *EAVServiceIntTestSuite) createProfile(entityType *models.EntityType, displayName string) *models.Entity {
	entity := &models.Entity{
		EntityTypeID: entityType.ID,
		EntityType:   entityType,
	}

	displayNameVal, _ := models.NewStringValue(ts.displayNameAttr, displayName)
	descriptionVal, _ := models.NewStringValue(ts.descriptionAttr, "something in description")

	entity.Fields = append(entity.Fields,
		displayNameVal,
		descriptionVal,
	)

	err := ts.entityRepository.Create(entity)
	ts.Nil(err)

	return entity
}
