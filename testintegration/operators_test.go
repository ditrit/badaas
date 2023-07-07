package testintegration

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/dynamic"
	"github.com/ditrit/badaas/badorm/dynamic/mysqldynamic"
	"github.com/ditrit/badaas/badorm/multitype"
	"github.com/ditrit/badaas/badorm/multitype/mysqlmultitype"
	"github.com/ditrit/badaas/badorm/mysql"
	"github.com/ditrit/badaas/badorm/psql"
	"github.com/ditrit/badaas/badorm/sqlite"
	"github.com/ditrit/badaas/badorm/sqlserver"
	"github.com/ditrit/badaas/badorm/unsafe"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type OperatorIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService badorm.CRUDService[models.Product, badorm.UUID]
}

func NewOperatorIntTestSuite(
	db *gorm.DB,
	crudProductService badorm.CRUDService[models.Product, badorm.UUID],
) *OperatorIntTestSuite {
	return &OperatorIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
	}
}

func (ts *OperatorIntTestSuite) TestEqNullableNullReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.Eq(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, badorm.ErrValueCantBeNull)
	ts.ErrorContains(err, "operator: Eq; model: models.Product, field: NullFloat")
}

func (ts *OperatorIntTestSuite) TestEqPointers() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullTNotNil() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullTNil() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullTNilOfType() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullNilPointer() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullNotNilPointer() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullNullableNil() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullNullableNotNil() {
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

func (ts *OperatorIntTestSuite) TestEqOrIsNullNotRelated() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductFloat(
			badorm.EqOrIsNull[float64]("not_related"),
		),
	)
	ts.ErrorIs(err, badorm.ErrNotRelated)
	ts.ErrorContains(err, "type: string, T: float64; operator: EqOrIsNull; model: models.Product, field: Float")
}

