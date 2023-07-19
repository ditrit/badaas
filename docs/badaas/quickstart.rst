==============================
Quickstart
==============================

To quickly get badaas up and running, you can head to the 
`example <https://github.com/ditrit/badaas-example>`_. 
By following its README.md, you will see how to use badaas and it will be util 
as a template to start your own project.

Step-by-step instructions
-----------------------------------

Once you have started your project with :code:`go init`, you must add the dependency to badaas.
To use badaas, your project must also use `fx <https://github.com/uber-go/fx>`_ and
`verdeter <https://github.com/ditrit/verdeter>`_:

.. code-block:: bash

    go get -u github.com/ditrit/badaas github.com/uber-go/fx github.com/ditrit/verdeter

Then, your application must be defined as a `verdeter command` and you have to call
the configuration of this command:

.. code-block:: go

    var command = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
      Use:   "badaas",
      Short: "Backend and Distribution as a Service",
      Run:   runCommandFunc,
    })

    func main() {
      err := configuration.NewCommandInitializer().Init(command)
      if err != nil {
        panic(err)
      }

      command.Execute()
    }

Then, in the Run function of your command, you must use `fx` and start the badaas functions:

.. code-block:: go

    func runCommandFunc(cmd *cobra.Command, args []string) {
      fx.New(
        badaas.BadaasModule,

        // Here you can add the functionalities provided by badaas
        // Here you can start the rest of the modules that your project uses.
      ).Run()
    }

You are free to choose which badaas functionalities you wish to use.
To add them, you must initialise the corresponding module:

.. code-block:: go

    func runCommandFunc(cmd *cobra.Command, args []string) {
      fx.New(
        badaas.BadaasModule,

        fx.Provide(NewAPIVersion),
        // add routes provided by badaas
        badaasControllers.InfoControllerModule,
        badaasControllers.AuthControllerModule,
        badaasControllers.EAVControllerModule,
        // Here you can start the rest of the modules that your project uses.
      ).Run()
    }

    func NewAPIVersion() *semver.Version {
      return semver.MustParse("0.0.0-unreleased")
    }

For details visit :doc:`functionalities`.

Once you have defined the functionalities of your project (an http api for example),
you can generate everything you need to run your application using `badctl`.

For installing it, use:

.. code-block:: bash

    go install github.com/ditrit/badaas/tools/badctl

Then generate files to make this project work with `cockroach` as database

.. code-block:: bash

    badctl gen docker --db_provider cockroachdb

For more information about `badctl` refer to :doc:`../badctl/index`.

Finally, you can run the api with

.. code-block:: bash

    make badaas_run

The api will be available at <http://localhost:8000>.