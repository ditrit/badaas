package selfreferential

import "github.com/ditrit/badaas/badorm"

type Employee struct {
	badorm.UUIDModel

	Boss   *Employee `gorm:"constraint:OnDelete:SET NULL;"` // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *badorm.UUID
}