func (ts *OperatorIntTestSuite) TestNotEqOrIsNotNullTNotNil() {
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

func (ts *OperatorIntTestSuite) TestNotEqOrIsNotNullTNil() {
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

func (ts *OperatorIntTestSuite) TestNotEq() {
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

func (ts *OperatorIntTestSuite) TestLt() {
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

func (ts *OperatorIntTestSuite) TestLtNullableNullReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			badorm.Lt(sql.NullFloat64{Valid: false}),
		),
	)
	ts.ErrorIs(err, badorm.ErrValueCantBeNull)
	ts.ErrorContains(err, "operator: Lt; model: models.Product, field: NullFloat")
}

func (ts *OperatorIntTestSuite) TestLtOrEq() {
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

func (ts *OperatorIntTestSuite) TestNotLt() {
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

func (ts *OperatorIntTestSuite) TestGt() {
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

func (ts *OperatorIntTestSuite) TestGtOrEq() {
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

func (ts *OperatorIntTestSuite) TestNotGt() {
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

func (ts *OperatorIntTestSuite) TestBetween() {
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

func (ts *OperatorIntTestSuite) TestNotBetween() {
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

func (ts *OperatorIntTestSuite) TestIsDistinct() {
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

func (ts *OperatorIntTestSuite) TestIsNotDistinct() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	var isNotEqualOperator badorm.Operator[int]

	switch getDBDialector() {
	case configuration.MySQL:
		isNotEqualOperator = mysql.IsEqual(3)
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		isNotEqualOperator = badorm.IsNotDistinct(3)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			isNotEqualOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestIsNotDistinctNullValue() {
	match := ts.createProduct("match", 3, 0, false, nil)

	notMatch := ts.createProduct("not_match", 4, 0, false, nil)
	notMatch.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	var isEqualOperator badorm.Operator[sql.NullFloat64]

	switch getDBDialector() {
	case configuration.MySQL:
		isEqualOperator = mysql.IsEqual(sql.NullFloat64{Valid: false})
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		isEqualOperator = badorm.IsNotDistinct(sql.NullFloat64{Valid: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			isEqualOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestIsNull() {
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

func (ts *OperatorIntTestSuite) TestIsNullNotPointers() {
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

func (ts *OperatorIntTestSuite) TestIsNotNull() {
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

func (ts *OperatorIntTestSuite) TestIsNotNullNotPointers() {
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

func (ts *OperatorIntTestSuite) TestIsTrue() {
	match := ts.createProduct("match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	var isTrueOperator badorm.Operator[bool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isTrueOperator = badorm.IsTrue[bool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsTrue
		isTrueOperator = badorm.Eq[bool](true)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			isTrueOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	var isFalseOperator badorm.Operator[bool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isFalseOperator = badorm.IsFalse[bool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsFalse
		isFalseOperator = badorm.Eq[bool](false)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			isFalseOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

//nolint:dupl // not really duplicated
func (ts *OperatorIntTestSuite) TestIsNotTrue() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err := ts.db.Save(match2).Error
	ts.Nil(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: true}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	var isNotTrueOperator badorm.Operator[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isNotTrueOperator = badorm.IsNotTrue[sql.NullBool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsNotTrue
		isNotTrueOperator = badorm.IsDistinct[sql.NullBool](sql.NullBool{Valid: true, Bool: true})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isNotTrueOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

//nolint:dupl // not really duplicated
func (ts *OperatorIntTestSuite) TestIsNotFalse() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match2).Error
	ts.Nil(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	var isNotFalseOperator badorm.Operator[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		isNotFalseOperator = badorm.IsNotFalse[sql.NullBool]()
	case configuration.SQLServer:
		// sqlserver doesn't support IsNotFalse
		isNotFalseOperator = badorm.IsDistinct[sql.NullBool](sql.NullBool{Valid: true, Bool: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isNotFalseOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorIntTestSuite) TestIsUnknown() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch1 := ts.createProduct("match", 0, 0, false, nil)
	notMatch1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(notMatch1).Error
	ts.Nil(err)

	notMatch2 := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch2).Error
	ts.Nil(err)

	var isUnknownOperator badorm.Operator[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL:
		isUnknownOperator = badorm.IsUnknown[sql.NullBool]()
	case configuration.SQLServer, configuration.SQLite:
		// sqlserver doesn't support IsUnknown
		isUnknownOperator = badorm.IsNotDistinct[sql.NullBool](sql.NullBool{Valid: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isUnknownOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestIsNotUnknown() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match1).Error
	ts.Nil(err)

	match2 := ts.createProduct("", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(match2).Error
	ts.Nil(err)

	ts.createProduct("", 0, 0, false, nil)

	var isNotUnknownOperator badorm.Operator[sql.NullBool]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL:
		isNotUnknownOperator = badorm.IsNotUnknown[sql.NullBool]()
	case configuration.SQLServer, configuration.SQLite:
		// sqlserver doesn't support IsNotUnknown
		isNotUnknownOperator = badorm.IsDistinct[sql.NullBool](sql.NullBool{Valid: false})
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			isNotUnknownOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorIntTestSuite) TestArrayIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.ArrayIn("s1", "s2", "s3"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorIntTestSuite) TestArrayNotIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			badorm.ArrayNotIn("ns1", "ns2"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorIntTestSuite) TestLike() {
	match1 := ts.createProduct("basd", 0, 0, false, nil)
	match2 := ts.createProduct("cape", 0, 0, false, nil)

	ts.createProduct("bbsd", 0, 0, false, nil)
	ts.createProduct("bbasd", 0, 0, false, nil)

	var likeOperator badorm.Operator[string]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		likeOperator = badorm.Like[string]("_a%")
	case configuration.SQLServer:
		likeOperator = badorm.Like[string]("[bc]a[^a]%")
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			likeOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorIntTestSuite) TestLikeEscape() {
	match1 := ts.createProduct("ba_sd", 0, 0, false, nil)
	match2 := ts.createProduct("ca_pe", 0, 0, false, nil)

	ts.createProduct("bb_sd", 0, 0, false, nil)
	ts.createProduct("bba_sd", 0, 0, false, nil)

	var likeOperator badorm.Operator[string]

	switch getDBDialector() {
	case configuration.MySQL, configuration.PostgreSQL, configuration.SQLite:
		likeOperator = badorm.Like[string]("_a!_%").Escape('!')
	case configuration.SQLServer:
		likeOperator = badorm.Like[string]("[bc]a!_[^a]%").Escape('!')
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			likeOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorIntTestSuite) TestLikeOnNumeric() {
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

func (ts *OperatorIntTestSuite) TestILike() {
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

func (ts *OperatorIntTestSuite) TestSimilarTo() {
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

func (ts *OperatorIntTestSuite) TestPosixRegexCaseSensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)
	ts.createProduct("AB", 0, 0, false, nil)

	var posixRegexOperator badorm.Operator[string]

	switch getDBDialector() {
	case configuration.SQLServer, configuration.MySQL:
		log.Println("PosixRegex not compatible")
	case configuration.PostgreSQL:
		posixRegexOperator = psql.POSIXMatch[string]("^a(b|x)")
	case configuration.SQLite:
		posixRegexOperator = sqlite.Glob[string]("a[bx]")
	}

	if posixRegexOperator != nil {
		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				posixRegexOperator,
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorIntTestSuite) TestPosixRegexCaseInsensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)
	match3 := ts.createProduct("AB", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)

	var posixRegexOperator badorm.Operator[string]

	switch getDBDialector() {
	case configuration.SQLServer, configuration.SQLite:
		log.Println("PosixRegex Case Insensitive not compatible")
	case configuration.MySQL:
		posixRegexOperator = mysql.RegexP[string]("^a(b|x)")
	case configuration.PostgreSQL:
		posixRegexOperator = psql.POSIXIMatch[string]("^a(b|x)")
	}

	if posixRegexOperator != nil {
		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				posixRegexOperator,
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
}

func (ts *OperatorIntTestSuite) TestPosixRegexNotPosix() {
	var posixRegexOperator badorm.Operator[string]

	switch getDBDialector() {
	case configuration.SQLServer:
		log.Println("PosixRegex not compatible")
	case configuration.MySQL:
		posixRegexOperator = mysql.RegexP[string]("^a(b|x")
	case configuration.PostgreSQL:
		posixRegexOperator = psql.POSIXMatch[string]("^a(b|x")
	case configuration.SQLite:
		posixRegexOperator = sqlite.Glob[string]("^a(b|x")
	}

	if posixRegexOperator != nil {
		_, err := ts.crudProductService.GetEntities(
			conditions.ProductString(
				posixRegexOperator,
			),
		)
		ts.ErrorContains(err, "error parsing regexp")
	}
}

func (ts *OperatorIntTestSuite) TestDynamicOperatorForBasicType() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			dynamic.Eq(conditions.ProductIntPointerField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorIntTestSuite) TestDynamicOperatorForCustomType() {
	match := ts.createProduct("salut,hola", 1, 0.0, false, nil)
	match.MultiString = models.MultiString{"salut", "hola"}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("salut,hola", 1, 0.0, false, nil)
	ts.createProduct("hola", 1, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductMultiString(
			dynamic.Eq(conditions.ProductMultiStringField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestDynamicOperatorForBadORMModelAttribute() {
	match := ts.createProduct("", 1, 0.0, false, nil)

	var isNotDistinctOperator badorm.Operator[gorm.DeletedAt]

	switch getDBDialector() {
	case configuration.MySQL:
		isNotDistinctOperator = mysqldynamic.IsEqual(conditions.ProductDeletedAtField)
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		isNotDistinctOperator = dynamic.IsNotDistinct(conditions.ProductDeletedAtField)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductDeletedAt(isNotDistinctOperator),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestMultitypeOperatorWithFieldOfAnotherTypeReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			multitype.Eq[int](conditions.ProductStringField),
		),
	)
	ts.ErrorIs(err, multitype.ErrFieldTypeDoesNotMatch)
	ts.ErrorContains(err, "field type: string, attribute type: int; operator: Eq; model: models.Product, field: Int")
}

func (ts *OperatorIntTestSuite) TestMultitypeOperatorForNullableTypeCanBeComparedWithNotNullType() {
	match := ts.createProduct("", 1, 1.0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.0}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("", 1, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullFloat(
			multitype.Eq[sql.NullFloat64](conditions.ProductFloatField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestMultitypeOperatorForNotNullTypeCanBeComparedWithNullableType() {
	match := ts.createProduct("", 1, 1.0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.0}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("", 1, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductFloat(
			multitype.Eq[float64](conditions.ProductNullFloatField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestMultitypeOperatorForBadORMModelAttribute() {
	match := ts.createProduct("", 1, 0.0, false, nil)

	var isDistinctCondition badorm.Condition[models.Product]

	switch getDBDialector() {
	case configuration.MySQL:
		isDistinctCondition = badorm.Not[models.Product](
			conditions.ProductDeletedAt(
				mysqlmultitype.IsEqual[gorm.DeletedAt](conditions.ProductCreatedAtField),
			),
		)
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		isDistinctCondition = conditions.ProductDeletedAt(
			multitype.IsDistinct[gorm.DeletedAt](conditions.ProductCreatedAtField),
		)
	}

	entities, err := ts.crudProductService.GetEntities(isDistinctCondition)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestMultitypeMultivalueOperatorWithValueOfAnotherTypeReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			multitype.Between[int, int]("hola", 1),
		),
	)
	ts.ErrorIs(err, multitype.ErrParamNotValueOrField)
	ts.ErrorContains(err, "parameter type: string, attribute type: int; operator: Between; model: models.Product, field: Int")
}

func (ts *OperatorIntTestSuite) TestMultitypeMultivalueOperatorWithFieldOfAnotherTypeReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			multitype.Between[int, int](1, conditions.ProductCreatedAtField),
		),
	)
	ts.ErrorIs(err, multitype.ErrParamNotValueOrField)
	ts.ErrorContains(err, "parameter type: badorm.FieldIdentifier[time.Time], attribute type: int; operator: Between; model: models.Product, field: Int")
}

func (ts *OperatorIntTestSuite) TestMultitypeMultivalueOperatorWithFieldOfNotRelatedTypeReturnsError() {
	_, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			multitype.Between[int, time.Time](1, conditions.ProductCreatedAtField),
		),
	)
	ts.ErrorIs(err, multitype.ErrFieldTypeDoesNotMatch)
	ts.ErrorContains(err, "field type: time.Time, attribute type: int; operator: Between; model: models.Product, field: Int")
}

func (ts *OperatorIntTestSuite) TestMultitypeMultivalueOperatorWithAFieldAndAValue() {
	match := ts.createProduct("", 1, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			multitype.Between[int, int](1, conditions.ProductIntField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestMultitypeMultivalueOperatorWithAFieldRelatedAndAValue() {
	match := ts.createProduct("", 1, 1.0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 2.0}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	notMatch1 := ts.createProduct("", 0, 0.0, false, nil)
	notMatch1.NullFloat = sql.NullFloat64{Valid: true, Float64: 2.0}
	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	notMatch2 := ts.createProduct("", 0, 5.0, false, nil)
	notMatch2.NullFloat = sql.NullFloat64{Valid: true, Float64: 2.0}
	err = ts.db.Save(notMatch2).Error
	ts.Nil(err)

	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductFloat(
			multitype.Between[float64, sql.NullFloat64](1.0, conditions.ProductNullFloatField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorIntTestSuite) TestUnsafeOperatorInCaseTypesNotMatch() {
	switch getDBDialector() {
	case configuration.MySQL:
		// in mysql comparisons between types are allowed
		match1 := ts.createProduct("0", 1, 0, false, nil)
		match2 := ts.createProduct("0.0", 2, 0.0, false, nil)
		ts.createProduct("0.0", 2, 1.0, false, nil)

		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductFloat(
				unsafe.Eq[float64]("string"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case configuration.PostgreSQL, configuration.SQLServer, configuration.SQLite:
		// on postgresql returns an error
		_, err := ts.crudProductService.GetEntities(
			conditions.ProductFloat(
				unsafe.Eq[float64]("string"),
			),
		)
		ts.ErrorContains(err, "string")
	}
}

func (ts *OperatorIntTestSuite) TestUnsafeOperatorCanCompareFieldsThatMapToTheSameType() {
	match := ts.createProduct("hola,chau", 1, 1.0, false, nil)
	match.MultiString = models.MultiString{"hola", "chau"}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	notMatch := ts.createProduct("chau", 0, 0.0, false, nil)
	notMatch.MultiString = models.MultiString{"hola", "chau"}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			unsafe.Eq[string](conditions.ProductMultiStringField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}
