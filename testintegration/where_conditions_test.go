package testintegration

import (
	"log"

	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/mysql"
	"github.com/ditrit/badaas/badorm/sqlserver"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type WhereConditionsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService badorm.CRUDService[models.Product, badorm.UUID]
	crudSaleService    badorm.CRUDService[models.Sale, badorm.UUID]
	crudBrandService   badorm.CRUDService[models.Brand, uint]
}

func NewWhereConditionsIntTestSuite(
	db *gorm.DB,
	crudProductService badorm.CRUDService[models.Product, badorm.UUID],
	crudSaleService badorm.CRUDService[models.Sale, badorm.UUID],
	crudBrandService badorm.CRUDService[models.Brand, uint],
) *WhereConditionsIntTestSuite {
	return &WhereConditionsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
		crudSaleService:    crudSaleService,
		crudBrandService:   crudBrandService,
	}
}

// ------------------------- GetEntity --------------------------------

func (ts *WhereConditionsIntTestSuite) TestGetEntityReturnsErrorIfNotEntityCreated() {
	_, err := ts.crudProductService.GetEntity(badorm.NilUUID)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *WhereConditionsIntTestSuite) TestGetEntityReturnsErrorIfNotEntityMatch() {
	ts.createProduct("", 0, 0, false, nil)

	_, err := ts.crudProductService.GetEntity(badorm.NewUUID())
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *WhereConditionsIntTestSuite) TestGetEntityReturnsTheEntityIfItIsCreate() {
	match := ts.createProduct("", 0, 0, false, nil)

	entity, err := ts.crudProductService.GetEntity(match.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), match, entity)
}

// ------------------------- GetEntities --------------------------------

func (ts *WhereConditionsIntTestSuite) TestGetEntitiesReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestGetEntitiesReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestGetEntitiesReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match2 := ts.createProduct("", 0, 0, false, nil)
	match3 := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.Eq("not_created"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct("something_else", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.Eq("not_match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsOneIfOnlyOneMatch() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsMultipleIfMultipleMatch() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfIntType() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfFloatType() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfBoolType() {
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

func (ts *WhereConditionsIntTestSuite) TestMultipleConditionsOfDifferentTypesWorks() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfID() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfCreatedAt() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductCreatedAt(badorm.Eq(match.CreatedAt)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestDeletedAtConditionIsAddedAutomatically() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	deleted := ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(deleted).Error)

	entities, err := ts.crudProductService.GetEntities()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

// TODO DeletedAt with nil value but not automatic

func (ts *WhereConditionsIntTestSuite) TestConditionOfDeletedAtNotNil() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(match).Error)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductDeletedAt(badorm.Eq(match.DeletedAt)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfEmbedded() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfGormEmbedded() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfPointerTypeWithValue() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)
	intNotMatch := 2
	ts.createProduct("not_match", 2, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(badorm.Eq(1)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfByteArrayWithContent() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfByteArrayEmpty() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfCustomType() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfRelationType() {
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

func (ts *WhereConditionsIntTestSuite) TestConditionOfRelationTypeOptionalWithValue() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSellerId(badorm.Eq(seller1.ID)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfRelationTypeOptionalByNil() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.GetEntities(
		conditions.SaleSellerId(badorm.IsNull[badorm.UUID]()),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsOnUIntModel() {
	match := ts.createBrand("match")
	ts.createBrand("not_match")

	entities, err := ts.crudBrandService.GetEntities(
		conditions.BrandName(badorm.Eq("match")),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Brand{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestMultipleConditionsAreConnectedByAnd() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 5, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(badorm.GtOrEq(3)),
		conditions.ProductInt(badorm.LtOrEq(4)),
		conditions.ProductString(badorm.Eq("match")),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestNot() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)

	ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		badorm.Not(conditions.ProductInt(badorm.Eq(2))),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestNotWithMultipleConditionsAreConnectedByAnd() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 5, 0, false, nil)

	ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		badorm.Not(
			conditions.ProductInt(badorm.Gt(1)),
			conditions.ProductInt(badorm.Lt(4)),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestOr() {
	match1 := ts.createProduct("match", 2, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)
	match3 := ts.createProduct("match_3", 3, 0, false, nil)

	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		badorm.Or(
			conditions.ProductInt(badorm.Eq(2)),
			conditions.ProductInt(badorm.Eq(3)),
			conditions.ProductString(badorm.Eq("match_3")),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestNotOr() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 5, 0, false, nil)
	match3 := ts.createProduct("match", 4, 0, false, nil)

	ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match_string", 3, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		badorm.Not[models.Product](
			badorm.Or(
				conditions.ProductInt(badorm.Eq(2)),
				conditions.ProductString(badorm.Eq("not_match_string")),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestXor() {
	match1 := ts.createProduct("", 1, 0, false, nil)
	match2 := ts.createProduct("", 7, 0, false, nil)

	ts.createProduct("", 5, 0, false, nil)
	ts.createProduct("", 4, 0, false, nil)

	var xorExpression badorm.WhereCondition[models.Product]

	switch getDBDialector() {
	case configuration.PostgreSQL, configuration.SQLite:
		log.Println("Xor not compatible")
	case configuration.MySQL:
		xorExpression = mysql.Xor(
			conditions.ProductInt(badorm.Lt(6)),
			conditions.ProductInt(badorm.Gt(3)),
		)
	case configuration.SQLServer:
		xorExpression = sqlserver.Xor(
			conditions.ProductInt(badorm.Lt(6)),
			conditions.ProductInt(badorm.Gt(3)),
		)
	}

	if xorExpression != nil {
		entities, err := ts.crudProductService.GetEntities(xorExpression)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *WhereConditionsIntTestSuite) TestMultipleConditionsDifferentExpressions() {
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
