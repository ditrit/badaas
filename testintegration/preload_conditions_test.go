package testintegration

import (
	"errors"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
	"github.com/ditrit/badaas/utils"
)

type PreloadConditionsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudSaleService     badorm.CRUDService[models.Sale, badorm.UUID]
	crudCompanyService  badorm.CRUDService[models.Company, badorm.UUID]
	crudSellerService   badorm.CRUDService[models.Seller, badorm.UUID]
	crudCountryService  badorm.CRUDService[models.Country, badorm.UUID]
	crudCityService     badorm.CRUDService[models.City, badorm.UUID]
	crudEmployeeService badorm.CRUDService[models.Employee, badorm.UUID]
	crudChildService    badorm.CRUDService[models.Child, badorm.UUID]
	crudPhoneService    badorm.CRUDService[models.Phone, badorm.UIntID]
}

func NewPreloadConditionsIntTestSuite(
	db *gorm.DB,
	crudSaleService badorm.CRUDService[models.Sale, badorm.UUID],
	crudCompanyService badorm.CRUDService[models.Company, badorm.UUID],
	crudSellerService badorm.CRUDService[models.Seller, badorm.UUID],
	crudCountryService badorm.CRUDService[models.Country, badorm.UUID],
	crudCityService badorm.CRUDService[models.City, badorm.UUID],
	crudEmployeeService badorm.CRUDService[models.Employee, badorm.UUID],
	crudChildService badorm.CRUDService[models.Child, badorm.UUID],
	crudPhoneService badorm.CRUDService[models.Phone, badorm.UIntID],
) *PreloadConditionsIntTestSuite {
	return &PreloadConditionsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudSaleService:    crudSaleService,
		crudCompanyService: crudCompanyService,
		crudSellerService:  crudSellerService,
		crudCountryService: crudCountryService,
		crudCityService:    crudCityService,
		// TODO hacer algun test de self reference preload
		crudEmployeeService: crudEmployeeService,
		crudChildService:    crudChildService,
		crudPhoneService:    crudPhoneService,
	}
}

func (ts *PreloadConditionsIntTestSuite) TestNoPreloadReturnsErrorOnGetRelation() {
	product := ts.createProduct("a_string", 1, 0.0, false, nil)
	seller := ts.createSeller("franco", nil)
	sale := ts.createSale(0, product, seller)

	entities, err := ts.crudSaleService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{sale}, entities)

	saleLoaded := entities[0]

	ts.False(saleLoaded.Product.IsLoaded())
	_, err = saleLoaded.GetProduct()
	ts.ErrorIs(err, badorm.ErrRelationNotLoaded)

	ts.Nil(saleLoaded.Seller)       // is nil but we cant determine why directly (not loaded or really null)
	_, err = saleLoaded.GetSeller() // GetSeller give us that information
	ts.ErrorIs(err, badorm.ErrRelationNotLoaded)
}

