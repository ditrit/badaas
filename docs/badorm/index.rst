==============================
BaDORM |version|
==============================

BaDORM stands for Backend and Distribution ORM (Object Relational Mapping). 
It's the BaDaaS' component that allows for easy and safe persistence and querying of objects but 
it can be used both within a BaDaaS application and independently.

BaDORM is built on top of `gorm <https://gorm.io/>`_, 
a library that actually provides the functionality of an ORM: mapping objects to tables in the SQL database. 
While gorm does this job well with its automatic migration 
then performing queries on these objects is somewhat limited, 
forcing us to write SQL queries directly when they are complex. 
BaDORM seeks to address these limitations with a query system that:

- Is compile-time safe: 
  the BaDORM query system is validated at compile time to avoid errors such as 
  comparing attributes that are of different types, 
  trying to use attributes or navigate relationships that do not exist, 
  using information from tables that are not included in the query, etc.
- Is easy to use: 
  the use of this system does not require knowledge of databases, 
  SQL languages or complex concepts. 
  Writing queries only requires programming in go and the result is easy to read.
- Is designed for real applications: 
  the query system is designed to work well in real-world cases where queries are complex, 
  require navigating multiple relationships, performing multiple comparisons, etc.
- Is designed so that developers can focus on the business model: 
  its queries allow easy retrieval of model relationships to apply business logic to the model 
  and it provides mechanisms to avoid errors in the business logic due to mistakes in loading 
  information from the database.
- It is designed for high performance: 
  the query system avoids as much as possible the use of reflection and aims 
  that all the necessary model data can be retrieved in a single query to the database.

To quickly see how BaDORM can be used you can read the :doc:`quickstart`.

.. TODO
.. conceptos
..    model
..    service
..    repositorio
..    transaccion
..    conditions: dynamic, multitype, unsafe
..    operators
..    stand-alone: stand-alone (otra vez?) y con fx
..    migracion
..    gorm tags
..    mysql, sqlserver, etc
..    preloading
..    coneccion a la base de datos
..    errores?
..    uuid

.. definicion del modelo: gorm, gorm tags, base models, uuid, relaciones, punteros, etc
.. persistencia: getModels y automigracion
.. connecion a la base de datos, bases de datos soportadas
.. Generacion de condiciones
.. Creacion de los servicios y repositorios
.. utilizacion de las condiciones y operadores
.. preloading y getters
