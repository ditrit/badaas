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

var EAVControllerModule = fx.Module(
	"eavController",
	fx.Provide(NewEAVController),
	fx.Invoke(AddEAVCRUDRoutes),
)

var CRUDControllerModule = fx.Module(
	"crudController",
	fx.Provide(NewGeneralCRUDController),
	fx.Invoke(AddCRUDRoutes),
)
