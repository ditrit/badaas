package integrationtests

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/google/uuid"
)

// TODO gorm
// column
// embedded
// embeddedPrefix
// todas las clases tienen sus conditions, repository and service menos las que esten embeded en otras
// aunque tambien podria ser que:
// haya embeded que tambien tengan su propia tabla
// haya cosas en el modelo que no van a terminar en una tabla, que son solo clases para llamar a metodos y eso

// podria meter alguna anotacion para que esa si vaya a modelos?
// directamente las clases que tienen un base model metido adentro son las que quiero

type Company struct {
	badorm.UUIDModel

	Name    string
	Sellers []Seller // Company HasMany Sellers (Company 1 -> 0..* Seller)
}

type Product struct {
	badorm.UUIDModel

	String string
	Int    int
	Float  float64
	Bool   bool
}

type Seller struct {
	badorm.UUIDModel

	Name      string
	CompanyID *uuid.UUID // Company HasMany Sellers (Company 1 -> 0..* Seller)
}

func SellerCompanyCondition(conditions ...badorm.Condition[Company]) badorm.Condition[Seller] {
	return badorm.JoinCondition[Seller, Company]{
		Field:      "company",
		Conditions: conditions,
	}
}

type Sale struct {
	badorm.UUIDModel

	Code        int
	Description string

	// Sale belongsTo Product (Sale 0..* -> 1 Product)
	Product   Product
	ProductID uuid.UUID

	// Sale HasOne Seller (Sale 0..* -> 0..1 Seller)
	Seller   *Seller
	SellerID *uuid.UUID
}

func SaleProductIDCondition(id uuid.UUID) badorm.Condition[Sale] {
	return badorm.WhereCondition[Sale]{
		Field: "product_id",
		Value: id.String(),
	}
}

func SaleProductCondition(conditions ...badorm.Condition[Product]) badorm.Condition[Sale] {
	return badorm.JoinCondition[Sale, Product]{
		Field:      "product",
		Conditions: conditions,
	}
}

func SaleSellerCondition(conditions ...badorm.Condition[Seller]) badorm.Condition[Sale] {
	return badorm.JoinCondition[Sale, Seller]{
		Field:      "seller",
		Conditions: conditions,
	}
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
	badorm.UUIDModel

	Name    string
	Capital City // Country HasOne City (Country 1 -> 1 City)
}

type City struct {
	badorm.UUIDModel

	Name      string
	CountryID uuid.UUID // Country HasOne City (Country 1 -> 1 City)
}

func (m Country) Equal(other Country) bool {
	return m.Name == other.Name
}

func (m City) Equal(other City) bool {
	return m.Name == other.Name
}

func CityCountryCondition(conditions ...badorm.Condition[Country]) badorm.Condition[City] {
	return badorm.JoinCondition[City, Country]{
		Field:      "country",
		Conditions: conditions,
	}
}

func CountryCapitalCondition(conditions ...badorm.Condition[City]) badorm.Condition[Country] {
	return badorm.JoinCondition[Country, City]{
		Field:      "capital",
		Conditions: conditions,
	}
}

type Employee struct {
	badorm.UUIDModel

	Name   string
	Boss   *Employee // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *uuid.UUID
}

func EmployeeBossCondition(conditions ...badorm.Condition[Employee]) badorm.Condition[Employee] {
	return badorm.JoinCondition[Employee, Employee]{
		Field:      "boss",
		Conditions: conditions,
	}
}

func (m Employee) Equal(other Employee) bool {
	return m.Name == other.Name
}

type Person struct {
	badorm.UUIDModel

	Name string
}

func (m Person) TableName() string {
	return "persons_and_more_name"
}

type Bicycle struct {
	badorm.UUIDModel

	Name string
	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner   Person
	OwnerID uuid.UUID
}

func BicycleOwnerCondition(conditions ...badorm.Condition[Person]) badorm.Condition[Bicycle] {
	return badorm.JoinCondition[Bicycle, Person]{
		Field:      "owner",
		Conditions: conditions,
	}
}

func (m Bicycle) Equal(other Bicycle) bool {
	return m.Name == other.Name
}