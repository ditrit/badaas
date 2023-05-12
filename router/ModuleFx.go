package router

import (
	"fmt"

	"github.com/ditrit/badaas/badorm"
	"github.com/ditrit/badaas/badorm/tabler"
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/ditrit/badaas/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

// RouterModule for fx
var RouterModule = fx.Module(
	"router",
	fx.Provide(NewRouter),
	fx.Invoke(
		fx.Annotate(
			AddCRUDRoutes,
			fx.ParamTags(`group:"crudControllers"`),
		),
	),
	// middlewares
	fx.Provide(middlewares.NewJSONController),
	fx.Provide(middlewares.NewMiddlewareLogger),
	fx.Invoke(middlewares.AddLoggerMiddleware),
)

var InfoRouteModule = fx.Module(
	"infoRoute",
	// controller
	fx.Provide(controllers.NewInfoController),
	// routes
	fx.Invoke(AddInfoRoutes),
)

var AuthRoutesModule = fx.Module(
	"authRoutes",
	// service
	services.AuthServiceModule,

	// controller
	fx.Provide(controllers.NewBasicAuthenticationController),

	// routes
	fx.Provide(middlewares.NewAuthenticationMiddleware),
	fx.Invoke(AddAuthRoutes),
)

var EAVRoutesModule = fx.Module(
	"eavRoutes",
	// service
	services.EAVServiceModule,

	// controller
	fx.Provide(
		fx.Annotate(
			controllers.NewEAVController,
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

func GetCRUDRoutesModule[T tabler.Tabler]() fx.Option {
	typeName := fmt.Sprintf("%T", *new(T))
	return fx.Module(
		typeName+"CRUDRoutesModule",
		// service
		badorm.GetCRUDServiceModule[T, uuid.UUID](),

		// controller
		fx.Provide(
			fx.Annotate(
				controllers.NewCRUDController[T],
				fx.ResultTags(`group:"crudControllers"`),
			),
		),
	)
}
