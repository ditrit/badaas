==============================
Functionalities
==============================

InfoControllerModule
-------------------------------

`InfoControllerModule` adds the path `/info`, where the api version will be answered. 
To set the version we want to be responded to we must provide the version using fx:

.. code-block:: go

    func runCommandFunc(cmd *cobra.Command, args []string) {
      fx.New(
        badaas.BadaasModule,

        // provide api version
        fx.Provide(NewAPIVersion),
        // add /info route provided by badaas
        badaasControllers.InfoControllerModule,
      ).Run()

    func NewAPIVersion() *semver.Version {
      return semver.MustParse("0.0.0-unreleased")
    }

AuthControllerModule
-------------------------------

`AuthControllerModule` adds `/login` and `/logout`, 
which allow us to add authentication to our application in a simple way:

.. code-block:: go

    func runCommandFunc(cmd *cobra.Command, args []string) {
      fx.New(
        badaas.BadaasModule,

        // add /login and /logout routes provided by badaas
        badaasControllers.AuthControllerModule,
      ).Run()
    }

EAVControllerModule
-------------------------------

`EAVControllerModule` adds `/eav/objects/{type}` and `/eav/objects/{type}/{id}`, 
where `{type}` is any defined type and `{id}` is any uuid. These routes allow us to create, 
read, update and remove objects using an EAV model. For more information on how to use them, 
see the `example <https://github.com/ditrit/badaas-example>`_.

.. code-block:: go

    func runCommandFunc(cmd *cobra.Command, args []string) {
      fx.New(
        badaas.BadaasModule,

        badaasControllers.EAVControllerModule,
      ).Run()
    }