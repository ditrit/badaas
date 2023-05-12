package integrationtests

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Company struct {
	models.BaseModel

	Name    string
	Sellers []Seller
}

type Product struct {
	models.BaseModel

	String string
	Int    int
	Float  float64
	Bool   bool
}

type Seller struct {
	models.BaseModel

	Name      string
	CompanyID *uuid.UUID
}

type Sale struct {
	models.BaseModel

	// belongsTo Product
	Product   *Product
	ProductID uuid.UUID

	// belongsTo Seller
	Seller   *Seller
	SellerID uuid.UUID
}

func (Product) TableName() string {
	return "products"
}

func (Sale) TableName() string {
	return "sales"
}

func (Company) TableName() string {
	return "companies"
}

func (Seller) TableName() string {
	return "sellers"
}

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
}

func (m Sale) Equal(other Sale) bool {
	return m.ID == other.ID
}

type CRUDServiceIntTestSuite struct {
	suite.Suite
	logger             *zap.Logger
	db                 *gorm.DB
	crudProductService badorm.CRUDService[Product, uuid.UUID]
	crudSaleService    badorm.CRUDService[Sale, uuid.UUID]
}

func NewCRUDServiceIntTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
	crudProductService badorm.CRUDService[Product, uuid.UUID],
	crudSaleService badorm.CRUDService[Sale, uuid.UUID],
) *CRUDServiceIntTestSuite {
	return &CRUDServiceIntTestSuite{
		logger:             logger,
		db:                 db,
		crudProductService: crudProductService,
		crudSaleService:    crudSaleService,
	}
}

func (ts *CRUDServiceIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

// ------------------------- GetEntities --------------------------------

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.GetEntities(map[string]any{})
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct(map[string]any{})

	entities, err := ts.crudProductService.GetEntities(map[string]any{})
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithoutConditionsReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct(map[string]any{})
	match2 := ts.createProduct(map[string]any{})
	match3 := ts.createProduct(map[string]any{})

	entities, err := ts.crudProductService.GetEntities(map[string]any{})
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match1, match2, match3}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNotEntitiesCreated() {
	params := map[string]any{
		"string": "not_created",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{}, entities)
}

// TODO mirar esto, el string del where se pone sin tabla, puede generar problemas cuando haya join
func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct(map[string]any{
		"string": "something_else",
	})

	params := map[string]any{
		"string": "not_match",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createProduct(map[string]any{
		"string": "match",
	})
	ts.createProduct(map[string]any{
		"string": "something_else",
	})

	params := map[string]any{
		"string": "match",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProduct(map[string]any{
		"string": "match",
	})
	match2 := ts.createProduct(map[string]any{
		"string": "match",
	})
	ts.createProduct(map[string]any{
		"string": "something_else",
	})

	params := map[string]any{
		"string": "match",
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match1, match2}, entities)
}

// TODO ver caso en el que la columna no existe
// func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatDoesNotExistReturnsEmpty() {
// 	ts.createProduct(map[string]any{
// 		"string": "match",
// 	})
// 	ts.createProduct(map[string]any{
// 		"string": "match",
// 	})
// 	ts.createProduct(map[string]any{
// 		"string": "match",
// 	})

// 	params := map[string]any{
// 		"not_exists": "not_exists",
// 	}
// 	entities, err := ts.crudProductService.GetEntities(params)
// 	ts.Nil(err)

