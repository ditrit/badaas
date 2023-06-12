package conditions

import (
	badorm "github.com/ditrit/badaas/badorm"
	models "github.com/ditrit/badaas/testintegration/models"
)

func ProductId(expressions ...badorm.Expression) badorm.WhereCondition[models.Product] {
	return badorm.WhereCondition[models.Product]{
		Field:       "ID",
		Expressions: expressions,
	}
}
