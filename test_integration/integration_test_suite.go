package integration_test

import (
	"log"
	"reflect"
	"sort"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

type IntegrationTestSuite struct {
	suite.Suite
	logger *zap.Logger
	db     *gorm.DB
}

func NewIntegrationTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		logger: logger,
		db:     db,
	}
}

var ListOfTables = []any{
	models.Session{},
	models.User{},
	models.Value{},
	models.Entity{},
	models.Attribute{},
	models.EntityType{},
}

func (ts *IntegrationTestSuite) SetupTest() {
	// clean database to ensure independency between tests
	for _, table := range ListOfTables {
		err := ts.db.Unscoped().Where("1 = 1").Delete(table).Error
		if err != nil {
			log.Fatalln("could not clean database: ", err)
		}
	}
}

func (ts *IntegrationTestSuite) equalList(expected, actual any) {
	v := reflect.ValueOf(expected)
	v2 := reflect.ValueOf(actual)

	ts.Len(actual, v.Len())

	for i := 0; i < v.Len(); i++ {
		j := 0
		for ; j < v.Len(); j++ {
			if is.DeepEqual(v2.Index(j).Interface(), v.Index(i).Interface())().Success() {
				break
			}
		}
		if j == v.Len() {
			ts.Fail("element %v not in list %v", v.Index(i).Interface(), actual)
		}
	}
}

func (ts *IntegrationTestSuite) equalEntityList(expected, actual []*models.Entity) {
	ts.Len(actual, len(expected))

	sort.SliceStable(expected, func(i, j int) bool {
		return expected[i].ID.String() < expected[j].ID.String()
	})

	sort.SliceStable(actual, func(i, j int) bool {
		return actual[i].ID.String() < actual[j].ID.String()
	})

	for i := range actual {
		ts.equalEntity(expected[i], actual[i])
	}
}

func (ts *IntegrationTestSuite) equalEntity(expected, actual *models.Entity) {
	assert.DeepEqual(ts.T(), expected, actual)
	ts.equalList(
		expected.Fields,
		actual.Fields,
	)
}
