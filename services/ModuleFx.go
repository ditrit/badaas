package services

import (
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

var AuthServiceModule = fx.Module(
	"authService",
	// models
	fx.Invoke(addAuthModels),
	// repositories
	fx.Provide(repository.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(repository.NewCRUDRepository[models.User, uuid.UUID]),

	// services
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)

func addAuthModels() {
	gormdatabase.ListOfTables = append(gormdatabase.ListOfTables,
		models.User{},
		models.Session{},
	)
}

var EAVServiceModule = fx.Module(
	"eavService",
	// models
	fx.Invoke(addEAVModels),
	// repositories
	fx.Provide(repository.NewValueRepository),
	fx.Provide(repository.NewEntityRepository),
	fx.Provide(repository.NewEntityTypeRepository),

	// service
	fx.Provide(NewEAVService),
)

func addEAVModels() {
	gormdatabase.ListOfTables = append(gormdatabase.ListOfTables,
		models.EntityType{},
		models.Entity{},
		models.Value{},
		models.Attribute{},
	)
}

func GetCRUDServiceModule[T models.Tabler]() fx.Option {
	return fx.Module(
		"crudServiceModule",
		// models
		fx.Invoke(addCRUDModel[T]),
		// repository
		fx.Provide(repository.NewCRUDRepository[T, uuid.UUID]),
		// service
		fx.Provide(NewCRUDService[T, uuid.UUID]),
	)
}

func addCRUDModel[T models.Tabler]() {
	gormdatabase.ListOfTables = append(gormdatabase.ListOfTables,
		*new(T),
	)
}
