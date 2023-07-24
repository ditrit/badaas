package overridereferencesinverse

import "github.com/ditrit/badaas/badorm"

type Computer struct {
	badorm.UUIDModel
	Name      string
	Processor Processor `gorm:"foreignKey:ComputerName;references:Name"`
}

type Processor struct {
	badorm.UUIDModel
	ComputerName string
}
