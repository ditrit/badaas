package integration_test

import (
	"fmt"
	"log"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/eavservice"
)

type EAVServiceIntTestSuite struct {
	IntegrationTestSuite
	eavService       eavservice.EAVService
	entityRepository *repository.EntityRepository
	profileType      *models.EntityType
	displayNameAttr  *models.Attribute
	descriptionAttr  *models.Attribute
}

func NewEAVServiceIntTestSuite(
	ts *IntegrationTestSuite,
	eavService eavservice.EAVService,
	entityRepository *repository.EntityRepository,
) *EAVServiceIntTestSuite {
	return &EAVServiceIntTestSuite{
		IntegrationTestSuite: *ts,
		eavService:           eavService,
		entityRepository:     entityRepository,
	}
}

func (ts *EAVServiceIntTestSuite) SetupTest() {
	ts.IntegrationTestSuite.SetupTest()

	// TODO duplicated code
	// CREATION OF THE PROFILE TYPE AND ASSOCIATED ATTRIBUTES
	ts.profileType = &models.EntityType{
		Name: "profile",
	}
	ts.displayNameAttr = &models.Attribute{
		EntityTypeID: ts.profileType.ID,
		Name:         "displayName",
		ValueType:    models.StringValueType,
		Required:     true,
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
	if err != nil {
		ts.Fail("Unable to create entity type: ", err)
	}

	log.Println(ts.profileType.ID)
	log.Println(ts.displayNameAttr.Name)
	log.Println(ts.displayNameAttr.ID)
	log.Println(ts.descriptionAttr.Name)
	log.Println(ts.descriptionAttr.ID)
}

// ------------------------- GetEntitiesWithParams --------------------------------

func (ts *EAVServiceIntTestSuite) TestWithoutParamsReturnsEmptyIfNotEntitiesCreated() {
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, make(map[string]string))

	ts.equalEntityList([]*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithoutParamsReturnsTheOnlyOneIfOneEntityCreated() {
	profile := ts.createProfile(ts.profileType, "profile")

	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, make(map[string]string))

	ts.equalEntityList([]*models.Entity{profile}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithoutParamsReturnsTheListWhenMultipleCreated() {
	profile1 := ts.createProfile(ts.profileType, "profile1")
	profile2 := ts.createProfile(ts.profileType, "profile2")
	profile3 := ts.createProfile(ts.profileType, "profile3")

	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, make(map[string]string))

	ts.equalEntityList([]*models.Entity{profile1, profile2, profile3}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsEmptyIfNotEntitiesCreated() {
	params := map[string]string{
		"displayName": "not_created",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	ts.equalEntityList([]*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsEmptyIfNothingMatch() {
	ts.createProfile(ts.profileType, "profile")

	params := map[string]string{
		"displayName": "not_match",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	ts.equalEntityList([]*models.Entity{}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsOneIfOnlyOneMatch() {
	matchProfile := ts.createProfile(ts.profileType, "match")
	ts.createProfile(ts.profileType, "something_else")

	params := map[string]string{
		"displayName": "match",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	ts.equalEntityList([]*models.Entity{matchProfile}, entities)
}

func (ts *EAVServiceIntTestSuite) TestWithParamsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProfile(ts.profileType, "match")
	match2 := ts.createProfile(ts.profileType, "match")
	ts.createProfile(ts.profileType, "something_else")

	params := map[string]string{
		"displayName": "match",
	}
	entities := ts.eavService.GetEntitiesWithParams(ts.profileType, params)

	ts.equalEntityList([]*models.Entity{match1, match2}, entities)
}

// TODO verificar cuando el atributo nisiquiera existe
// TODO verificar con otros tipos de atributos
// TODO verificar cuando hay otros entityTypes

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
	if err != nil {
		ts.Fail("Unable to create entity: ", err)
	}

	log.Println("After create")
	for _, field := range entity.Fields {
		log.Println(field.Attribute.Name)
		log.Println(field.AttributeID)
	}

	return entity
}
