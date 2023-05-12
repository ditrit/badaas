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
	fx.Provide(badorm.AddModel[models.User]),
	fx.Provide(badorm.AddModel[models.Session]),
	// repositories
	fx.Provide(badorm.NewCRUDRepository[models.Session, uuid.UUID]),
	fx.Provide(badorm.NewCRUDRepository[models.User, uuid.UUID]),

	// services
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
)

var EAVServiceModule = fx.Module(
	"eavService",
	// models
	fx.Provide(badorm.AddModel[models.EntityType]),
	fx.Provide(badorm.AddModel[models.Entity]),
	fx.Provide(badorm.AddModel[models.Value]),
	fx.Provide(badorm.AddModel[models.Attribute]),
	// repositories
	fx.Provide(repository.NewValueRepository),
	fx.Provide(repository.NewEntityRepository),
	fx.Provide(repository.NewEntityTypeRepository),

	// service
	fx.Provide(NewEAVService),
)
