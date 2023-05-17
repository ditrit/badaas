package services

import (
	"github.com/ditrit/badaas/badorm"
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
	fx.Provide(getAuthModels),
	fx.Invoke(badorm.AddModel[models.User]),
	fx.Invoke(badorm.AddModel[models.Session]),
	// repositories
	fx.Provide(badorm.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(badorm.NewCRUDRepository[models.User, uuid.UUID]),

	// services
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)

func getAuthModels() badorm.GetModelsResult {
	return badorm.GetModelsResult{
		Models: []any{
			models.Session{},
			models.User{},
		},
	}
}

var EAVServiceModule = fx.Module(
	"eavService",
	// models
	fx.Provide(getEAVModels),
	fx.Invoke(badorm.AddModel[models.EntityType]),
	fx.Invoke(badorm.AddModel[models.Entity]),
	fx.Invoke(badorm.AddModel[models.Value]),
	fx.Invoke(badorm.AddModel[models.Attribute]),
	// repositories
	fx.Provide(repository.NewValueRepository),
	fx.Provide(repository.NewEntityRepository),
	fx.Provide(repository.NewEntityTypeRepository),

	// service
	fx.Provide(NewEAVService),
)

func getEAVModels() badorm.GetModelsResult {
	return badorm.GetModelsResult{
		Models: []any{
			models.EntityType{},
			models.Entity{},
			models.Attribute{},
			models.Value{},
		},
	}
}
