package persistence

import (
	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"go.uber.org/fx"
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
