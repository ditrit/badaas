package router

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/ditrit/badaas/services"
)

// RouterModule for fx
var RouterModule = fx.Module(
	"router",
	fx.Provide(NewRouter),
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
