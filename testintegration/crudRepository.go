package testintegration

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type CRUDRepositoryIntTestSuite struct {
	suite.Suite
	db                    *badorm.DB
	crudProductRepository badorm.CRUDRepository[models.Product, badorm.UUID]
}

func NewCRUDRepositoryIntTestSuite(
	db *badorm.DB,
	crudProductRepository badorm.CRUDRepository[models.Product, badorm.UUID],
) *CRUDRepositoryIntTestSuite {
	return &CRUDRepositoryIntTestSuite{
		db:                    db,
		crudProductRepository: crudProductRepository,
	}
}

func (ts *CRUDRepositoryIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *CRUDRepositoryIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

// ------------------------- GetByID --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsErrorIfIDDontMatch() {
	ts.createProduct(0)
	_, err := ts.crudProductRepository.GetByID(ts.db.GormDB, badorm.NilUUID)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsEntityIfIDMatch() {
	product := ts.createProduct(0)
	ts.createProduct(0)
	productReturned, err := ts.crudProductRepository.GetByID(ts.db.GormDB, product.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

// ------------------------- QueryOne --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetReturnsErrorIfConditionsDontMatch() {
	ts.createProduct(0)
	_, err := ts.crudProductRepository.QueryOne(
		ts.db.GormDB,
		conditions.ProductInt(badorm.Eq(1)),
	)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetReturnsEntityIfConditionsMatch() {
	product := ts.createProduct(1)
	productReturned, err := ts.crudProductRepository.QueryOne(
		ts.db.GormDB,
		conditions.ProductInt(badorm.Eq(1)),
	)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

// ------------------------- utils -------------------------

func (ts *CRUDRepositoryIntTestSuite) createProduct(intV int) *models.Product {
	entity := &models.Product{
		Int: intV,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}
