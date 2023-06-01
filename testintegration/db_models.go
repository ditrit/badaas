package testintegration

import (
	"log"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/testintegration/models"
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

var ListOfTables = []any{
	models.Product{},
	models.Company{},
	models.Seller{},
	models.Sale{},
	models.Country{},
	models.City{},
	models.Employee{},
	models.Person{},
	models.Bicycle{},
}

func GetModels() badorm.GetModelsResult {
	return badorm.GetModelsResult{
		Models: ListOfTables,
	}
}

func CleanDB(db *gorm.DB) {
	CleanDBTables(db, pie.Reverse(ListOfTables))
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
