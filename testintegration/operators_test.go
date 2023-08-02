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
