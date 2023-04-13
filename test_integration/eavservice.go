package integration_test

import (
	"log"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/eavservice"
)

type EAVServiceIntTestSuite struct {
	IntegrationTestSuite
	eavService      eavservice.EAVService
	profileType     *models.EntityType
	displayNameAttr *models.Attribute
	descriptionAttr *models.Attribute
}

func NewEAVServiceIntTestSuite(
	ts *IntegrationTestSuite,
	eavService eavservice.EAVService,
) *EAVServiceIntTestSuite {
	return &EAVServiceIntTestSuite{
		IntegrationTestSuite: *ts,
		eavService:           eavService,
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

	// err := ts.db.Create(ts.profileType).Error
	// if err != nil {
	// 	ts.Fail("Unable to create entity type: ", err)
	// }

	log.Println(ts.displayNameAttr.ID)
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

// TODO verificar cuando el atributo nisiquiera existe
// TODO verificar con otros tipos de atributos
// TODO verificar cuando hay otros entityTypes

func (ts *EAVServiceIntTestSuite) createProfile(entityType *models.EntityType, displayName string) *models.Entity {
	entity := &models.Entity{
		EntityTypeID: entityType.ID,
		EntityType:   entityType,
	}

	displayNameVal, _ := models.NewStringValue(ts.displayNameAttr, displayName)
	// descriptionVal, _ := models.NewNullValue(ts.descriptionAttr)
	entity.Fields = append(entity.Fields,
		displayNameVal,
		// descriptionVal,
	)

	err := ts.db.Create(entity).Error
	if err != nil {
		ts.Fail("Unable to create entity: ", err)
	}

	return entity
}
