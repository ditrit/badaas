package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ditrit/badaas/badorm"
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
	switch typedSrc := src.(type) {
	case string:
		*s = strings.Split(typedSrc, ",")
		return nil
	case []byte:
		str := string(typedSrc)
		*s = strings.Split(str, ",")
		return nil
	default:
		return fmt.Errorf("failed to scan multistring field - source is not a string, is %T", src)
	}
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

func (MultiString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlserver":
		return "varchar(255)"
	default:
		return "text"
	}
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
	ToBeEmbedded
	GormEmbedded ToBeGormEmbedded `gorm:"embedded;embeddedPrefix:gorm_embedded_"`
}

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
}

type Seller struct {
	badorm.UUIDModel

	Name      string
	CompanyID *badorm.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type Sale struct {
	badorm.UUIDModel

	Code        int
	Description string

	// Sale belongsTo Product (Sale 0..* -> 1 Product)
	Product   Product
	ProductID badorm.UUID

	// Sale HasOne Seller (Sale 0..* -> 0..1 Seller)
	Seller   *Seller
	SellerID *badorm.UUID
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
	CountryID badorm.UUID // Country HasOne City (Country 1 -> 1 City)
}

func (m Country) Equal(other Country) bool {
	return m.Name == other.Name
}

func (m City) Equal(other City) bool {
	return m.Name == other.Name
}

type Person struct {
	badorm.UUIDModel

	Name string `gorm:"unique;type:VARCHAR(255)"`
}

func (m Person) TableName() string {
	return "persons_and_more_name"
}

type Bicycle struct {
	badorm.UUIDModel

	Name string
	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner     Person `gorm:"references:Name;foreignKey:OwnerName"`
	OwnerName string
}

func (m Bicycle) Equal(other Bicycle) bool {
	return m.Name == other.Name
}

type Brand struct {
	badorm.UIntModel

	Name string
}

func (m Brand) Equal(other Brand) bool {
	return m.Name == other.Name
}

type Phone struct {
	badorm.UIntModel

	Name string
	// Phone belongsTo Brand (Phone 0..* -> 1 Brand)
	Brand   Brand
	BrandID uint
}

func (m Phone) Equal(other Phone) bool {
	return m.Name == other.Name
}
