package services

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
)

var AuthServiceModule = fx.Module(
	"authService",
	// models
	fx.Provide(getAuthModels),
	// repositories
	fx.Provide(badorm.NewCRUDRepository[models.Session, badorm.UUID]),
	fx.Provide(badorm.NewCRUDRepository[models.User, badorm.UUID]),

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