func (ts *PreloadConditionsIntTestSuite) TestNoPreloadWhenItsNullKnowsItsReallyNull() {
	product := ts.createProduct("a_string", 1, 0.0, false, nil)
	sale := ts.createSale(10, product, nil)

	entities, err := ts.crudSaleService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{sale}, entities)

	saleLoaded := entities[0]

	ts.False(saleLoaded.Product.IsLoaded())
	_, err = saleLoaded.GetProduct()
	ts.ErrorIs(err, badorm.ErrRelationNotLoaded)

	ts.Nil(saleLoaded.Seller)                 // is nil but we cant determine why directly (not loaded or really null)
	saleSeller, err := saleLoaded.GetSeller() // GetSeller give us that information
	ts.Nil(err)
	ts.Nil(saleSeller)
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadWithoutWhereConditionDoesNotFilter() {
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
		saleSeller, err := sale.GetSeller()
		return err == nil && saleSeller.Equal(*seller1)
	}))
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		// in this case sale.Seller will also be nil
		// but we can now it's really null in the db because err is nil
		saleSeller, err := sale.GetSeller()
		return err == nil && saleSeller == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadNullableAtSecondLevel() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	company := ts.createCompany("ditrit")

	withCompany := ts.createSeller("with", company)
	withoutCompany := ts.createSeller("without", nil)

	sale1 := ts.createSale(0, product1, withoutCompany)
	sale2 := ts.createSale(0, product2, withCompany)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSeller(
			conditions.SellerPreloadCompany,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{sale1, sale2}, entities)
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if err != nil {
			return false
		}

		sellerCompany, err := saleSeller.GetCompany()
		return err == nil && saleSeller.Name == "with" && sellerCompany != nil && sellerCompany.Equal(*company)
	}))
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if err != nil {
			return false
		}

		sellerCompany, err := saleSeller.GetCompany()
		return err == nil && saleSeller.Name == "without" && sellerCompany == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadAtSecondLevelWorksWithManualPreload() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	company := ts.createCompany("ditrit")

	withCompany := ts.createSeller("with", company)
	withoutCompany := ts.createSeller("without", nil)

	sale1 := ts.createSale(0, product1, withoutCompany)
	sale2 := ts.createSale(0, product2, withCompany)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSeller(
			conditions.SellerPreloadAttributes,
			conditions.SellerPreloadCompany,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{sale1, sale2}, entities)
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if err != nil {
			return false
		}

		sellerCompany, err := saleSeller.GetCompany()
		return err == nil && saleSeller.Name == "with" && sellerCompany != nil && sellerCompany.Equal(*company)
	}))
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if err != nil {
			return false
		}

		sellerCompany, err := saleSeller.GetCompany()
		return err == nil && saleSeller.Name == "without" && sellerCompany == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestNoPreloadNullableAtSecondLevel() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	company := ts.createCompany("ditrit")

	withCompany := ts.createSeller("with", company)
	withoutCompany := ts.createSeller("without", nil)

	sale1 := ts.createSale(0, product1, withoutCompany)
	sale2 := ts.createSale(0, product2, withCompany)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SalePreloadSeller,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{sale1, sale2}, entities)
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if err != nil {
			return false
		}

		// the not null one is not loaded
		sellerCompany, err := saleSeller.GetCompany()
		return errors.Is(err, badorm.ErrRelationNotLoaded) && sellerCompany == nil
	}))
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if err != nil {
			return false
		}

		// we can be sure the null one is null
		sellerCompany, err := saleSeller.GetCompany()
		return err == nil && sellerCompany == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadWithoutWhereConditionDoesNotFilterAtSecondLevel() {
	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)

	withSeller := ts.createSale(0, product1, seller1)
	withoutSeller := ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSeller(
			conditions.SellerPreloadCompany,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{withSeller, withoutSeller}, entities)
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		saleSeller, err := sale.GetSeller()
		if saleSeller == nil || err != nil {
			return false
		}

		sellerCompany, err := saleSeller.GetCompany()

		return err == nil && saleSeller.Equal(*seller1) && sellerCompany == nil
	}))
	ts.True(pie.Any(entities, func(sale *models.Sale) bool {
		// in this case sale.Seller will also be nil
		// but we can now it's really null in the db because err is nil
		saleSeller, err := sale.GetSeller()
		return err == nil && saleSeller == nil
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadUIntModel() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	phone1 := ts.createPhone("pixel", *brand1)
	phone2 := ts.createPhone("iphone", *brand2)

	entities, err := ts.crudPhoneService.GetEntities(
		conditions.PhonePreloadBrand,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Phone{phone1, phone2}, entities)
	ts.True(pie.Any(entities, func(phone *models.Phone) bool {
		phoneBrand, err := phone.GetBrand()
		return err == nil && phoneBrand.Equal(*brand1)
	}))
	ts.True(pie.Any(entities, func(phone *models.Phone) bool {
		phoneBrand, err := phone.GetBrand()
		return err == nil && phoneBrand.Equal(*brand2)
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

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductPreloadAttributes,
			conditions.ProductInt(badorm.Eq(1)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
	saleProduct, err := entities[0].GetProduct()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), product1, saleProduct)
	ts.Equal("a_string", saleProduct.String)
	ts.Equal(1, saleProduct.EmbeddedInt)
	ts.Equal(2, saleProduct.GormEmbedded.Int)
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

	entities, err := ts.crudCityService.GetEntities(
		conditions.CityPreloadCountry,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.City{&capital1, &capital2}, entities)
	ts.True(pie.Any(entities, func(city *models.City) bool {
		cityCountry, err := city.GetCountry()
		if err != nil {
			return false
		}

		return cityCountry.Equal(*country1)
	}))
	ts.True(pie.Any(entities, func(city *models.City) bool {
		cityCountry, err := city.GetCountry()
		if err != nil {
			return false
		}

		return cityCountry.Equal(*country2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestNoPreloadOneToOne() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}

	ts.createCountry("Argentina", capital1)

	entities, err := ts.crudCityService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.City{&capital1}, entities)
	_, err = entities[0].GetCountry()
	ts.ErrorIs(err, badorm.ErrRelationNotLoaded)
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

	entities, err := ts.crudCountryService.GetEntities(
		conditions.CountryPreloadCapital,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Country{country1, country2}, entities)
	ts.True(pie.Any(entities, func(country *models.Country) bool {
		countryCapital, err := country.GetCapital()
		return err == nil && countryCapital.Equal(capital1)
	}))
	ts.True(pie.Any(entities, func(country *models.Country) bool {
		countryCapital, err := country.GetCapital()
		return err == nil && countryCapital.Equal(capital2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadHasManyReversed() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	entities, err := ts.crudSellerService.GetEntities(
		conditions.SellerPreloadCompany,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{seller1, seller2}, entities)
	ts.True(pie.Any(entities, func(seller *models.Seller) bool {
		sellerCompany, err := seller.GetCompany()
		return err == nil && sellerCompany.Equal(*company1)
	}))
	ts.True(pie.Any(entities, func(seller *models.Seller) bool {
		sellerCompany, err := seller.GetCompany()
		return err == nil && sellerCompany.Equal(*company2)
	}))
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
	saleProduct, err := entities[0].GetProduct()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), product1, saleProduct)

	saleSeller, err := entities[0].GetSeller()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), seller1, saleSeller)
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildPreloadParent1,
		conditions.ChildPreloadParent2,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	childParent1, err := entities[0].GetParent1()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent1, childParent1)

	childParent2, err := entities[0].GetParent2()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent2, childParent2)
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildPreloadRelations...,
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	childParent1, err := entities[0].GetParent1()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent1, childParent1)

	childParent2, err := entities[0].GetParent2()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent2, childParent2)
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildParent1(
			conditions.Parent1PreloadAttributes,
			conditions.Parent1PreloadParentParent,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	childParent1, err := entities[0].GetParent1()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent1, childParent1)

	childParentParent, err := childParent1.GetParentParent()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parentParent, childParentParent)
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
	childParent1, err := entities[0].GetParent1()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent11, childParent1)

	childParentParent, err := childParent1.GetParentParent()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parentParent1, childParentParent)
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

	entities, err := ts.crudChildService.GetEntities(
		conditions.ChildParent1(
			conditions.Parent1PreloadParentParent,
		),
		conditions.ChildParent2(
			conditions.Parent2PreloadParentParent,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
	childParent1, err := entities[0].GetParent1()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent1, childParent1)

	childParent2, err := entities[0].GetParent2()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parent2, childParent2)

	childParent1Parent, err := childParent1.GetParentParent()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parentParent, childParent1Parent)

	childParent2Parent, err := childParent2.GetParentParent()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), parentParent, childParent2Parent)
}

