package persistence

import (
	"github.com/google/uuid"

	"go.uber.org/fx"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
)

// PersistanceModule for fx
//
// Provides:
//
// - The database connection
// - badaas-orm auto-migration
// - The repositories
var PersistanceModule = fx.Module(
	"persistence",
	// Database connection
	fx.Provide(gormdatabase.SetupDatabaseConnection),
	// auto-migrate
	orm.AutoMigrate,
	// repositories
	fx.Provide(repository.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(repository.NewCRUDRepository[models.User, uuid.UUID]),
)
