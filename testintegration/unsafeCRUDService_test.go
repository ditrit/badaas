package testintegration

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/unsafe"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/testintegration/models"
)

type CRUDUnsafeServiceIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService  unsafe.CRUDService[models.Product, badorm.UUID]
	crudSaleService     unsafe.CRUDService[models.Sale, badorm.UUID]
	crudSellerService   unsafe.CRUDService[models.Seller, badorm.UUID]
	crudCountryService  unsafe.CRUDService[models.Country, badorm.UUID]
	crudCityService     unsafe.CRUDService[models.City, badorm.UUID]
	crudEmployeeService unsafe.CRUDService[models.Employee, badorm.UUID]
	crudBicycleService  unsafe.CRUDService[models.Bicycle, badorm.UUID]
	crudBrandService    unsafe.CRUDService[models.Brand, badorm.UIntID]
	crudPhoneService    unsafe.CRUDService[models.Phone, badorm.UIntID]
}

func NewCRUDUnsafeServiceIntTestSuite(
	db *gorm.DB,
	crudProductService unsafe.CRUDService[models.Product, badorm.UUID],
	crudSaleService unsafe.CRUDService[models.Sale, badorm.UUID],
	crudSellerService unsafe.CRUDService[models.Seller, badorm.UUID],
	crudCountryService unsafe.CRUDService[models.Country, badorm.UUID],
	crudCityService unsafe.CRUDService[models.City, badorm.UUID],
	crudEmployeeService unsafe.CRUDService[models.Employee, badorm.UUID],
	crudBicycleService unsafe.CRUDService[models.Bicycle, badorm.UUID],
	crudBrandService unsafe.CRUDService[models.Brand, badorm.UIntID],
	crudPhoneService unsafe.CRUDService[models.Phone, badorm.UIntID],
) *CRUDUnsafeServiceIntTestSuite {
	return &CRUDUnsafeServiceIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService:  crudProductService,
		crudSaleService:     crudSaleService,
		crudSellerService:   crudSellerService,
		crudCountryService:  crudCountryService,
		crudCityService:     crudCityService,
		crudEmployeeService: crudEmployeeService,
		crudBicycleService:  crudBicycleService,
		crudBrandService:    crudBrandService,
		crudPhoneService:    crudPhoneService,
	}
}

// ------------------------- GetEntities --------------------------------

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithoutConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities(map[string]any{})
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithoutConditionsReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(map[string]any{})
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithoutConditionsReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match2 := ts.createProduct("", 0, 0, false, nil)
	match3 := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(map[string]any{})
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionsReturnsEmptyIfNotEntitiesCreated() {
	params := map[string]any{
		"string_something_else": "not_created",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct("something_else", 0, 0, false, nil)

	params := map[string]any{
		"string_something_else": "not_match",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	params := map[string]any{
		"string_something_else": "match",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	params := map[string]any{
		"string_something_else": "match",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatDoesNotExistReturnsDBError() {
	ts.createProduct("match", 0, 0, false, nil)

	params := map[string]any{
		"not_exists": "not_exists",
	}
	_, err := ts.crudProductService.GetEntities(params)
	ts.NotNil(err)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionOfIntType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	params := map[string]any{
		"int": 1,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionOfIncorrectType() {
	ts.createProduct("not_match", 1, 0, false, nil)

	params := map[string]any{
		"int": "not_an_int",
	}
	result, err := ts.crudProductService.GetEntities(params)

	switch getDBDialector() {
	case configuration.MySQL, configuration.SQLite:
		// mysql and sqlite simply matches nothing
		ts.Nil(err)
		ts.Len(result, 0)
	case configuration.PostgreSQL, configuration.SQLServer:
		// postgresql and sqlserver do the verification
		ts.NotNil(err)
	}
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionOfFloatType() {
	match := ts.createProduct("match", 0, 1.1, false, nil)
	ts.createProduct("not_match", 0, 2.2, false, nil)

	params := map[string]any{
		"float": 1.1,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionOfBoolType() {
	match := ts.createProduct("match", 0, 0.0, true, nil)
	ts.createProduct("not_match", 0, 0.0, false, nil)

	params := map[string]any{
		"bool": true,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionOfRelationType() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	params := map[string]any{
		"product_id": product1.ID.String(),
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 1, 0.0, true, nil)
	ts.createProduct("match", 2, 0.0, true, nil)

	params := map[string]any{
		"string_something_else": "match",
		"int":                   1,
		"bool":                  true,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionsOnUIntModel() {
	match := ts.createBrand("match")
	ts.createBrand("not_match")

	params := map[string]any{
		"name": "match",
	}
	entities, err := ts.crudBrandService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Brand{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsUintBelongsTo() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	match := ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	params := map[string]any{
		"Brand": map[string]any{
			"name": "google",
		},
	}
	entities, err := ts.crudPhoneService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Phone{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	params := map[string]any{
		"Product": map[string]any{
			"int": 1,
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsHasOneOptional() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	params := map[string]any{
		"Seller": map[string]any{
			"name": "franco",
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsHasOneSelfReferential() {
	boss1 := &models.Employee{
		Name: "Xavier",
	}
	boss2 := &models.Employee{
		Name: "Vincent",
	}

	match := ts.createEmployee("franco", boss1)
	ts.createEmployee("pierre", boss2)

	params := map[string]any{
		"Boss": map[string]any{
			"name": "Xavier",
		},
	}
	entities, err := ts.crudEmployeeService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Employee{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsOneToOne() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	params := map[string]any{
		"Country": map[string]any{
			"name": "Argentina",
		},
	}
	entities, err := ts.crudCityService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.City{&capital1}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsOneToOneReversed() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	params := map[string]any{
		"Capital": map[string]any{
			"name": "Buenos Aires",
		},
	}
	entities, err := ts.crudCountryService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Country{country1}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsReturnsErrorIfNoRelation() {
	params := map[string]any{
		"NotExists": map[string]any{
			"int": 1,
		},
	}
	_, err := ts.crudSaleService.GetEntities(params)
	ts.ErrorContains(err, "Sale has not attribute named NotExists or NotExistsID")
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	params := map[string]any{
		"Company": map[string]any{
			"name": "ditrit",
		},
	}
	entities, err := ts.crudSellerService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	params := map[string]any{
		"Product": map[string]any{
			"int":                   1,
			"string_something_else": "match",
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsAndFiltersTheMainEntity() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(1, product1, seller1)
	ts.createSale(2, product2, seller2)
	ts.createSale(2, product1, seller2)

	params := map[string]any{
		"Product": map[string]any{
			"int": 1,
		},
		"code": 1,
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsDifferentEntities() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	params := map[string]any{
		"Product": map[string]any{
			"int": 1,
		},
		"Seller": map[string]any{
			"name": "franco",
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDUnsafeServiceIntTestSuite) TestGetEntitiesUnsafeWithConditionThatJoinsMultipleTimes() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	params := map[string]any{
		"Seller": map[string]any{
			"name": "franco",
			"Company": map[string]any{
				"name": "ditrit",
			},
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}