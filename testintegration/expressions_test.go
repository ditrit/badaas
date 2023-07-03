package testintegration

import (
	"database/sql"
	"log"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/mysql"
	"github.com/ditrit/badaas/badorm/psql"
	"github.com/ditrit/badaas/badorm/sqlite"
	"github.com/ditrit/badaas/badorm/sqlserver"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type ExpressionIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService badorm.CRUDService[models.Product, badorm.UUID]
}

func NewExpressionsIntTestSuite(
	db *gorm.DB,
	crudProductService badorm.CRUDService[models.Product, badorm.UUID],
) *ExpressionIntTestSuite {
	return &ExpressionIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
	}
}

func (ts *ExpressionIntTestSuite) TestEqNullableNullReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.Eq(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, badorm.ErrValueCantBeNull)
}

func (ts *ExpressionIntTestSuite) TestEqPointers() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)

	intNotMatch := 2
	ts.createProduct("match", 3, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(
			badorm.Eq(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullTNotNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.EqOrIsNull[int](1),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductByteArray(
			badorm.EqOrIsNull[[]byte](nil),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductByteArray(
			badorm.EqOrIsNull[[]byte](nilOfType),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(
			badorm.EqOrIsNull[int](intPointer),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNotNilPointer() {
	matchInt := 1
	match := ts.createProduct("match", 1, 0, false, &matchInt)

	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.EqOrIsNull[int](&matchInt),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.EqOrIsNull[sql.NullFloat64](sql.NullFloat64{Valid: false}),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.EqOrIsNull[sql.NullFloat64](sql.NullFloat64{Valid: true, Float64: 6}),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestEqOrIsNullNotRelated() {
	notRelated := "not_related"

	_, err := ts.crudProductService.GetEntities(
		conditions.ProductFloat(
			badorm.EqOrIsNull[float64](&notRelated),
		),
	)
	ts.ErrorIs(err, badorm.ErrNotRelated)
}

func (ts *ExpressionIntTestSuite) TestNotEqOrIsNotNullTNotNil() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("match", 3, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.NotEqOrIsNotNull[int](3),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductByteArray(
			badorm.NotEqOrIsNotNull[[]byte](nil),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestNotEq() {
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

func (ts *ExpressionIntTestSuite) TestLt() {
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

func (ts *ExpressionIntTestSuite) TestLtNullableNullReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.Lt(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, badorm.ErrValueCantBeNull)
}

func (ts *ExpressionIntTestSuite) TestLtOrEq() {
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

func (ts *ExpressionIntTestSuite) TestNotLt() {
	switch getDBDialector() {
	case configuration.SQLServer:
		match1 := ts.createProduct("match", 3, 0, false, nil)
		match2 := ts.createProduct("match", 4, 0, false, nil)
		ts.createProduct("not_match", 1, 0, false, nil)
		ts.createProduct("not_match", 2, 0, false, nil)

		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductInt(
				sqlserver.NotLt(3),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case configuration.PostgreSQL, configuration.MySQL, configuration.SQLite:
		log.Println("NotLt not supported")
	}
}

func (ts *ExpressionIntTestSuite) TestGt() {
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

func (ts *ExpressionIntTestSuite) TestGtOrEq() {
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

func (ts *ExpressionIntTestSuite) TestNotGt() {
	switch getDBDialector() {
	case configuration.SQLServer:
		match1 := ts.createProduct("match", 1, 0, false, nil)
		match2 := ts.createProduct("match", 2, 0, false, nil)
		ts.createProduct("not_match", 3, 0, false, nil)
		ts.createProduct("not_match", 4, 0, false, nil)

		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductInt(
				sqlserver.NotGt(2),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case configuration.PostgreSQL, configuration.MySQL, configuration.SQLite:
		log.Println("NotGt not supported")
	}
}

func (ts *ExpressionIntTestSuite) TestBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 6, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.Between(3, 5),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.NotBetween(0, 2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsDistinct() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	switch getDBDialector() {
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductInt(
				badorm.IsDistinct(2),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case configuration.MySQL:
		entities, err := ts.crudProductService.GetEntities(
			badorm.Not[models.Product](
				conditions.ProductInt(mysql.IsEqual(2)),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *ExpressionIntTestSuite) TestIsNotDistinct() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	var isNotEqualExpression badorm.Expression[int]

	switch getDBDialector() {
	case configuration.MySQL:
		isNotEqualExpression = mysql.IsEqual(3)
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		isNotEqualExpression = badorm.IsNotDistinct(3)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			isNotEqualExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNotDistinctNullValue() {
	match := ts.createProduct("match", 3, 0, false, nil)

	notMatch := ts.createProduct("not_match", 4, 0, false, nil)
	notMatch.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	var isEqualExpression badorm.Expression[sql.NullFloat64]

	switch getDBDialector() {
	case configuration.MySQL:
		isEqualExpression = mysql.IsEqual(sql.NullFloat64{Valid: false})
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		isEqualExpression = badorm.IsNotDistinct(sql.NullFloat64{Valid: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			isEqualExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsNull() {
	match := ts.createProduct("match", 0, 0, false, nil)
	int1 := 1
	int2 := 2

	ts.createProduct("not_match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, &int2)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(
			badorm.IsNull[int](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.IsNull[sql.NullFloat64](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductIntPointer(
			badorm.IsNotNull[int](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.IsNotNull[sql.NullFloat64](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsTrue() {
	match := ts.createProduct("match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	var isTrueExpression badorm.Expression[bool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isTrueExpression = badorm.IsTrue[bool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsTrue
		isTrueExpression = badorm.Eq[bool](true)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			isTrueExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	var isFalseExpression badorm.Expression[bool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isFalseExpression = badorm.IsFalse[bool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsFalse
		isFalseExpression = badorm.Eq[bool](false)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			isFalseExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

//nolint:dupl // not really duplicated
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

	var isNotTrueExpression badorm.Expression[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isNotTrueExpression = badorm.IsNotTrue[sql.NullBool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsNotTrue
		isNotTrueExpression = badorm.IsDistinct[sql.NullBool](sql.NullBool{Valid: true, Bool: true})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isNotTrueExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

//nolint:dupl // not really duplicated
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

	var isNotFalseExpression badorm.Expression[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isNotFalseExpression = badorm.IsNotFalse[sql.NullBool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsNotFalse
		isNotFalseExpression = badorm.IsDistinct[sql.NullBool](sql.NullBool{Valid: true, Bool: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isNotFalseExpression,
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

	var isUnknownExpression badorm.Expression[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL:
		isUnknownExpression = badorm.IsUnknown[sql.NullBool]()
	case configuration.SQLServer, configuration.SQLite:
		// sqlserver doesn't support IsUnknown
		isUnknownExpression = badorm.IsNotDistinct[sql.NullBool](sql.NullBool{Valid: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isUnknownExpression,
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

	var isNotUnknownExpression badorm.Expression[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL:
		isNotUnknownExpression = badorm.IsNotUnknown[sql.NullBool]()
	case configuration.SQLServer, configuration.SQLite:
		// sqlserver doesn't support IsNotUnknown
		isNotUnknownExpression = badorm.IsDistinct[sql.NullBool](sql.NullBool{Valid: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isNotUnknownExpression,
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

	var arrayInExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.MySQL:
		arrayInExpression = mysql.ArrayIn("s1", "s2", "s3")
	case configuration.PostgreSQL:
		arrayInExpression = psql.ArrayIn("s1", "s2", "s3")
	case configuration.SQLServer:
		arrayInExpression = sqlserver.ArrayIn("s1", "s2", "s3")
	case configuration.SQLite:
		arrayInExpression = sqlite.ArrayIn("s1", "s2", "s3")
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			arrayInExpression,
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

	var arrayNotInExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.MySQL:
		arrayNotInExpression = mysql.ArrayNotIn("ns1", "ns2")
	case configuration.PostgreSQL:
		arrayNotInExpression = psql.ArrayNotIn("ns1", "ns2")
	case configuration.SQLServer:
		arrayNotInExpression = sqlserver.ArrayNotIn("ns1", "ns2")
	case configuration.SQLite:
		arrayNotInExpression = sqlite.ArrayNotIn("ns1", "ns2")
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			arrayNotInExpression,
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

	var likeExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		likeExpression = badorm.Like[string]("_a%")
	case configuration.SQLServer:
		likeExpression = badorm.Like[string]("[bc]a[^a]%")
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			likeExpression,
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

	var likeExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		likeExpression = badorm.Like[string]("_a!_%").Escape('!')
	case configuration.SQLServer:
		likeExpression = badorm.Like[string]("[bc]a!_[^a]%").Escape('!')
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			likeExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestLikeOnNumeric() {
	switch getDBDialector() {
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		log.Println("Like with numeric not compatible")
	case configuration.MySQL:
		match1 := ts.createProduct("", 10, 0, false, nil)
		match2 := ts.createProduct("", 100, 0, false, nil)

		ts.createProduct("", 20, 0, false, nil)
		ts.createProduct("", 3, 0, false, nil)

		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductInt(
				mysql.Like[int]("1%"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *ExpressionIntTestSuite) TestILike() {
	switch getDBDialector() {
	case configuration.MySQL, configuration.SQLServer, configuration.SQLite:
		log.Println("ILike not compatible")
	case configuration.PostgreSQL:
		match1 := ts.createProduct("basd", 0, 0, false, nil)
		match2 := ts.createProduct("cape", 0, 0, false, nil)
		match3 := ts.createProduct("bAsd", 0, 0, false, nil)

		ts.createProduct("bbsd", 0, 0, false, nil)
		ts.createProduct("bbasd", 0, 0, false, nil)

		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				psql.ILike[string]("_a%"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
}

func (ts *ExpressionIntTestSuite) TestSimilarTo() {
	switch getDBDialector() {
	case configuration.MySQL, configuration.SQLServer, configuration.SQLite:
		log.Println("SimilarTo not compatible")
	case configuration.PostgreSQL:
		match1 := ts.createProduct("abc", 0, 0, false, nil)
		match2 := ts.createProduct("aabcc", 0, 0, false, nil)

		ts.createProduct("aec", 0, 0, false, nil)
		ts.createProduct("aaaaa", 0, 0, false, nil)

		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				psql.SimilarTo[string]("%(b|d)%"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *ExpressionIntTestSuite) TestPosixRegexCaseSensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)
	ts.createProduct("AB", 0, 0, false, nil)

	var posixRegexExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.SQLServer, configuration.MySQL:
		log.Println("PosixRegex not compatible")
	case configuration.PostgreSQL:
		posixRegexExpression = psql.POSIXMatch[string]("^a(b|x)")
	case configuration.SQLite:
		posixRegexExpression = sqlite.Glob[string]("a[bx]")
	}

	if posixRegexExpression != nil {
		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				posixRegexExpression,
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *ExpressionIntTestSuite) TestPosixRegexCaseInsensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)
	match3 := ts.createProduct("AB", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)

	var posixRegexExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.SQLServer, configuration.SQLite:
		log.Println("PosixRegex Case Insensitive not compatible")
	case configuration.MySQL:
		posixRegexExpression = mysql.RegexP[string]("^a(b|x)")
	case configuration.PostgreSQL:
		posixRegexExpression = psql.POSIXIMatch[string]("^a(b|x)")
	}

	if posixRegexExpression != nil {
		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				posixRegexExpression,
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
}

func (ts *ExpressionIntTestSuite) TestPosixRegexNotPosix() {
	var posixRegexExpression badorm.Expression[string]

	switch getDBDialector() {
	case configuration.SQLServer:
		log.Println("PosixRegex not compatible")
	case configuration.MySQL:
		posixRegexExpression = mysql.RegexP[string]("^a(b|x")
	case configuration.PostgreSQL:
		posixRegexExpression = psql.POSIXMatch[string]("^a(b|x")
	case configuration.SQLite:
		posixRegexExpression = sqlite.Glob[string]("^a(b|x")
	}

	if posixRegexExpression != nil {
		_, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				posixRegexExpression,
			),
		)
		ts.ErrorContains(err, "error parsing regexp")
	}
}

func (ts *ExpressionIntTestSuite) TestDynamicExpressionOver1Table() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			badorm.NewDynamicExpression(
				badorm.Eq[int],
				conditions.ProductIntPointerField,
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}
