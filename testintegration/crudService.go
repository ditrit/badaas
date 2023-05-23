package testintegration

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gotest.tools/assert"
)

type CRUDServiceIntTestSuite struct {
	suite.Suite
	logger              *zap.Logger
	db                  *gorm.DB
	crudProductService  badorm.CRUDService[Product, uuid.UUID]
	crudSaleService     badorm.CRUDService[Sale, uuid.UUID]
	crudSellerService   badorm.CRUDService[Seller, uuid.UUID]
	crudCountryService  badorm.CRUDService[Country, uuid.UUID]
	crudCityService     badorm.CRUDService[City, uuid.UUID]
	crudEmployeeService badorm.CRUDService[Employee, uuid.UUID]
	crudBicycleService  badorm.CRUDService[Bicycle, uuid.UUID]
}

func NewCRUDServiceIntTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
	crudProductService badorm.CRUDService[Product, uuid.UUID],
	crudSaleService badorm.CRUDService[Sale, uuid.UUID],
	crudSellerService badorm.CRUDService[Seller, uuid.UUID],
	crudCountryService badorm.CRUDService[Country, uuid.UUID],
	crudCityService badorm.CRUDService[City, uuid.UUID],
	crudEmployeeService badorm.CRUDService[Employee, uuid.UUID],
	crudBicycleService badorm.CRUDService[Bicycle, uuid.UUID],
) *CRUDServiceIntTestSuite {
	return &CRUDServiceIntTestSuite{
		logger:              logger,
		db:                  db,
		crudProductService:  crudProductService,
		crudSaleService:     crudSaleService,
		crudSellerService:   crudSellerService,
		crudCountryService:  crudCountryService,
		crudCityService:     crudCityService,
		crudEmployeeService: crudEmployeeService,
		crudBicycleService:  crudBicycleService,
	}
}

