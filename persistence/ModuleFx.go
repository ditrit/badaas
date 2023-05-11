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
	fx.Provide(gormdatabase.SetupDatabaseConnection),
)
