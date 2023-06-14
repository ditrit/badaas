package testintegration

import (
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type CRUDServiceIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService  badorm.CRUDService[models.Product, badorm.UUID]
	crudSaleService     badorm.CRUDService[models.Sale, badorm.UUID]
	crudSellerService   badorm.CRUDService[models.Seller, badorm.UUID]
	crudCountryService  badorm.CRUDService[models.Country, badorm.UUID]
	crudCityService     badorm.CRUDService[models.City, badorm.UUID]
	crudEmployeeService badorm.CRUDService[models.Employee, badorm.UUID]
	crudBicycleService  badorm.CRUDService[models.Bicycle, badorm.UUID]
	crudBrandService    badorm.CRUDService[models.Brand, uint]
	crudPhoneService    badorm.CRUDService[models.Phone, uint]
}

func NewCRUDServiceIntTestSuite(
	db *gorm.DB,
	crudProductService badorm.CRUDService[models.Product, badorm.UUID],
	crudSaleService badorm.CRUDService[models.Sale, badorm.UUID],
	crudSellerService badorm.CRUDService[models.Seller, badorm.UUID],
	crudCountryService badorm.CRUDService[models.Country, badorm.UUID],
	crudCityService badorm.CRUDService[models.City, badorm.UUID],
	crudEmployeeService badorm.CRUDService[models.Employee, badorm.UUID],
	crudBicycleService badorm.CRUDService[models.Bicycle, badorm.UUID],
	crudBrandService badorm.CRUDService[models.Brand, uint],
	crudPhoneService badorm.CRUDService[models.Phone, uint],
) *CRUDServiceIntTestSuite {
	return &CRUDServiceIntTestSuite{
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

// ------------------------- GetEntity --------------------------------

func (ts *CRUDServiceIntTestSuite) TestGetEntityReturnsErrorIfNotEntityCreated() {
	_, err := ts.crudProductService.GetEntity(badorm.NilUUID)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntityReturnsErrorIfNotEntityMatch() {
	ts.createProduct("", 0, 0, false, nil)

	_, err := ts.crudProductService.GetEntity(badorm.NewUUID())
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntityReturnsTheEntityIfItIsCreate() {
	match := ts.createProduct("", 0, 0, false, nil)

	entity, err := ts.crudProductService.GetEntity(match.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), match, entity)
}

// ------------------------- GetEntities --------------------------------

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match2 := ts.createProduct("", 0, 0, false, nil)
	match3 := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.Eq("not_created"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct("something_else", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.Eq("not_match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.Eq("match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.Eq("match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.Eq(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntTypeNotEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.NotEq(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntTypeLt() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.Lt(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntTypeLtOrEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.LtOrEq(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntTypeGt() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.Gt(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntTypeGtOrEq() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.GtOrEq(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionWithExprIsDistinct() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.IsDistinct(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionWithExprIsNotDistinct() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.IsNotDistinct(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionWithExprIsNull() {
	match := ts.createProduct("match", 0, 0, false, nil)
	int1 := 1
	int2 := 2
	ts.createProduct("not_match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, &int2)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNull[*int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionWithExprIsNotNull() {
	int1 := 1
	match := ts.createProduct("match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNotNull[*int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionWithMultipleExpressions() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 5, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.GtOrEq(3),
			badorm.LtOrEq(4),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfFloatType() {
	match := ts.createProduct("match", 0, 1.1, false, nil)
	ts.createProduct("not_match", 0, 2.2, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductFloat(
			badorm.Eq(1.1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfBoolType() {
	match := ts.createProduct("match", 0, 0.0, true, nil)
	ts.createProduct("not_match", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			badorm.Eq(true),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 1, 0.0, true, nil)
	ts.createProduct("match", 2, 0.0, true, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(badorm.Eq("match")),
		conditions.ProductInt(badorm.Eq(1)),
		conditions.ProductBool(badorm.Eq(true)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithMultipleConditionsOfDifferentTypesWithDifferentExpressionsWorks() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 1, 0.0, true, nil)
	ts.createProduct("match", 2, 0.0, true, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(badorm.Eq("match")),
		conditions.ProductInt(badorm.Lt(2)),
		conditions.ProductBool(badorm.NotEq(false)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfID() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductId(
			badorm.Eq(match.ID),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfCreatedAt() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductCreatedAt(badorm.Eq(match.CreatedAt)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesDeletedAtConditionIsAddedAutomatically() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	deleted := ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(deleted).Error)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

// TODO DeletedAt with nil value but not automatic

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfDeletedAtNotNil() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(match).Error)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductDeletedAt(badorm.Eq(match.DeletedAt)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfEmbedded() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	match.EmbeddedInt = 1

	err := ts.db.Save(match).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductEmbeddedInt(badorm.Eq(1)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfGormEmbedded() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	match.GormEmbedded.Int = 1

	err := ts.db.Save(match).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductGormEmbeddedInt(badorm.Eq(1)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfPointerTypeWithValue() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)
	intNotMatch := 2
	ts.createProduct("not_match", 2, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(badorm.Eq(&intMatch)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfPointerTypeByNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	intNotMatch := 2
	ts.createProduct("not_match", 2, 0, false, &intNotMatch)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(badorm.Eq[*int](nil)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfByteArrayWithContent() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	match.ByteArray = []byte{1, 2}
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductByteArray(badorm.Eq([]byte{1, 2})),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfByteArrayEmpty() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	match.ByteArray = []byte{}
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductByteArray(badorm.Eq([]byte{})),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfByteArrayNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductByteArray(badorm.Eq[[]uint8](nil)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfCustomType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	match.MultiString = models.MultiString{"salut", "hola"}
	notMatch1.MultiString = models.MultiString{"salut", "hola", "hello"}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductMultiString(badorm.Eq(models.MultiString{"salut", "hola"})),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfRelationType() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProductId(badorm.Eq(product1.ID)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfRelationTypeOptionalWithValue() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSellerId(badorm.Eq(&seller1.ID)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfRelationTypeOptionalByNil() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSellerId(badorm.Eq[*badorm.UUID](nil)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsOnUIntModel() {
	match := ts.createBrand("match")
	ts.createBrand("not_match")

	entities, err := ts.crudBrandService.GetEntities(
		conditions.BrandName(badorm.Eq("match")),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Brand{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsUintBelongsTo() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsBelongsTo() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsAndFiltersTheMainEntity() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsHasOneOptional() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsHasOneSelfReferential() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOneToOne() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOneToOneReversed() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsWithEntityThatDefinesTableName() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnHasMany() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnDifferentAttributes() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsAddsDeletedAtAutomatically() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnDeletedAt() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsAndFiltersByNil() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	intProduct2 := 2
	product2 := ts.createProduct("", 2, 0.0, false, &intProduct2)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleProduct(
			conditions.ProductIntPointer(badorm.Eq[*int](nil)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsDifferentEntities() {
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsMultipleTimes() {
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
