package controllers

import "go.uber.org/fx"

// ControllerModule for fx
var ControllerModule = fx.Module(
	"controllers",
	fx.Provide(NewInfoController),
	fx.Provide(NewBasicAuthentificationController),
	fx.Provide(NewEAVController),
)
