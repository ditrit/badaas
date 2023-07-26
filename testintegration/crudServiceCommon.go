package testintegration

import (
	"github.com/stretchr/testify/suite"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/testintegration/models"
)

type CRUDServiceCommonIntTestSuite struct {
	suite.Suite
	db *badorm.DB
}

func (ts *CRUDServiceCommonIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *CRUDServiceCommonIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

func (ts *CRUDServiceCommonIntTestSuite) createProduct(stringV string, intV int, floatV float64, boolV bool, intP *int) *models.Product {
	entity := &models.Product{
		String:     stringV,
		Int:        intV,
		Float:      floatV,
		Bool:       boolV,
		IntPointer: intP,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createSale(code int, product *models.Product, seller *models.Seller) *models.Sale {
	entity := &models.Sale{
		Code:    code,
		Product: *product,
		Seller:  seller,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createSeller(name string, company *models.Company) *models.Seller {
	var companyID *badorm.UUID
	if company != nil {
		companyID = &company.ID
	}

	entity := &models.Seller{
		Name:      name,
		CompanyID: companyID,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createCompany(name string) *models.Company {
	entity := &models.Company{
		Name: name,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createCountry(name string, capital models.City) *models.Country {
	entity := &models.Country{
		Name:    name,
		Capital: capital,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createEmployee(name string, boss *models.Employee) *models.Employee {
	entity := &models.Employee{
		Name: name,
		Boss: boss,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createBicycle(name string, owner models.Person) *models.Bicycle {
	entity := &models.Bicycle{
		Name:  name,
		Owner: owner,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createBrand(name string) *models.Brand {
	entity := &models.Brand{
		Name: name,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createPhone(name string, brand models.Brand) *models.Phone {
	entity := &models.Phone{
		Name:  name,
		Brand: brand,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceCommonIntTestSuite) createUniversity(name string) *models.University {
	entity := &models.University{
		Name: name,
	}
	err := ts.db.GormDB.Create(entity).Error
	ts.Nil(err)

	return entity
}