func (ts *CRUDServiceIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

// ------------------------- GetEntity --------------------------------

func (ts *CRUDServiceIntTestSuite) TestGetEntityReturnsErrorIfNotEntityCreated() {
	_, err := ts.crudProductService.GetEntity(uuid.Nil)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntityReturnsErrorIfNotEntityMatch() {
	ts.createProduct("", 0, 0, false)

	_, err := ts.crudProductService.GetEntity(uuid.New())
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntityReturnsTheEntityIfItIsCreate() {
	match := ts.createProduct("", 0, 0, false)

	entity, err := ts.crudProductService.GetEntity(match.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), match, entity)
}

// ------------------------- GetEntities --------------------------------

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct("", 0, 0, false)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct("", 0, 0, false)
	match2 := ts.createProduct("", 0, 0, false)
	match3 := ts.createProduct("", 0, 0, false)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match1, match2, match3}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities(
		ProductStringCondition("not_created"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct("something_else", 0, 0, false)

	entities, err := ts.crudProductService.GetEntities(
		ProductStringCondition("not_match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createProduct("match", 0, 0, false)
	ts.createProduct("not_match", 0, 0, false)

	entities, err := ts.crudProductService.GetEntities(
		ProductStringCondition("match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProduct("match", 0, 0, false)
	match2 := ts.createProduct("match", 0, 0, false)
	ts.createProduct("not_match", 0, 0, false)

	entities, err := ts.crudProductService.GetEntities(
		ProductStringCondition("match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntType() {
	match := ts.createProduct("match", 1, 0, false)
	ts.createProduct("not_match", 2, 0, false)

	entities, err := ts.crudProductService.GetEntities(
		ProductIntCondition(1),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfFloatType() {
	match := ts.createProduct("match", 0, 1.1, false)
	ts.createProduct("not_match", 0, 2.2, false)

	entities, err := ts.crudProductService.GetEntities(
		ProductFloatCondition(1.1),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfBoolType() {
	match := ts.createProduct("match", 0, 0.0, true)
	ts.createProduct("not_match", 0, 0.0, false)

	entities, err := ts.crudProductService.GetEntities(
		ProductBoolCondition(true),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createProduct("match", 1, 0.0, true)
	match2 := ts.createProduct("match", 1, 0.0, true)

	ts.createProduct("not_match", 1, 0.0, true)
	ts.createProduct("match", 2, 0.0, true)

	entities, err := ts.crudProductService.GetEntities(
		ProductStringCondition("match"),
		ProductIntCondition(1),
		ProductBoolCondition(true),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfRelationType() {
	product1 := ts.createProduct("", 0, 0.0, false)
	product2 := ts.createProduct("", 0, 0.0, false)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		SaleProductIDCondition(product1.ID),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct("", 1, 0.0, false)
	product2 := ts.createProduct("", 2, 0.0, false)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		SaleProductCondition(
			ProductIntCondition(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsAndFiltersTheMainEntity() {
	product1 := ts.createProduct("", 1, 0.0, false)
	product2 := ts.createProduct("", 2, 0.0, false)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(1, product1, seller1)
	ts.createSale(2, product2, seller2)
	ts.createSale(2, product1, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		SaleCodeCondition(1),
		SaleProductCondition(
			ProductIntCondition(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsHasOneOptional() {
	product1 := ts.createProduct("", 1, 0.0, false)
	product2 := ts.createProduct("", 2, 0.0, false)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		SaleSellerCondition(
			SellerNameCondition("franco"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsHasOneSelfReferential() {
	boss1 := &Employee{
		Name: "Xavier",
	}
	boss2 := &Employee{
		Name: "Vincent",
	}

	match := ts.createEmployee("franco", boss1)
	ts.createEmployee("pierre", boss2)

	entities, err := ts.crudEmployeeService.GetEntities(
		EmployeeBossCondition(
			EmployeeNameCondition("Xavier"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Employee{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOneToOne() {
	capital1 := City{
		Name: "Buenos Aires",
	}
	capital2 := City{
		Name: "Paris",
	}

	ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := ts.crudCityService.GetEntities(
		CityCountryCondition(
			CountryNameCondition("Argentina"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*City{&capital1}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOneToOneReversed() {
	capital1 := City{
		Name: "Buenos Aires",
	}
	capital2 := City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := ts.crudCountryService.GetEntities(
		CountryCapitalCondition(
			CityNameCondition("Buenos Aires"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Country{country1}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsWithEntityThatDefinesTableName() {
	person1 := Person{
		Name: "franco",
	}
	person2 := Person{
		Name: "xavier",
	}

	match := ts.createBicycle("BMX", person1)
	ts.createBicycle("Shimano", person2)

	entities, err := ts.crudBicycleService.GetEntities(
		BicycleOwnerCondition(
			PersonNameCondition("franco"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Bicycle{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	entities, err := ts.crudSellerService.GetEntities(
		SellerCompanyCondition(
			CompanyNameCondition("ditrit"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Seller{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct("match", 1, 0.0, false)
	product2 := ts.createProduct("match", 2, 0.0, false)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		SaleProductCondition(
			ProductIntCondition(1),
			ProductStringCondition("match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsDifferentEntities() {
	product1 := ts.createProduct("", 1, 0.0, false)
	product2 := ts.createProduct("", 2, 0.0, false)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	entities, err := ts.crudSaleService.GetEntities(
		SaleProductCondition(
			ProductIntCondition(1),
		),
		SaleSellerCondition(
			SellerNameCondition("franco"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsMultipleTimes() {
	product1 := ts.createProduct("", 0, 0.0, false)
	product2 := ts.createProduct("", 0, 0.0, false)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		SaleSellerCondition(
			SellerNameCondition("franco"),
			SellerCompanyCondition(
				CompanyNameCondition("ditrit"),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

// ------------------------- utils -------------------------

func (ts *CRUDServiceIntTestSuite) createProduct(stringV string, intV int, floatV float64, boolV bool) *Product {
	entity := &Product{
		String: stringV,
		Int:    intV,
		Float:  floatV,
		Bool:   boolV,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createSale(code int, product *Product, seller *Seller) *Sale {
	entity := &Sale{
		Code:    code,
		Product: *product,
		Seller:  seller,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createSeller(name string, company *Company) *Seller {
	var companyID *uuid.UUID
	if company != nil {
		companyID = &company.ID
	}
	entity := &Seller{
		Name:      name,
		CompanyID: companyID,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createCompany(name string) *Company {
	entity := &Company{
		Name: name,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createCountry(name string, capital City) *Country {
	entity := &Country{
		Name:    name,
		Capital: capital,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createEmployee(name string, boss *Employee) *Employee {
	entity := &Employee{
		Name: name,
		Boss: boss,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createBicycle(name string, owner Person) *Bicycle {
	entity := &Bicycle{
		Name:  name,
		Owner: owner,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}
