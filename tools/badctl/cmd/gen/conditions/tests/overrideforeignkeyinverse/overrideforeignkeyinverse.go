package overrideforeignkeyinverse

import (
	"github.com/ditrit/badaas/badorm"
)

type User struct {
	badorm.UUIDModel
	CreditCard CreditCard `gorm:"foreignKey:UserReference"`
}

type CreditCard struct {
	badorm.UUIDModel
	UserReference badorm.UUID
}
