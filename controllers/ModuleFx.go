package controllers

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

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
	fx.Provide(
		fx.Annotate(
			NewEAVController,
			fx.ResultTags(`name:"eavController"`),
		),
	),
	fx.Invoke(
		fx.Annotate(
			AddEAVCRUDRoutes,
			fx.ParamTags(`name:"eavController"`),
		),
	),
)

var CRUDControllerModule = fx.Module(
	"crudController",
	fx.Provide(NewGeneralCRUDController),
	fx.Invoke(AddCRUDRoutes),
)

func GetCRUDModule[T models.Tabler]() fx.Option {
	return fx.Module(
		"crudModule",
		fx.Provide(repository.NewCRUDRepository[T, uuid.UUID]),
		fx.Provide(services.NewCRUDService[T, uuid.UUID]),
		fx.Provide(
			fx.Annotate(
				NewCRUDController[T],
				fx.ResultTags(
					fmt.Sprintf(
						`name:"%TCRUDController"`,
						*new(T),
					),
				),
			),
		),
	)
}
