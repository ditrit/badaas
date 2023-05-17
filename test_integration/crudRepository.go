package integrationtests

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gotest.tools/assert"
)

type CRUDRepositoryIntTestSuite struct {
	suite.Suite
	logger                *zap.Logger
	db                    *gorm.DB
	crudProductRepository badorm.CRUDRepository[Product, uuid.UUID]
}

func NewCRUDRepositoryIntTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
	crudProductRepository badorm.CRUDRepository[Product, uuid.UUID],
) *CRUDRepositoryIntTestSuite {
	return &CRUDRepositoryIntTestSuite{
		logger:                logger,
		db:                    db,
		crudProductRepository: crudProductRepository,
	}
}

func (ts *CRUDRepositoryIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

// ------------------------- GetByID --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsErrorIfIDDontMatch() {
	ts.createProduct(0)
	_, err := ts.crudProductRepository.GetByID(ts.db, uuid.Nil)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsEntityIfIDMatch() {
	product := ts.createProduct(0)
	ts.createProduct(0)
	productReturned, err := ts.crudProductRepository.GetByID(ts.db, product.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

// ------------------------- Get --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetReturnsErrorIfConditionsDontMatch() {
	ts.createProduct(0)
	_, err := ts.crudProductRepository.Get(ts.db, ProductIntCondition(1))
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetReturnsEntityIfConditionsMatch() {
	product := ts.createProduct(1)
	productReturned, err := ts.crudProductRepository.Get(ts.db, ProductIntCondition(1))
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

// ------------------------- GetOptional --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetOptionalReturnsNilIfConditionsDontMatch() {
	ts.createProduct(0)
	productReturned, err := ts.crudProductRepository.GetOptional(ts.db, ProductIntCondition(1))
	ts.Nil(err)
	ts.Nil(productReturned)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetOptionalReturnsEntityIfConditionsMatch() {
	product := ts.createProduct(0)
	productReturned, err := ts.crudProductRepository.GetOptional(ts.db, ProductIntCondition(0))
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetOptionalReturnsErrorIfMoreThanOneMatchConditions() {
	ts.createProduct(0)
	ts.createProduct(0)
	_, err := ts.crudProductRepository.GetOptional(ts.db, ProductIntCondition(0))
	ts.Error(err, badorm.ErrMoreThanOneObjectFound)
}

// ------------------------- GetAll --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetAllReturnsEmptyIfNotEntitiesCreated() {
	productsReturned, err := ts.crudProductRepository.GetAll(ts.db)
	ts.Nil(err)
	EqualList(&ts.Suite, []*Product{}, productsReturned)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetAllReturnsAllEntityIfConditionsMatch() {
	product1 := ts.createProduct(0)
	product2 := ts.createProduct(0)
	productsReturned, err := ts.crudProductRepository.GetAll(ts.db)
	ts.Nil(err)
	EqualList(&ts.Suite, []*Product{product1, product2}, productsReturned)
}

// ------------------------- utils -------------------------

func (ts *CRUDRepositoryIntTestSuite) createProduct(intV int) *Product {
	entity := &Product{
		Int: intV,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}
