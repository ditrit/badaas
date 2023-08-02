package testintegration

import (
	"database/sql"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type OperatorsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService orm.CRUDService[models.Product, orm.UUID]
}

func NewOperatorsIntTestSuite(
	db *gorm.DB,
	crudProductService orm.CRUDService[models.Product, orm.UUID],
) *OperatorsIntTestSuite {
	return &OperatorsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
	}
}

func (ts *OperatorsIntTestSuite) TestEqNullableNullReturnsError() {
	_, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.Eq(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, orm.ErrValueCantBeNull)
}

func (ts *OperatorsIntTestSuite) TestEqPointers() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)

	intNotMatch := 2
	ts.createProduct("match", 3, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			orm.Eq(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestNotEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.NotEq(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLt() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.Lt(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLtNullableNullReturnsError() {
	_, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.Lt(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, orm.ErrValueCantBeNull)
}

func (ts *OperatorsIntTestSuite) TestLtOrEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.LtOrEq(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestGt() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.Gt(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestGtOrEq() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.GtOrEq(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}
