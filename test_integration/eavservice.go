package integration_test

import (
	"fmt"
	"log"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/eavservice"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	log.Println(ts.profileType.ID)
	log.Println(ts.displayNameAttr.Name)
	log.Println(ts.displayNameAttr.ID)
	log.Println(ts.descriptionAttr.Name)
	log.Println(ts.descriptionAttr.ID)
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
		"displayName": "match",
		"relation":    otherEntity1.ID.String(),
	})
	ts.Nil(err)

	_, err = ts.eavService.CreateEntity(ts.profileType, map[string]any{
		"displayName": "not_match",
		"relation":    otherEntity2.ID.String(),
	})
	ts.Nil(err)

	params := map[string]string{
		"relation": otherEntity1.ID.String(),
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
			} else {
				ts.Equal(initialDescriptionID, value.AttributeID)
			}
		}
	}
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
			} else {
				ts.Equal(initialDescriptionID, value.AttributeID)
			}
		}
	}
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

	log.Println("Before create")
	log.Println(displayNameVal.Attribute.Name)
	log.Println(displayNameVal.AttributeID)
	log.Println(descriptionVal.Attribute.Name)
	log.Println(descriptionVal.AttributeID)

	err := ts.entityRepository.Save(entity)
	ts.Nil(err)

	log.Println("After create")
	for _, field := range entity.Fields {
		log.Println(field.Attribute.Name)
		log.Println(field.AttributeID)
	}

	return entity
}
