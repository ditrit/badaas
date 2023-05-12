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
	Sellers []Seller // Company HasMany Sellers (Company 1 -> 0..* Seller)
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
	CompanyID *uuid.UUID // Company HasMany Sellers (Company 1 -> 0..* Seller)
}

type Sale struct {
	models.BaseModel

	// Sale belongsTo Product (Sale 0..* -> 1 Product)
	Product   Product
	ProductID uuid.UUID

	// Sale HasOne Seller (Sale 0..* -> 0..1 Seller)
	Seller   *Seller
	SellerID *uuid.UUID
}

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
}

func (m Sale) Equal(other Sale) bool {
	return m.ID == other.ID
}

func (m Seller) Equal(other Seller) bool {
	return m.Name == other.Name
}

type Country struct {
	models.BaseModel

	Name    string
	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	models.BaseModel

	Name      string
	CountryID uuid.UUID // Country HasOne City (Country 1 -> 1 City)
}

func (m Country) Equal(other Country) bool {
	return m.Name == other.Name
}

func (m City) Equal(other City) bool {
	return m.Name == other.Name
}

type Employee struct {
	models.BaseModel

	Name   string
	Boss   *Employee // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *uuid.UUID
}

func (m Employee) Equal(other Employee) bool {
	return m.Name == other.Name
}

type Person struct {
	models.BaseModel

	Name string
}

func (m Person) TableName() string {
	return "persons_and_more_name"
}

type Bicycle struct {
	models.BaseModel

	Name string
	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner   Person
	OwnerID uuid.UUID
}

func (m Bicycle) Equal(other Bicycle) bool {
	return m.Name == other.Name
}

type CRUDServiceIntTestSuite struct {
	suite.Suite
	logger              *zap.Logger
	db                  *gorm.DB
	crudProductService  badorm.CRUDService[Product, uuid.UUID]
	crudSaleService     badorm.CRUDService[Sale, uuid.UUID]
	crudSellerService   badorm.CRUDService[Seller, uuid.UUID]
	crudCountryService  badorm.CRUDService[Country, uuid.UUID]
	crudCityService     badorm.CRUDService[City, uuid.UUID]
	crudEmployeeService badorm.CRUDService[Employee, uuid.UUID]
	crudBicycleService  badorm.CRUDService[Bicycle, uuid.UUID]
}

func NewCRUDServiceIntTestSuite(
	logger *zap.Logger,
	db *gorm.DB,
	crudProductService badorm.CRUDService[Product, uuid.UUID],
	crudSaleService badorm.CRUDService[Sale, uuid.UUID],
	crudSellerService badorm.CRUDService[Seller, uuid.UUID],
	crudCountryService badorm.CRUDService[Country, uuid.UUID],
	crudCityService badorm.CRUDService[City, uuid.UUID],
	crudEmployeeService badorm.CRUDService[Employee, uuid.UUID],
	crudBicycleService badorm.CRUDService[Bicycle, uuid.UUID],
) *CRUDServiceIntTestSuite {
	return &CRUDServiceIntTestSuite{
		logger:              logger,
		db:                  db,
		crudProductService:  crudProductService,
		crudSaleService:     crudSaleService,
		crudSellerService:   crudSellerService,
		crudCountryService:  crudCountryService,
		crudCityService:     crudCityService,
		crudEmployeeService: crudEmployeeService,
		crudBicycleService:  crudBicycleService,
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

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct(map[string]any{
		"int": 1,
	})
	product2 := ts.createProduct(map[string]any{
		"int": 2,
	})

	match := ts.createSale(product1, nil)
	ts.createSale(product2, nil)

	params := map[string]any{
		"Product": map[string]any{
			"int": 1.0,
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsHasOneOptional() {
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
		"Seller": map[string]any{
			"name": "franco",
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsHasOneSelfReferential() {
	boss1 := &Employee{
		Name: "Xavier",
	}
	boss2 := &Employee{
		Name: "Vincent",
	}

	match := ts.createEmployee("franco", boss1)
	ts.createEmployee("pierre", boss2)

	params := map[string]any{
		"Boss": map[string]any{
			"name": "Xavier",
		},
	}
	entities, err := ts.crudEmployeeService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Employee{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOneToOne() {
	capital1 := City{
		Name: "Buenos Aires",
	}
	capital2 := City{
		Name: "Paris",
	}

	ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	params := map[string]any{
		"Country": map[string]any{
			"name": "Argentina",
		},
	}
	entities, err := ts.crudCityService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*City{&capital1}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOneToOneReversed() {
	capital1 := City{
		Name: "Buenos Aires",
	}
	capital2 := City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	params := map[string]any{
		"Capital": map[string]any{
			"name": "Buenos Aires",
		},
	}
	entities, err := ts.crudCountryService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Country{country1}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsReturnsErrorIfNoRelation() {
	params := map[string]any{
		"NotExists": map[string]any{
			"int": 1.0,
		},
	}
	_, err := ts.crudSaleService.GetEntities(params)
	ts.ErrorContains(err, "Sale has not attribute named NotExists or NotExistsID")
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsWithEntityThatDefinesTableName() {
	person1 := Person{
		Name: "franco",
	}
	person2 := Person{
		Name: "xavier",
	}

	match := ts.createBicycle("BMX", person1)
	ts.createBicycle("Shimano", person2)

	params := map[string]any{
		"Owner": map[string]any{
			"name": "franco",
		},
	}
	entities, err := ts.crudBicycleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Bicycle{match}, entities)
}

func (ts *CRUDServiceIntTestSuite) TestGetEntitiesWithConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	params := map[string]any{
		"Company": map[string]any{
			"name": "ditrit",
		},
	}
	entities, err := ts.crudSellerService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Seller{match}, entities)
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
		"Product": map[string]any{
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
		"Product": map[string]any{
			"int": 1.0,
		},
		"Seller": map[string]any{
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
		"Seller": map[string]any{
			"name": "franco",
			"Company": map[string]any{
				"name": "ditrit",
			},
		},
	}
	entities, err := ts.crudSaleService.GetEntities(params)
	ts.Nil(err)

	EqualList(&ts.Suite, []*Sale{match}, entities)
}

// TODO test de que hay que poner el mismo nombre que el atributo y no como en la base de datos, aunque es lo mismo que el de product_id
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
		Product: *product,
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

func (ts *CRUDServiceIntTestSuite) createCountry(name string, capital City) *Country {
	entity := &Country{
		Name:    name,
		Capital: capital,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createEmployee(name string, boss *Employee) *Employee {
	entity := &Employee{
		Name: name,
		Boss: boss,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *CRUDServiceIntTestSuite) createBicycle(name string, owner Person) *Bicycle {
	entity := &Bicycle{
		Name:  name,
		Owner: owner,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}
