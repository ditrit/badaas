//go:build mysql
// +build mysql

package models

import "github.com/ditrit/badaas/badorm"

type Employee struct {
	badorm.UUIDModel

	Name string
	// mysql needs OnDelete to work with self-referential fk
	Boss   *Employee `gorm:"constraint:OnDelete:SET NULL"` // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *badorm.UUID
}

func (m Employee) Equal(other Employee) bool {
	return m.Name == other.Name
}
