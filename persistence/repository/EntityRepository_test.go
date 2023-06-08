package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/persistence/models"
)

type EntityRepositorySuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	repository EntityRepository
	query      *gorm.DB
	uuid       uuid.UUID
}

func (s *EntityRepositorySuite) SetupSuite() {
	s.repository = EntityRepository{}
}

func (s *EntityRepositorySuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.query = s.DB.Select("entities.*")
	s.uuid = uuid.New()
}

func Test(t *testing.T) {
	suite.Run(t, new(EntityRepositorySuite))
}

type AttributeNameAndValue struct {
	AttributeName  string
	AttributeValue any
}

func (s *EntityRepositorySuite) TestAddValueCheckAddsValueCheckForString() {
	attributeName := "attrName"
	attributeValue := "a string"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName ON
			attributes_attrName.entity_type_id = entities.entity_type_id AND
			attributes_attrName.name = $1
		JOIN values_ values_attrName ON
			values_attrName.attribute_id = attributes_attrName.id AND
			values_attrName.entity_id = entities.id AND
			attributes_attrName.value_type = 'string' AND
			values_attrName.string_val = $2
		WHERE entities.entity_type_id = $3 AND
			"entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName, attributeValue, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(AttributeNameAndValue{attributeName, attributeValue})
}

func (s *EntityRepositorySuite) TestAddValueCheckAddsValueCheckForStringUUID() {
	attributeName := "attrName"
	attributeValue := uuid.New().String()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName ON
			attributes_attrName.entity_type_id = entities.entity_type_id AND
			attributes_attrName.name = $1
		JOIN values_ values_attrName ON
			values_attrName.attribute_id = attributes_attrName.id AND
			values_attrName.entity_id = entities.id AND
			attributes_attrName.value_type = 'relation' AND
			values_attrName.relation_val = $2
		WHERE entities.entity_type_id = $3 AND
			"entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName, attributeValue, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(AttributeNameAndValue{attributeName, attributeValue})
}

func (s *EntityRepositorySuite) TestAddValueCheckAddsValueCheckForBool() {
	attributeName := "attrName"
	attributeValue := true
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName ON
			attributes_attrName.entity_type_id = entities.entity_type_id AND
			attributes_attrName.name = $1
		JOIN values_ values_attrName ON
			values_attrName.attribute_id = attributes_attrName.id AND
			values_attrName.entity_id = entities.id AND
			attributes_attrName.value_type = 'bool' AND
			values_attrName.bool_val = $2
		WHERE entities.entity_type_id = $3 AND
			"entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName, attributeValue, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(AttributeNameAndValue{attributeName, attributeValue})
}

func (s *EntityRepositorySuite) TestAddValueCheckAddsValueCheckForNil() {
	attributeName := "attrName"
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName ON
			attributes_attrName.entity_type_id = entities.entity_type_id AND
			attributes_attrName.name = $1
		JOIN values_ values_attrName ON
			values_attrName.attribute_id = attributes_attrName.id AND
			values_attrName.entity_id = entities.id AND
			values_attrName.is_null = $2
		WHERE entities.entity_type_id = $3 AND
			"entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName, true, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(AttributeNameAndValue{attributeName, nil})
}

func (s *EntityRepositorySuite) TestAddValueCheckAddsValueCheckIntAndFloatForFloat() {
	attributeName := "attrName"
	attributeValue := 1.2
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName ON
			attributes_attrName.entity_type_id = entities.entity_type_id AND
			attributes_attrName.name = $1
		JOIN values_ values_attrName ON
			values_attrName.attribute_id = attributes_attrName.id AND
			values_attrName.entity_id = entities.id AND
			((attributes_attrName.value_type = 'int' AND values_attrName.int_val = $2) OR
			(attributes_attrName.value_type = 'float' AND values_attrName.float_val = $3))
		WHERE entities.entity_type_id = $4 AND "entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName, attributeValue, attributeValue, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(AttributeNameAndValue{attributeName, attributeValue})
}

func (s *EntityRepositorySuite) TestAddValueCheck2AddsValueCheckFor2Values() {
	attributeName1 := "attrName1"
	attributeValue1 := "a string"
	attributeName2 := "attrName2"
	attributeValue2 := true
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName1 ON
			attributes_attrName1.entity_type_id = entities.entity_type_id AND
			attributes_attrName1.name = $1
		JOIN values_ values_attrName1 ON
			values_attrName1.attribute_id = attributes_attrName1.id AND
			values_attrName1.entity_id = entities.id AND
			attributes_attrName1.value_type = 'string' AND
			values_attrName1.string_val = $2
		JOIN attributes attributes_attrName2 ON
			attributes_attrName2.entity_type_id = entities.entity_type_id AND
			attributes_attrName2.name = $3
		JOIN values_ values_attrName2 ON
			values_attrName2.attribute_id = attributes_attrName2.id AND
			values_attrName2.entity_id = entities.id AND
			attributes_attrName2.value_type = 'bool' AND
			values_attrName2.bool_val = $4
		WHERE entities.entity_type_id = $5 AND
			"entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName1, attributeValue1, attributeName2, attributeValue2, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(
		AttributeNameAndValue{attributeName1, attributeValue1},
		AttributeNameAndValue{attributeName2, attributeValue2},
	)
}

func (s *EntityRepositorySuite) TestAddValueCheckAddsJoinWithEntitiesForMap() {
	attributeName := "attrName"
	innerAttributeName := "innerAttrName"
	innerAttributeValue := "a string"
	attributeValue := map[string]any{innerAttributeName: innerAttributeValue}
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT entities.* FROM "entities"
		JOIN attributes attributes_attrName ON
			attributes_attrName.entity_type_id = entities.entity_type_id AND
			attributes_attrName.name = $1
		JOIN values_ values_attrName ON
			values_attrName.attribute_id = attributes_attrName.id AND
			values_attrName.entity_id = entities.id AND
			attributes_attrName.value_type = 'relation'
		JOIN entities entities_attrName ON
			entities_attrName.id = values_attrName.relation_val AND
			entities_attrName.deleted_at IS NULL
		JOIN attributes attributes_attrName_innerAttrName ON
			attributes_attrName_innerAttrName.entity_type_id = entities_attrName.entity_type_id AND
			attributes_attrName_innerAttrName.name = $2
		JOIN values_ values_attrName_innerAttrName ON
			values_attrName_innerAttrName.attribute_id = attributes_attrName_innerAttrName.id AND
			values_attrName_innerAttrName.entity_id = entities_attrName.id AND
			attributes_attrName_innerAttrName.value_type = 'string' AND
			values_attrName_innerAttrName.string_val = $3
		WHERE entities.entity_type_id = $4 AND
			"entities"."deleted_at" IS NULL`)).
		WithArgs(attributeName, innerAttributeName, innerAttributeValue, s.uuid).
		WillReturnRows(sqlmock.NewRows(nil))

	s.execQuery(AttributeNameAndValue{attributeName, attributeValue})
}

func (s *EntityRepositorySuite) execQuery(attributes ...AttributeNameAndValue) {
	for _, attribute := range attributes {
		err := s.repository.AddValueCheckToQuery(s.query, attribute.AttributeName, attribute.AttributeValue)
		require.NoError(s.T(), err)
	}

	s.query.Where(
		"entities.entity_type_id = ?",
		s.uuid,
	)

	var entities []*models.Entity
	err := s.query.Find(&entities).Error
	require.NoError(s.T(), err)
}
