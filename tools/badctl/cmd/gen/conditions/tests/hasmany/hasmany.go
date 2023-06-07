package hasone

import "github.com/ditrit/badaas/badorm"

type Company struct {
	badorm.UUIDModel

	Sellers []Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type Seller struct {
	badorm.UUIDModel

	CompanyID *badorm.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}
