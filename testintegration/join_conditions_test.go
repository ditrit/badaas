package testintegration

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type JoinConditionsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudSaleService     badorm.CRUDService[models.Sale, badorm.UUID]
	crudSellerService   badorm.CRUDService[models.Seller, badorm.UUID]
	crudCountryService  badorm.CRUDService[models.Country, badorm.UUID]
	crudCityService     badorm.CRUDService[models.City, badorm.UUID]
	crudEmployeeService badorm.CRUDService[models.Employee, badorm.UUID]
	crudBicycleService  badorm.CRUDService[models.Bicycle, badorm.UUID]
	crudPhoneService    badorm.CRUDService[models.Phone, uint]
	crudChildService    badorm.CRUDService[models.Child, badorm.UUID]
}

func NewJoinConditionsIntTestSuite(
	db *gorm.DB,
	crudSaleService badorm.CRUDService[models.Sale, badorm.UUID],
	crudSellerService badorm.CRUDService[models.Seller, badorm.UUID],
	crudCountryService badorm.CRUDService[models.Country, badorm.UUID],
	crudCityService badorm.CRUDService[models.City, badorm.UUID],
	crudEmployeeService badorm.CRUDService[models.Employee, badorm.UUID],
	crudBicycleService badorm.CRUDService[models.Bicycle, badorm.UUID],
	crudPhoneService badorm.CRUDService[models.Phone, uint],
	crudChildService badorm.CRUDService[models.Child, badorm.UUID],
) *JoinConditionsIntTestSuite {
	return &JoinConditionsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudSaleService:     crudSaleService,
		crudSellerService:   crudSellerService,
		crudCountryService:  crudCountryService,
		crudCityService:     crudCityService,
		crudEmployeeService: crudEmployeeService,
		crudBicycleService:  crudBicycleService,
		crudPhoneService:    crudPhoneService,
		crudChildService:    crudChildService,
	}
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsUintBelongsTo() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	match := ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	entities, err := ts.crudPhoneService.GetEntities(
		conditions.PhoneBrand(
			conditions.BrandName(badorm.Eq("google")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Phone{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductInt(badorm.Eq(1)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAndFiltersTheMainEntity() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(1, product1, seller1)
	ts.createSale(2, product2, seller2)
	ts.createSale(2, product1, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleCode(badorm.Eq(1)),
		conditions.SaleProduct(
			conditions.ProductInt(badorm.Eq(1)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsHasOneOptional() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSeller(
			conditions.SellerName(badorm.Eq("franco")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsHasOneSelfReferential() {
	boss1 := &models.Employee{
		Name: "Xavier",
	}
	boss2 := &models.Employee{
		Name: "Vincent",
	}

	match := ts.createEmployee("franco", boss1)
	ts.createEmployee("pierre", boss2)

	entities, err := ts.crudEmployeeService.GetEntities(
		conditions.EmployeeBoss(
			conditions.EmployeeName(badorm.Eq("Xavier")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Employee{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOneToOne() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := ts.crudCityService.GetEntities(
		conditions.CityCountry(
			conditions.CountryName(badorm.Eq("Argentina")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.City{&capital1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOneToOneReversed() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := ts.crudCountryService.GetEntities(
		conditions.CountryCapital(
			conditions.CityName(badorm.Eq("Buenos Aires")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Country{country1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsWithEntityThatDefinesTableName() {
	person1 := models.Person{
		Name: "franco",
	}
	person2 := models.Person{
		Name: "xavier",
	}

	match := ts.createBicycle("BMX", person1)
	ts.createBicycle("Shimano", person2)

	entities, err := ts.crudBicycleService.GetEntities(
		conditions.BicycleOwner(
			conditions.PersonName(badorm.Eq("franco")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Bicycle{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	entities, err := ts.crudSellerService.GetEntities(
		conditions.SellerCompany(
			conditions.CompanyName(badorm.Eq("ditrit")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductInt(badorm.Eq(1)),
			conditions.ProductString(badorm.Eq("match")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAddsDeletedAtAutomatically() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.Nil(ts.db.Delete(product2).Error)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductString(badorm.Eq("match")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnDeletedAt() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.Nil(ts.db.Delete(product1).Error)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductDeletedAt(badorm.Eq(product1.DeletedAt)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAndFiltersByNil() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	intProduct2 := 2
	product2 := ts.createProduct("", 2, 0.0, false, &intProduct2)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductIntPointer(badorm.IsNull[int]()),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsDifferentEntities() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductInt(badorm.Eq(1)),
		),
		conditions.SaleSeller(
			conditions.SellerName(badorm.Eq("franco")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsMultipleTimes() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSeller(
			conditions.SellerName(badorm.Eq("franco")),
			conditions.SellerCompany(
				conditions.CompanyName(badorm.Eq("ditrit")),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestJoinWithUnsafeCondition() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSeller(
			conditions.SellerCompany(
				badorm.NewUnsafeCondition[models.Company]("%s.name = sales__Seller.name", []any{}),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestJoinWithEmptyConnectionConditionMakesNothing() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match1 := ts.createSale(0, product1, nil)
	match2 := ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			badorm.And[models.Product](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestJoinWithEmptyContainerConditionMakesNothing() {
	_, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			badorm.Not[models.Product](),
		),
	)
	ts.ErrorIs(err, badorm.ErrEmptyConditions)
}

func (ts *JoinConditionsIntTestSuite) TestJoinAndPreloadWithoutWhereConditionDoesNotFilter() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)

	withSeller := ts.createSale(0, product1, seller1)
	withoutSeller := ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
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

func (ts *JoinConditionsIntTestSuite) TestJoinAndPreloadWithWhereConditionFilters() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product1.EmbeddedInt = 1
	product1.GormEmbedded.Int = 2
	err := ts.db.Save(product1).Error
	ts.Nil(err)

	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductPreloadAttributes,
			conditions.ProductInt(badorm.Eq(1)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
	assert.DeepEqual(ts.T(), *product1, entities[0].Product)
	ts.Equal("a_string", entities[0].Product.String)
	ts.Equal(1, entities[0].Product.EmbeddedInt)
	ts.Equal(2, entities[0].Product.GormEmbedded.Int)
}

func (ts *JoinConditionsIntTestSuite) TestJoinAndPreloadDifferentEntities() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductPreloadAttributes,
			conditions.ProductInt(badorm.Eq(1)),
		),
		conditions.SaleSeller(
			conditions.SellerPreloadAttributes,
			conditions.SellerName(badorm.Eq("franco")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
	assert.DeepEqual(ts.T(), *product1, entities[0].Product)
	assert.DeepEqual(ts.T(), seller1, entities[0].Seller)
}

func (ts *JoinConditionsIntTestSuite) TestPreloadDifferentEntities() {
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildPreloadParent1,
		conditions.ChildPreloadParent2,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parent2, entities[0].Parent2)
}

func (ts *JoinConditionsIntTestSuite) TestPreloadAll() {
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildPreloadRelations...,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parent2, entities[0].Parent2)
}

func (ts *JoinConditionsIntTestSuite) TestJoinMultipleTimesAndPreloadWithoutCondition() {
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildParent1(
			// TODO ver esto, no se si me gusta que esten separados
			conditions.Parent1PreloadAttributes,
			conditions.Parent1PreloadParentParent,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	assert.DeepEqual(ts.T(), *parent1, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parentParent, entities[0].Parent1.ParentParent)
}

func (ts *JoinConditionsIntTestSuite) TestJoinMultipleTimesAndPreloadWithCondition() {
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildParent1(
			conditions.Parent1PreloadAttributes,
			conditions.Parent1ParentParent(
				conditions.ParentParentPreloadAttributes,
				conditions.ParentParentName(badorm.Eq("parentParent1")),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child1}, entities)
	assert.DeepEqual(ts.T(), *parent11, entities[0].Parent1)
	assert.DeepEqual(ts.T(), *parentParent1, entities[0].Parent1.ParentParent)
}

// TODO que pasa si usan preaload y ademas hacen el join
// idem ahora que pasa si hago el join dos veces con la misma entidad, deberian agruparse los join?

func (ts *JoinConditionsIntTestSuite) TestJoinMultipleTimesAndPreloadDiamond() {
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

	entities, err := ts.crudChildService.GetEntities(
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
