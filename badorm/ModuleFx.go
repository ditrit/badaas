package badorm

import (
	"fmt"

	"github.com/ditrit/badaas/badorm/tabler"
	"go.uber.org/fx"
)

var BaDORMModule = fx.Module(
	"BaDORM",
	fx.Invoke(
		fx.Annotate(
			autoMigrate,
			fx.ParamTags(`group:"modelTables"`),
		),
	),
)

func GetCRUDServiceModule[T tabler.Tabler, ID BadaasID]() fx.Option {
	return fx.Module(
		fmt.Sprintf(
			"%TCRUDServiceModule",
			*new(T),
		),
		// models
		fx.Provide(AddModel[T]),
		// repository
		fx.Provide(NewCRUDRepository[T, ID]),
		// service
		fx.Provide(NewCRUDService[T, ID]),
	)
}

type AddModelResult struct {
	fx.Out

	Model any `group:"modelTables"`
}

func AddModel[T any]() AddModelResult {
	return AddModelResult{
		Model: *new(T),
	}
}
