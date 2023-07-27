package badaas

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"
)

func TestInvokeFunctionsWithProvidedValues(t *testing.T) {
	mockObject := mockObject{}

	mockObject.On("Function", 1).Return(1)

	badaas := BaDaaSInitializer{}
	badaas.Provide(
		newIntValue,
	).Invoke(
		mockObject.Function,
		shutdown,
	).Start()
}

func TestAddModulesAreExecuted(t *testing.T) {
	mockObject := mockObject{}

	mockObject.On("Function", 1).Return(1)

	badaas := BaDaaSInitializer{}
	badaas.AddModules(
		fx.Module(
			"test module",
			fx.Provide(newIntValue),
			fx.Invoke(mockObject.Function),
		),
	).Invoke(
		shutdown,
	).Start()
}

func newIntValue() int {
	return 1
}

type mockObject struct {
	mock.Mock
}

func (o *mockObject) Function(intValue int) int {
	args := o.Called(intValue)
	return args.Int(0)
}

func shutdown(
	shutdowner fx.Shutdowner,
) {
	shutdowner.Shutdown()
}
