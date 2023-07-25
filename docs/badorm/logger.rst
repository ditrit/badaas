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
Any logger that complies with `gorm.io/gorm/logger.Interface` can be configured 
(for more information visit <https://gorm.io/docs/logger.html>).

Log levels
------------------------------

The log levels provided by badorm are the same as those of gorm:

- `logger.Error`: To only view error messages in case they occur during the execution of a sql query
- `logger.Warn`: The previous level plus warnings for execution of queries that take 
  longer than a certain time (configurable with SlowThreshold, 200ms by default).
- `logger.Info`: The previous level plus information messages for each query executed

Default logger
-------------------------------

BaDORM provides a default logger that will print Slow SQL and happening errors. 
Although in principle this logger may look the same as the gorm default logger, 
the badorm default logger will prevent stacktrace from displaying internal BaDORM files.

You can create one with the default configuration using 
(take into account that logger is github.com/ditrit/badaas/badorm/logger 
and gormLogger is gorm.io/gorm/logger):

.. code-block:: go

  logger.Default

or use `logger.New` to customize it:

.. code-block:: go

  logger.New(gormLogger.Config {
    SlowThreshold:             200 * time.Millisecond,
    LogLevel:                  gormLogger.Warn,
    IgnoreRecordNotFoundError: false,
    Colorful:                  true,
  })

The LogLevel is also configurable via the `LogMode` method. 

**Example**

.. code-block:: bash

  standalone/example.go:30 [10.392ms] [rows:1] INSERT INTO "products" ("id","created_at","updated_at","deleted_at","string","int","float","bool") VALUES ('4e6d837b-5641-45c9-a028-e5251e1a18b1','2023-07-21 17:19:59.563','2023-07-21 17:19:59.563',NULL,'',1,0.000000,false)

Zap logger
------------------------------

BaDORM provides the possibility to use `zap <https://github.com/uber-go/zap>`_ as logger. 
For this, there is a package called `gormzap` which is the compatibility layer between both loggers. 
The information displayed by the zap logger will be the same as if we were using the default logger 
but in a structured form, with the following information:

* level: ERROR, WARN or DEBUG
* message:

  * query_error for errors during the execution of a query
  * query_slow for slow queries
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
Note that the zap logger has its own level setting, so the lower of the two settings 
will be the one finally used.

**Example**

.. code-block:: bash

  DEBUG	fx/example.go:107	query_exec	{"elapsed_time": "3.291981ms", "rows_affected": "1", "sql": "SELECT products.* FROM \"products\" WHERE products.int = 1 AND \"products\".\"deleted_at\" IS NULL"}



