==============================
Logger
==============================

The logs are made directly by gorm, 
but it is possible to configure the type of logger to use, 
the logging level, among others. 
This configuration is done when connecting to the database, i.e. 
when creating the :ref:`GormDB <badorm/concepts:GormDB>` object.

As explained in the :ref:`connection section <badorm/connecting_to_a_database:Connection>`, 
this can be done by using gorm directly or by using the `badorm.ConnectToDialector` method. 
Any logger that complies with `logger.Interface` can be configured.

GORM logger
-------------------------------

Gorm provides a default logger that will print Slow SQL and happening errors.

For more information and customization visit <https://gorm.io/docs/logger.html>.

Zap logger
------------------------------

BaDORM provides the possibility to use `zap <https://github.com/uber-go/zap>`_ as logger. 
For this, there is a package called `gormzap` which is the compatibility layer between both loggers. 
The information displayed by the zap logger will be the same as if we were using the gorm logger 
but in a structured form, with the following information:

* level:

  * ERROR for errors during the execution of a query
  * WARN for slow queries: the query took longer than the SlowThreshold configured (200ms by default)
  * DEBUG for query execution
* message:

  * query_error for errors during the execution of a query
  * query_slow for slow queries: the query took longer than the SlowThreshold configured (200ms by default)
  * query_exec for query execution
* error: <error_message> (for errors only)
* elapsed_time: query execution time
* rows_affected: number of rows affected by the query
* sql: query executed

You can create one with the default configuration using:

.. code-block:: go

  gormzap.NewDefault(zapLogger)

where `zapLogger` is a zap logger, or use `gormzap.New` to customize it:

.. code-block:: go

  gormzap.New(zapLogger, gormzap.Config {
    LogLevel:                  logger.Warn,
    SlowThreshold:             200 * time.Millisecond,
    IgnoreRecordNotFoundError: false,
    ParameterizedQueries:      false,
  })

The LogLevel is also configurable via the `LogMode` method. 
Any configuration of the zap logger is done directly during its creation following the 
`zap documentation <https://pkg.go.dev/go.uber.org/zap#hdr-Configuring_Zap>`_.

.. TODO aclarar que no solo es estructurado sino que anda mejor porque no te muestra el path interno del badorm al loggear, aunque eso tambien podria intentar hacerlo con el de gorm