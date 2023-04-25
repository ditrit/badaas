package services

import (
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"go.uber.org/fx"
)

var ServicesModule = fx.Module(
	"services",
	fx.Provide(userservice.NewUserService),
	fx.Provide(sessionservice.NewSessionService),
	fx.Provide(NewEAVService),
)
