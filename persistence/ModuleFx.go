package persistence

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/gormdatabase"
)

// PersistanceModule for fx
//
// Provides:
//
// - The database connection
// - BaDORM
var PersistanceModule = fx.Module(
	"persistence",
	// Database connection
	fx.Provide(gormdatabase.SetupDatabaseConnection),
	// activate BaDORM
	badorm.BaDORMModule,
)
