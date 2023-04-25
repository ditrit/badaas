package persistence

import (
	uuid "github.com/google/uuid"

	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"go.uber.org/fx"
)

// PersistanceModule for fx
//
// Provides:
//
// - The database connection
//
// - The repositories
var PersistanceModule = fx.Module(
	"persistence",
	// Database connection
	fx.Provide(gormdatabase.SetupDatabaseConnection),

	// repositories
	fx.Provide(repository.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(repository.NewCRUDRepository[models.User, uuid.UUID]),

	fx.Provide(repository.NewValueRepository),
	fx.Provide(repository.NewEntityRepository),
	fx.Provide(repository.NewEntityTypeRepository),
)
