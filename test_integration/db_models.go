package integrationtests

import (
	"log"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

var ListOfTables = []any{
	Product{},
	Company{},
	Seller{},
	Sale{},
	Country{},
	City{},
	Employee{},
	Person{},
	Bicycle{},
}

func GetModels() badorm.GetModelsResult {
	return badorm.GetModelsResult{
		Models: ListOfTables,
	}
}

func CleanDB(db *gorm.DB) {
	CleanDBTables(db, append(
		pie.Reverse(ListOfTables),
		[]any{
			models.Value{},
			models.Attribute{},
			models.Entity{},
			models.EntityType{},
		}...,
	))
}

func CleanDBTables(db *gorm.DB, listOfTables []any) {
	// clean database to ensure independency between tests
	for _, table := range listOfTables {
		err := db.Unscoped().Where("1 = 1").Delete(table).Error
		if err != nil {
			log.Fatalln("could not clean database: ", err)
		}
	}
}
