package hasmanywithpointers

import "github.com/ditrit/badaas/badorm"

type CompanyWithPointers struct {
	badorm.UUIDModel

	Sellers *[]*SellerInPointers // CompanyWithPointers HasMany SellerInPointers
}

type SellerInPointers struct {
	badorm.UUIDModel

	Company   *CompanyWithPointers
	CompanyID *badorm.UUID // Company HasMany Seller
}
