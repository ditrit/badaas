package services

import (
	"fmt"

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
	fx.Provide(addModel[models.User]),
	fx.Provide(addModel[models.Session]),
	// repositories
	fx.Provide(repository.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(repository.NewCRUDRepository[models.User, uuid.UUID]),

	// services
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)

var EAVServiceModule = fx.Module(
	"eavService",
	// models
	fx.Provide(addModel[models.EntityType]),
	fx.Provide(addModel[models.Entity]),
	fx.Provide(addModel[models.Value]),
	fx.Provide(addModel[models.Attribute]),
	// repositories
	fx.Provide(repository.NewValueRepository),
	fx.Provide(repository.NewEntityRepository),
	fx.Provide(repository.NewEntityTypeRepository),

	// service
	fx.Provide(NewEAVService),
)

func GetCRUDServiceModule[T models.Tabler]() fx.Option {
	return fx.Module(
		fmt.Sprintf(
			"%TCRUDServiceModule",
			*new(T),
		),
		// models
		fx.Provide(addModel[T]),
		// repository
		fx.Provide(repository.NewCRUDRepository[T, uuid.UUID]),
		// service
		fx.Provide(NewCRUDService[T, uuid.UUID]),
	)
}

type addModelResult struct {
	fx.Out

	Model any `group:"modelTables"`
}

func addModel[T any]() addModelResult {
	return addModelResult{
		Model: *new(T),
	}
}