// 	EqualList(&ts.Suite, []*Product{}, entities)
// }

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIntType() {
	match := ts.createProduct(map[string]any{
		"string": "match",
		"int":    1,
	})
	ts.createProduct(map[string]any{
		"string": "not_match",
		"int":    2,
	})

	params := map[string]any{
		"int": 1.0,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

// TODO ver este caso
// func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfIncorrectTypeReturnsEmptyList() {
// 	ts.createProduct(map[string]any{
// 		"string": "not_match",
// 		"int":    1,
// 	})

// 	params := map[string]any{
// 		"int": "not_an_int",
// 	}
// 	entities, err := ts.crudProductService.GetEntities(params)
// 	ts.Nil(err)
// 	ts.Len(entities, 0)
// }

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfFloatType() {
	match := ts.createProduct(map[string]any{
		"string": "match",
		"float":  1.1,
	})
	ts.createProduct(map[string]any{
		"string": "not_match",
		"float":  2.0,
	})

	params := map[string]any{
		"float": 1.1,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfBoolType() {
	match := ts.createProduct(map[string]any{
		"string": "match",
		"bool":   true,
	})
	ts.createProduct(map[string]any{
		"string": "not_match",
		"bool":   false,
	})

	params := map[string]any{
		"bool": true,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match}, entities)
}

// TODO testear la creacion directo con entidades
func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionOfRelationType() {
	product1 := ts.createProduct(map[string]any{})
	product2 := ts.createProduct(map[string]any{})

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(product1, seller1)
	ts.createSale(product2, seller2)

	params := map[string]any{
		// TODO ver esto, es un poco confuso que lo muestra como ProductID pero para queries es product_id
		"product_id": product1.ID.String(),
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

// TODO ver esto
// func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionFilterByNull() {
// 	match := ts.createProduct(map[string]any{})
// 	ts.createProduct(map[string]any{
// 		"string": "something",
// 	})

// 	params := map[string]any{
// 		"string": nil,
// 	}
// 	entities, err := ts.crudProductService.GetEntities(params)
// 	ts.Nil(err)

// 	EqualList(&ts.Suite, []*Product{match}, entities)
// }

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createProduct(map[string]any{
		"string": "match",
		"int":    1,
		"bool":   true,
	})
	match2 := ts.createProduct(map[string]any{
		"string": "match",
		"int":    1,
		"bool":   true,
	})

	ts.createProduct(map[string]any{
		"string": "not_match",
		"int":    1,
		"bool":   true,
	})
	ts.createProduct(map[string]any{
		"string": "match",
		"int":    2,
		"bool":   true,
	})

	params := map[string]any{
		"string": "match",
		"int":    1.0,
		"bool":   true,
	}
	entities, err := ts.crudProductService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Product{match1, match2}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoins() {
	product1 := ts.createProduct(map[string]any{
		"int": 1,
	})
	product2 := ts.createProduct(map[string]any{
		"int": 2,
	})

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(product1, seller1)
	ts.createSale(product2, seller2)

	params := map[string]any{
		"product": map[string]any{
			"int": 1.0,
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct(map[string]any{
		"int":    1,
		"string": "match",
	})
	product2 := ts.createProduct(map[string]any{
		"int":    2,
		"string": "match",
	})

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(product1, seller1)
	ts.createSale(product2, seller2)

	params := map[string]any{
		"product": map[string]any{
			"int":    1.0,
			"string": "match",
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

// // TODO aca deberia agregar casos en lo que 1 matchea pero el otro no
func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsDifferentEntities() {
	product1 := ts.createProduct(map[string]any{
		"int": 1,
	})
	product2 := ts.createProduct(map[string]any{
		"int": 2,
	})

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(product1, seller1)
	ts.createSale(product2, seller2)

	params := map[string]any{
		"product": map[string]any{
			"int": 1.0,
		},
		"seller": map[string]any{
			"name": "franco",
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsMultipleTimes() {
	product1 := ts.createProduct(map[string]any{})
	product2 := ts.createProduct(map[string]any{})

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(product1, seller1)
	ts.createSale(product2, seller2)

	params := map[string]any{
		// TODO meter algo aca tambien, que sea join y mio
		"seller": map[string]any{
			"name": "franco",
			"company": map[string]any{
				"name": "ditrit",
			},
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

// TODO falta test de que joinea dos veces con la misma entidad
// TODO faltan test para otros tipos de relaciones

// ------------------------- utils -------------------------

func (ts *CRUDServiceIntTestSuite) createProduct(values map[string]any) *Product {
	entity, err := ts.crudProductService.CreateEntity(values)
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createSale(product *Product, seller *Seller) *Sale {
	entity := &Sale{
		Product: product,
		Seller:  seller,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createSeller(name string, company *Company) *Seller {
	var companyID *uuid.UUID
	if company != nil {
		companyID = &company.ID
	}
	entity := &Seller{
		Name:      name,
		CompanyID: companyID,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createCompany(name string) *Company {
	entity := &Company{
		Name: name,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}
