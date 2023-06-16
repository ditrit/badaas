package testintegration

import (
	"database/sql"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type ExpressionIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService orm.CRUDService[models.Product, orm.UUID]
}

func NewExpressionsIntTestSuite(
	db *gorm.DB,
	crudProductService orm.CRUDService[models.Product, orm.UUID],
) *ExpressionIntTestSuite {
	return &ExpressionIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
	}
}

func (ts *ExpressionIntTestSuite) TestEqNullableNullReturnsError() {
	_, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.Eq(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, orm.ErrValueCantBeNull)
}

func (ts *ExpressionIntTestSuite) TestEqPointers() {
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

func (ts *ExpressionIntTestSuite) TestEqOrIsNullTNotNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.EqOrIsNull[int](1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullTNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch := ts.createProduct("match", 3, 0, false, nil)
	notMatch.ByteArray = []byte{2, 3}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductByteArray(
			orm.EqOrIsNull[[]byte](nil),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullTNilOfType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch := ts.createProduct("match", 3, 0, false, nil)
	notMatch.ByteArray = []byte{2, 3}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	var nilOfType []byte

	entities, err := ts.crudProductService.Query(
		conditions.ProductByteArray(
			orm.EqOrIsNull[[]byte](nilOfType),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNilPointer() {
	match := ts.createProduct("match", 1, 0, false, nil)

	notMatchInt := 1
	ts.createProduct("match", 3, 0, false, &notMatchInt)

	var intPointer *int

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			orm.EqOrIsNull[int](intPointer),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNotNilPointer() {
	matchInt := 1
	match := ts.createProduct("match", 1, 0, false, &matchInt)

	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.EqOrIsNull[int](&matchInt),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNullableNil() {
	match := ts.createProduct("match", 1, 0, false, nil)

	notMatch := ts.createProduct("match", 3, 0, false, nil)
	notMatch.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.EqOrIsNull[sql.NullFloat64](sql.NullFloat64{Valid: false}),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNullableNotNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.EqOrIsNull[sql.NullFloat64](sql.NullFloat64{Valid: true, Float64: 6}),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNotRelated() {
	notRelated := "not_related"

	_, err := ts.crudProductService.Query(
		conditions.ProductFloat(
			orm.EqOrIsNull[float64](&notRelated),
		),
	)
	ts.ErrorIs(err, orm.ErrNotRelated)
}

func (ts *ExpressionIntTestSuite) TestNotEqOrIsNotNullTNotNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.NotEqOrIsNotNull[int](3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestNotEqOrIsNotNullTNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	match.ByteArray = []byte{2, 3}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductByteArray(
			orm.NotEqOrIsNotNull[[]byte](nil),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestNotEq() {
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

func (ts *ExpressionIntTestSuite) TestLt() {
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

func (ts *ExpressionIntTestSuite) TestLtNullableNullReturnsError() {
	_, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.Lt(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, orm.ErrValueCantBeNull)
}

func (ts *ExpressionIntTestSuite) TestLtOrEq() {
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

func (ts *ExpressionIntTestSuite) TestGt() {
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

func (ts *ExpressionIntTestSuite) TestGtOrEq() {
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

func (ts *ExpressionIntTestSuite) TestBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 6, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.Between(3, 5),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestNotBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.NotBetween(0, 2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNull() {
	match := ts.createProduct("match", 0, 0, false, nil)
	int1 := 1
	int2 := 2

	ts.createProduct("not_match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, &int2)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNull[int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNullNotPointers() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNull[sql.NullFloat64](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNotNull() {
	int1 := 1
	match := ts.createProduct("match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNotNull[int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNotNullNotPointers() {
	match := ts.createProduct("match", 0, 0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNotNull[sql.NullFloat64](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsTrue() {
	match := ts.createProduct("match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsTrue[bool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsFalse[bool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNotTrue() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err := ts.db.Save(match2).Error
	ts.Nil(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: true}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNotTrue[sql.NullBool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNotFalse() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match2).Error
	ts.Nil(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNotFalse[sql.NullBool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsUnknown() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch1 := ts.createProduct("match", 0, 0, false, nil)
	notMatch1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(notMatch1).Error
	ts.Nil(err)

	notMatch2 := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch2).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsUnknown[sql.NullBool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNotUnknown() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match1).Error
	ts.Nil(err)

	match2 := ts.createProduct("", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(match2).Error
	ts.Nil(err)

	ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			orm.IsNotUnknown[sql.NullBool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestArrayIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.ArrayIn("s1", "s2", "s3"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestArrayNotIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.ArrayNotIn("ns1", "ns2"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestLike() {
	match1 := ts.createProduct("basd", 0, 0, false, nil)
	match2 := ts.createProduct("cape", 0, 0, false, nil)

	ts.createProduct("bbsd", 0, 0, false, nil)
	ts.createProduct("bbasd", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.Like[string]("_a%"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestLikeEscape() {
	match1 := ts.createProduct("ba_sd", 0, 0, false, nil)
	match2 := ts.createProduct("ca_pe", 0, 0, false, nil)

	ts.createProduct("bb_sd", 0, 0, false, nil)
	ts.createProduct("bba_sd", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.Like[string]("_a!_%").Escape('!'),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}
