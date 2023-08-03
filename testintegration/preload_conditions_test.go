package testintegration

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type PreloadConditionsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudSaleService     orm.CRUDService[models.Sale, orm.UUID]
	crudSellerService   orm.CRUDService[models.Seller, orm.UUID]
	crudCountryService  orm.CRUDService[models.Country, orm.UUID]
	crudCityService     orm.CRUDService[models.City, orm.UUID]
	crudEmployeeService orm.CRUDService[models.Employee, orm.UUID]
	crudPhoneService    orm.CRUDService[models.Phone, uint]
	crudChildService    orm.CRUDService[models.Child, orm.UUID]
}

func NewPreloadConditionsIntTestSuite(
	db *gorm.DB,
	crudSaleService orm.CRUDService[models.Sale, orm.UUID],
	crudSellerService orm.CRUDService[models.Seller, orm.UUID],
	crudCountryService orm.CRUDService[models.Country, orm.UUID],
	crudCityService orm.CRUDService[models.City, orm.UUID],
	crudEmployeeService orm.CRUDService[models.Employee, orm.UUID],
	crudPhoneService orm.CRUDService[models.Phone, uint],
	crudChildService orm.CRUDService[models.Child, orm.UUID],
) *PreloadConditionsIntTestSuite {
	return &PreloadConditionsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudSaleService:     crudSaleService,
		crudSellerService:   crudSellerService,
		crudCountryService:  crudCountryService,
		crudCityService:     crudCityService,
		crudEmployeeService: crudEmployeeService,
		crudPhoneService:    crudPhoneService,
		crudChildService:    crudChildService,
	}
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadWithoutWhereConditionDoesNotFilter() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)

	withSeller := ts.createSale(0, product1, seller1)
	withoutSeller := ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.Query(
		conditions.SalePreloadSeller,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{withSeller, withoutSeller}, entities)
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		return sale.Seller.Equal(*seller1)
	}))
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		return sale.Seller == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadUIntModel() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	phone1 := ts.createPhone("pixel", *brand1)
	phone2 := ts.createPhone("iphone", *brand2)

	entities, err := ts.crudPhoneService.Query(
		conditions.PhonePreloadBrand,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Phone{phone1, phone2}, entities)
	ts.True(pie.Any(entities, func(phone *models.Phone) bool {
		return phone.Brand.Equal(*brand1)
	}))
	ts.True(pie.Any(entities, func(phone *models.Phone) bool {
		return phone.Brand.Equal(*brand2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadWithWhereConditionFilters() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product1.EmbeddedInt = 1
	product1.GormEmbedded.Int = 2
	err := ts.db.Save(product1).Error
	ts.Nil(err)

	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductPreloadAttributes,
			conditions.ProductInt(orm.Eq(1)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
	assert.DeepEqual(ts.T(), *product1, entities[0].Product)
	ts.Equal("a_string", entities[0].Product.String)
	ts.Equal(1, entities[0].Product.EmbeddedInt)
	ts.Equal(2, entities[0].Product.GormEmbedded.Int)
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadOneToOne() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	country2 := ts.createCountry("France", capital2)

	entities, err := ts.crudCityService.Query(
		conditions.CityPreloadCountry,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.City{&capital1, &capital2}, entities)
	ts.True(pie.Any(entities, func(city *models.City) bool {
		return city.Country.Equal(*country1)
	}))
	ts.True(pie.Any(entities, func(city *models.City) bool {
		return city.Country.Equal(*country2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	entities, err := ts.crudSellerService.Query(
		conditions.SellerPreloadCompany,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{seller1, seller2}, entities)
	ts.True(pie.Any(entities, func(seller *models.Seller) bool {
		return seller.Company.Equal(*company1)
	}))
	ts.True(pie.Any(entities, func(seller *models.Seller) bool {
		return seller.Company.Equal(*company2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadOneToOneReversed() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	country2 := ts.createCountry("France", capital2)

	entities, err := ts.crudCountryService.Query(
		conditions.CountryPreloadCapital,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Country{country1, country2}, entities)
	ts.True(pie.Any(entities, func(country *models.Country) bool {
		return country.Capital.Equal(capital1)
	}))
	ts.True(pie.Any(entities, func(country *models.Country) bool {
		return country.Capital.Equal(capital2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadSelfReferential() {
	boss1 := &models.Employee{
		Name: "Xavier",
	}

	employee1 := ts.createEmployee("franco", boss1)
	employee2 := ts.createEmployee("pierre", nil)

	entities, err := ts.crudEmployeeService.Query(
		conditions.EmployeePreloadBoss,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Employee{boss1, employee1, employee2}, entities)

	ts.True(pie.Any(entities, func(employee *models.Employee) bool {
		return employee.Boss != nil && employee.Boss.Equal(*boss1)
	}))
	ts.True(pie.Any(entities, func(employee *models.Employee) bool {
		return employee.Boss == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadSelfReferentialAtSecondLevel() {
	bossBoss := &models.Employee{
		Name: "Xavier",
	}
	boss := &models.Employee{
		Name: "Vincent",
		Boss: bossBoss,
	}
	employee := ts.createEmployee("franco", boss)

	entities, err := ts.crudEmployeeService.Query(
		conditions.EmployeeBoss(
			conditions.EmployeeBoss(
				conditions.EmployeePreloadAttributes,
			),
		),
		conditions.EmployeeName(orm.Eq("franco")),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Employee{employee}, entities)

	bossLoaded := entities[0].Boss
	ts.True(bossLoaded.Equal(*boss))

	bossBossLoaded := bossLoaded.Boss
	ts.True(bossBossLoaded.Equal(*bossBoss))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadDifferentEntitiesWithConditions() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductPreloadAttributes,
			conditions.ProductInt(orm.Eq(1)),
		),
		conditions.SaleSeller(
			conditions.SellerPreloadAttributes,
			conditions.SellerName(orm.Eq("franco")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
	assert.DeepEqual(ts.T(), *product1, entities[0].Product)
	assert.DeepEqual(ts.T(), seller1, entities[0].Seller)
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadDifferentEntitiesWithoutConditions() {
	parentParent := &models.ParentParent{}
	err := ts.db.Create(parentParent).Error
	ts.Nil(err)

	parent1 := &models.Parent1{ParentParent: *parentParent}
	err = ts.db.Create(parent1).Error
	ts.Nil(err)

	parent2 := &models.Parent2{ParentParent: *parentParent}
	err = ts.db.Create(parent2).Error
	ts.Nil(err)

	child := &models.Child{Parent1: *parent1, Parent2: *parent2}
	err = ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildPreloadParent1,
		conditions.ChildPreloadParent2,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parent2, entities[0].Parent2)
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadRelations() {
	parentParent := &models.ParentParent{}
	err := ts.db.Create(parentParent).Error
	ts.Nil(err)

	parent1 := &models.Parent1{ParentParent: *parentParent}
	err = ts.db.Create(parent1).Error
	ts.Nil(err)

	parent2 := &models.Parent2{ParentParent: *parentParent}
	err = ts.db.Create(parent2).Error
	ts.Nil(err)

	child := &models.Child{Parent1: *parent1, Parent2: *parent2}
	err = ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildPreloadRelations...,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parent2, entities[0].Parent2)
}

func (ts *PreloadConditionsIntTestSuite) TestJoinMultipleTimesAndPreloadWithoutCondition() {
	parentParent := &models.ParentParent{}
	err := ts.db.Create(parentParent).Error
	ts.Nil(err)

	parent1 := &models.Parent1{ParentParent: *parentParent}
	err = ts.db.Create(parent1).Error
	ts.Nil(err)

	parent2 := &models.Parent2{ParentParent: *parentParent}
	err = ts.db.Create(parent2).Error
	ts.Nil(err)

	child := &models.Child{Parent1: *parent1, Parent2: *parent2}
	err = ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1PreloadAttributes,
			conditions.Parent1PreloadParentParent,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parentParent, entities[0].Parent1.ParentParent)
}

func (ts *PreloadConditionsIntTestSuite) TestJoinMultipleTimesAndPreloadWithCondition() {
	parentParent1 := &models.ParentParent{
		Name: "parentParent1",
	}
	err := ts.db.Create(parentParent1).Error
	ts.Nil(err)

	parent11 := &models.Parent1{ParentParent: *parentParent1}
	err = ts.db.Create(parent11).Error
	ts.Nil(err)

	parent21 := &models.Parent2{ParentParent: *parentParent1}
	err = ts.db.Create(parent21).Error
	ts.Nil(err)

	child1 := &models.Child{Parent1: *parent11, Parent2: *parent21}
	err = ts.db.Create(child1).Error
	ts.Nil(err)

	parentParent2 := &models.ParentParent{}
	err = ts.db.Create(parentParent2).Error
	ts.Nil(err)

	parent12 := &models.Parent1{ParentParent: *parentParent2}
	err = ts.db.Create(parent12).Error
	ts.Nil(err)

	parent22 := &models.Parent2{ParentParent: *parentParent2}
	err = ts.db.Create(parent22).Error
	ts.Nil(err)

	child2 := &models.Child{Parent1: *parent12, Parent2: *parent22}
	err = ts.db.Create(child2).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1PreloadAttributes,
			conditions.Parent1ParentParent(
				conditions.ParentParentPreloadAttributes,
				conditions.ParentParentName(orm.Eq("parentParent1")),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child1}, entities)
	assert.DeepEqual(ts.T(), *parent11, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parentParent1, entities[0].Parent1.ParentParent)
}

func (ts *PreloadConditionsIntTestSuite) TestJoinMultipleTimesAndPreloadDiamond() {
	parentParent := &models.ParentParent{}
	err := ts.db.Create(parentParent).Error
	ts.Nil(err)

	parent1 := &models.Parent1{ParentParent: *parentParent}
	err = ts.db.Create(parent1).Error
	ts.Nil(err)

	parent2 := &models.Parent2{ParentParent: *parentParent}
	err = ts.db.Create(parent2).Error
	ts.Nil(err)

	child := &models.Child{Parent1: *parent1, Parent2: *parent2}
	err = ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1PreloadAttributes,
			conditions.Parent1PreloadParentParent,
		),
		conditions.ChildParent2(
			conditions.Parent2PreloadAttributes,
			conditions.Parent2PreloadParentParent,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parent2, entities[0].Parent2)
	assert.DeepEqual(ts.T(), *parentParent, entities[0].Parent1.ParentParent)
	assert.DeepEqual(ts.T(), *parentParent, entities[0].Parent2.ParentParent)
}