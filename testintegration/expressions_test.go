package testintegration

import (
	"database/sql"
	"log"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/mysql"
	"github.com/ditrit/badaas/badorm/psql"
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
	case configuration.PostgreSQL:
		entities, err := ts.crudProductService.GetEntities(
			conditions.ProductInt(
				psql.IsDistinct(2),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case configuration.MySQL, configuration.SQLServer, configuration.SQLite:
		// TODO
		log.Println("TODO")
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
	case configuration.PostgreSQL:
		isNotEqualExpression = psql.IsNotDistinct(3)
	case configuration.SQLServer, configuration.SQLite:
		// TODO esto no va a andar en todos
		isNotEqualExpression = psql.IsNotDistinct(3)
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductInt(
			isNotEqualExpression,
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
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNull[*int](),
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
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNotNull[*int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsTrue() {
	match := ts.createProduct("match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsTrue[bool](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *ExpressionIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsFalse[bool](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNotTrue[sql.NullBool](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNotFalse[sql.NullBool](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsUnknown[sql.NullBool](),
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

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductNullBool(
			// TODO esto no queda muy lindo que hay que ponerlo asi
			badorm.IsNotUnknown[sql.NullBool](),
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
	case configuration.SQLServer, configuration.SQLite:
		// TODO esto no va a andar en todos
		arrayInExpression = psql.ArrayIn("s1", "s2", "s3")
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
	case configuration.SQLServer, configuration.SQLite:
		// TODO esto no va a andar en todos
		arrayNotInExpression = psql.ArrayNotIn("ns1", "ns2")
	}

	entities, err := ts.crudProductService.GetEntities(
		conditions.ProductString(
			arrayNotInExpression,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *ExpressionIntTestSuite) TestMultipleExpressions() {
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

func (ts *ExpressionIntTestSuite) TestMultipleConditionsDifferentExpressions() {
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
