package controllers

import "go.uber.org/fx"

var InfoControllerModule = fx.Module(
	"infoController",
	fx.Provide(NewInfoController),
	fx.Invoke(AddInfoRoutes),
)

var AuthControllerModule = fx.Module(
	"authController",
	fx.Provide(NewBasicAuthenticationController),
	fx.Invoke(AddAuthRoutes),
)