// TODO generacion automatica
func CompanyPreloadSellers(nestedPreloads ...badorm.IJoinCondition[models.Seller]) badorm.Condition[models.Company] {
	return badorm.NewCollectionPreloadCondition[models.Company, models.Seller](
		"Sellers",
		nestedPreloads,
	)
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadList() {
	company := ts.createCompany("ditrit")
	seller1 := ts.createSeller("1", company)
	seller2 := ts.createSeller("2", company)

	entities, err := ts.crudCompanyService.GetEntities(
		CompanyPreloadSellers(),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Company{company}, entities)
	// TODO tener el getter
	EqualList(&ts.Suite, []models.Seller{*seller1, *seller2}, entities[0].Sellers)
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadListAndNestedAttributes() {
	company := ts.createCompany("ditrit")

	university1 := ts.createUniversity("uni1")
	seller1 := ts.createSeller("1", company)
	seller1.University = university1
	err := ts.db.Save(seller1).Error
	ts.Nil(err)

	university2 := ts.createUniversity("uni1")
	seller2 := ts.createSeller("2", company)
	seller2.University = university2
	err = ts.db.Save(seller2).Error
	ts.Nil(err)

	entities, err := ts.crudCompanyService.GetEntities(
		CompanyPreloadSellers(
			conditions.SellerPreloadUniversity,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Company{company}, entities)
	// TODO tener el getter
	EqualList(&ts.Suite, []models.Seller{*seller1, *seller2}, entities[0].Sellers)

	ts.True(pie.Any(entities[0].Sellers, func(seller models.Seller) bool {
		sellerUniversity, err := seller.GetUniversity()
		return err == nil && sellerUniversity.Equal(*university1)
	}))
	ts.True(pie.Any(entities[0].Sellers, func(seller models.Seller) bool {
		sellerUniversity, err := seller.GetUniversity()
		return err == nil && sellerUniversity.Equal(*university2)
	}))
}

func (ts *PreloadConditionsIntTestSuite) TestPreloadMultipleListsAndNestedAttributes() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	university1 := ts.createUniversity("uni1")
	seller1 := ts.createSeller("1", company1)
	seller1.University = university1
	err := ts.db.Save(seller1).Error
	ts.Nil(err)

	university2 := ts.createUniversity("uni1")
	seller2 := ts.createSeller("2", company1)
	seller2.University = university2
	err = ts.db.Save(seller2).Error
	ts.Nil(err)

	seller3 := ts.createSeller("3", company2)
	seller3.University = university1
	err = ts.db.Save(seller3).Error
	ts.Nil(err)

	seller4 := ts.createSeller("4", company2)
	seller4.University = university2
	err = ts.db.Save(seller4).Error
	ts.Nil(err)

	entities, err := ts.crudCompanyService.GetEntities(
		CompanyPreloadSellers(
			conditions.SellerPreloadUniversity,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Company{company1, company2}, entities)

	company1Loaded := *utils.FindFirst(entities, func(company *models.Company) bool {
		return company.Equal(*company1)
	})
	company2Loaded := *utils.FindFirst(entities, func(company *models.Company) bool {
		return company.Equal(*company2)
	})

	// TODO tener el getter
	EqualList(&ts.Suite, []models.Seller{*seller1, *seller2}, company1Loaded.Sellers)

	ts.True(pie.Any(company1Loaded.Sellers, func(seller models.Seller) bool {
		sellerUniversity, err := seller.GetUniversity()
		return err == nil && sellerUniversity.Equal(*university1)
	}))
	ts.True(pie.Any(company1Loaded.Sellers, func(seller models.Seller) bool {
		sellerUniversity, err := seller.GetUniversity()
		return err == nil && sellerUniversity.Equal(*university2)
	}))

	// TODO tener el getter
	EqualList(&ts.Suite, []models.Seller{*seller3, *seller4}, company2Loaded.Sellers)

	ts.True(pie.Any(company2Loaded.Sellers, func(seller models.Seller) bool {
		sellerUniversity, err := seller.GetUniversity()
		return err == nil && sellerUniversity.Equal(*university1)
	}))
	ts.True(pie.Any(company2Loaded.Sellers, func(seller models.Seller) bool {
		sellerUniversity, err := seller.GetUniversity()
		return err == nil && sellerUniversity.Equal(*university2)
	}))
}
