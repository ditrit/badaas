package overrideforeignkey

import "github.com/ditrit/badaas/badorm"

type Person struct {
	badorm.UUIDModel
}

type Bicycle struct {
	badorm.UUIDModel

	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner            Person `gorm:"foreignKey:OwnerSomethingID"`
	OwnerSomethingID string
}
