package persistence

import (
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"go.uber.org/fx"
)

// PersistanceModule for fx
//
// Provides:
//
// - The database connection
var PersistanceModule = fx.Module(
	"persistence",
	// Database connection
	fx.Provide(
		fx.Annotate(
			gormdatabase.SetupDatabaseConnection,
			fx.ParamTags(`group:"modelTables"`),
		),
	),
)

type AddModelResult struct {
	fx.Out

	Model any `group:"modelTables"`
}

func AddModel[T any]() AddModelResult {
	return AddModelResult{
		Model: *new(T),
	}
}
