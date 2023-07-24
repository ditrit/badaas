package testintegration

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/dynamic"
	"github.com/ditrit/badaas/badorm/multitype"
	"github.com/ditrit/badaas/badorm/unsafe"
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
	crudPhoneService    badorm.CRUDService[models.Phone, badorm.UIntID]
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
	crudPhoneService badorm.CRUDService[models.Phone, badorm.UIntID],
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

	entities, err := ts.crudPhoneService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudEmployeeService.Query(
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

	entities, err := ts.crudCityService.Query(
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

	entities, err := ts.crudCountryService.Query(
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

	entities, err := ts.crudBicycleService.Query(
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

	entities, err := ts.crudSellerService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

	entities, err := ts.crudSaleService.Query(
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

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorOver2Tables() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	ts.createSeller("agustin", company2)

	entities, err := ts.crudSellerService.Query(
		conditions.SellerCompany(
			conditions.CompanyName(
				dynamic.Eq(conditions.SellerNameField),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{seller1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorOver2TablesAtMoreLevel() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleSeller(
			conditions.SellerCompany(
				conditions.CompanyName(
					dynamic.Eq(conditions.SellerNameField),
				),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorWithNotJoinedModelReturnsError() {
	_, err := ts.crudChildService.Query(
		conditions.ChildId(dynamic.Eq(conditions.ParentParentIdField)),
	)
	ts.ErrorIs(err, badorm.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "not concerned model: models.ParentParent; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithoutSelectJoinReturnsError() {
	_, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1ParentParent(),
		),
		conditions.ChildParent2(
			conditions.Parent2ParentParent(),
		),
		conditions.ChildId(dynamic.Eq(conditions.ParentParentIdField)),
	)
	ts.ErrorIs(err, badorm.ErrJoinMustBeSelected)
	ts.ErrorContains(err, "joined multiple times model: models.ParentParent; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithSelectJoin() {
	parentParent := &models.ParentParent{Name: "franco"}
	parent1 := &models.Parent1{ParentParent: *parentParent}
	parent2 := &models.Parent2{ParentParent: *parentParent}
	child := &models.Child{Parent1: *parent1, Parent2: *parent2, Name: "franco"}
	err := ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1ParentParent(),
		),
		conditions.ChildParent2(
			conditions.Parent2ParentParent(),
		),
		conditions.ChildName(
			dynamic.Eq(conditions.ParentParentNameField).SelectJoin(0, 0),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithoutSelectJoinOnMultivalueOperatorReturnsError() {
	_, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1ParentParent(),
		),
		conditions.ChildParent2(
			conditions.Parent2ParentParent(),
		),
		conditions.ChildId(
			dynamic.Between(conditions.ParentParentIdField, conditions.ParentParentIdField),
		),
	)
	ts.ErrorIs(err, badorm.ErrJoinMustBeSelected)
	ts.ErrorContains(err, "joined multiple times model: models.ParentParent; operator: Between; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithSelectJoinOnMultivalueOperator() {
	parentParent := &models.ParentParent{Name: "franco", Number: 3}
	parent1 := &models.Parent1{ParentParent: *parentParent}
	parent2 := &models.Parent2{ParentParent: *parentParent}
	child := &models.Child{Parent1: *parent1, Parent2: *parent2, Name: "franco", Number: 2}
	err := ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := ts.crudChildService.Query(
		conditions.ChildParent1(
			conditions.Parent1ParentParent(),
		),
		conditions.ChildParent2(
			conditions.Parent2ParentParent(),
		),
		conditions.ChildNumber(
			multitype.Between[int, int](1, conditions.ParentParentNumberField).SelectJoin(1, 0),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestJoinWithUnsafeCondition() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleSeller(
			conditions.SellerCompany(
				unsafe.NewCondition[models.Company]("%s.name = Seller.name"),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestJoinWithEmptyConnectionConditionMakesNothing() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match1 := ts.createSale(0, product1, nil)
	match2 := ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			badorm.And[models.Product](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match1, match2}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestJoinWithEmptyContainerConditionMakesNothing() {
	_, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			badorm.Not[models.Product](),
		),
	)
	ts.ErrorIs(err, badorm.ErrEmptyConditions)
	ts.ErrorContains(err, "connector: Not; model: models.Product")
}
