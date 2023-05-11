package controllers

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/ditrit/badaas/services"
	"go.uber.org/fx"
)

var InfoControllerModule = fx.Module(
	"infoController",
	// controller
	fx.Provide(NewInfoController),
	// routes
	fx.Invoke(AddInfoRoutes),
)

var AuthControllerModule = fx.Module(
	"authController",
	// service
	services.AuthServiceModule,

	// controller
	fx.Provide(NewBasicAuthenticationController),

	// routes
	fx.Provide(middlewares.NewAuthenticationMiddleware),
	fx.Invoke(AddAuthRoutes),
)

var EAVControllerModule = fx.Module(
	"eavController",
	// service
	services.EAVServiceModule,

	// controller
	fx.Provide(
		fx.Annotate(
			NewEAVController,
			fx.ResultTags(`name:"eavController"`),
		),
	),

	// routes
	fx.Invoke(
		fx.Annotate(
			AddEAVCRUDRoutes,
			fx.ParamTags(`name:"eavController"`),
		),
	),
)

var CRUDControllerModule = fx.Module(
	"crudController",
	// TODO cambiar el nombre de esto
	fx.Provide(NewGeneralCRUDController),
	fx.Invoke(AddCRUDRoutes),
)

func GetCRUDControllerModule[T models.Tabler]() fx.Option {
	typeName := fmt.Sprintf("%T", *new(T))
	return fx.Module(
		typeName+"CRUDControllerModule",
		// service
		services.GetCRUDServiceModule[T](),

		// controller
		fx.Provide(
			fx.Annotate(
				NewCRUDController[T],
				fx.ResultTags(`name:"`+typeName+`CRUDController"`),
			),
		),
	)
}
