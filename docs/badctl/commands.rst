==============================
Commands
==============================

You can see the available commands by running::

    badctl help

For more information about the functionality provided and how to use each command use::

    badctl help [command]

gen docker
---------------------------

gen docker is the command you can use to generate the files and configurations 
necessary for your project to use BadAss in a simple way.

Depending of the `db_provider` chosen `gen` will generate the docker and 
configuration files needed to run the application in the `badaas/docker/db` 
and `badaas/config` folders respectively. It will also generate docker files 
to run a http api in `badaas/docker/api`.

The possible values for `db_provider` are `cockroachdb` and `postgres`. 
CockroachDB is recommended since it's a distributed database from its 
conception and postgres compatible.

All these files can be modified in case you need different values than 
those provided by default. For more information about the configuration 
head to :doc:`/badaas/configuration`

A Makefile will be generated for the execution of a badaas server, with the command::

    make badaas_run

gen conditions
---------------------------

gen conditions is the command you can use to generate 
conditions to query your objects using BaDORM. 
For each BaDORM Model found in the input packages a file 
containing all possible Conditions on that object will be generated, 
allowing you to use BaDORM in an easy way.

Its use is recommended through `go generate`. 
For that, you will only need to create a file with the following content::

    package conditions

    //go:generate badctl gen conditions ../models

An example can be found `here <https://github.com/ditrit/badorm-example/blob/main/standalone/conditions/badorm.go>`_.