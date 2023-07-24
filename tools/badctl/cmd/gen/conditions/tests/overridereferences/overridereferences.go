package overridereferences

import "github.com/ditrit/badaas/badorm"

type Brand struct {
	badorm.UUIDModel

	Name string `gorm:"unique;type:VARCHAR(255)"`
}

type Phone struct {
	badorm.UUIDModel

	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Brand     Brand `gorm:"references:Name;foreignKey:BrandName"`
	BrandName string
}
