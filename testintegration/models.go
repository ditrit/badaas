package testintegration

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/ditrit/badaas/badorm"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// TODO testear tambien en otras bases de datos

// TODO todas las clases badorm tienen sus conditions, repository and service

type Company struct {
	badorm.UUIDModel

	Name    string
	Sellers []Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}
	*s = strings.Split(str, ",")
	return nil
}

func (s MultiString) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

func (MultiString) GormDataType() string {
	return "text"
}

type ToBeEmbedded struct {
	EmbeddedInt int
}

type ToBeGormEmbedded struct {
	Int int
}

type Product struct {
	badorm.UUIDModel

	String      string `gorm:"column:string_something_else"`
	Int         int
	IntPointer  *int
	Float       float64
	Bool        bool
	ByteArray   []byte
	MultiString MultiString
	StringArray pq.StringArray `gorm:"type:text[]"`
	ToBeEmbedded
	GormEmbedded ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
}

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
}

type Seller struct {
	badorm.UUIDModel

	Name      string
	CompanyID *uuid.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
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

type Employee struct {
	badorm.UUIDModel

	Name   string
	Boss   *Employee // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *uuid.UUID
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

func (m Bicycle) Equal(other Bicycle) bool {
	return m.Name == other.Name
}
