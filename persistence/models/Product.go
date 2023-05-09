package models

import "github.com/google/uuid"

// TODO esto no deberia estar aca
type Company struct {
	BaseModel

	Name    string
	Sellers []Seller
}

type Product struct {
	BaseModel

	String string
	Int    int
	Float  float64
	Bool   bool
}

type Seller struct {
	BaseModel

	Name      string
	CompanyID *uuid.UUID
}

type Sale struct {
	BaseModel

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

func (m Product) Equal(other Product) bool {
	return m.ID == other.ID
}

func (m Sale) Equal(other Sale) bool {
	return m.ID == other.ID
}
