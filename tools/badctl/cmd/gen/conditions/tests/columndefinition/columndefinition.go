package columndefinition

import "github.com/ditrit/badaas/badorm"

type ColumnDefinition struct {
	badorm.UUIDModel

	String string `gorm:"column:string_something_else"`
}
