package integrationtests

import (
	"log"

	"github.com/ditrit/badaas/persistence/models"
	"gorm.io/gorm"
)

var ListOfTables = []any{
	models.Session{},
	models.User{},
	models.Value{},
	models.Attribute{},
	models.Entity{},
	models.EntityType{},
	Sale{},
	Product{},
	Seller{},
	Company{},
}

func CleanDB(db *gorm.DB) {
	// clean database to ensure independency between tests
	for _, table := range ListOfTables {
		err := db.Unscoped().Where("1 = 1").Delete(table).Error
		if err != nil {
			log.Fatalln("could not clean database: ", err)
		}
	}
}
