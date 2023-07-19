==============================
Introduction
==============================

Badaas enables the effortless construction of **distributed, resilient, 
highly available and secure applications by design**, while ensuring very simple 
deployment and management (NoOps).

.. warning::
   BaDaaS is still under development and each of its components can have a different state of evolution

Features and components
=================================

Badaas provides several key features, 
each provided by a component that can be used independently and has a different state of evolution:

- **Authentication** (unstable): Badaas can authenticate users using its internal 
  authentication scheme or externally by using protocols such as OIDC, SAML, Oauth2...
- **Authorization** (wip_unstable): On resource access, Badaas will check if the user 
  is authorized using a RBAC model.
- **Distribution** (todo): Badaas is built to run in clusters by default. 
  Communications between nodes are TLS encrypted using `shoset <https://github.com/ditrit/shoset>`_.
- **Persistence** (wip_unstable): Applicative objects are persisted as well as user files. 
  Those resources are shared across the clusters to increase resiliency. 
  To achieve this, BaDaaS uses the :doc:`BaDORM <badorm/index>` component.
- **Querying Resources** (unstable): Resources are accessible via a REST API.
- **Posix compliant** (stable): Badaas strives towards being a good unix citizen and 
  respecting commonly accepted norms. (see :doc:`badaas/configuration`)
- **Advanced logs management** (todo): Badaas provides an interface to interact with 
  the logs produced by the clusters. Logs are formatted in json by default.

Learn how to use BaDaaS following the :doc:`badaas/quickstart`.

.. toctree::
   :caption: BaDaaS

   self
   badaas/quickstart
   badaas/functionalities
   badaas/configuration

.. toctree::
   :caption: BaDctl

   badctl/index
   badctl/installation
   badctl/commands

.. toctree::
   :caption: BaDORM

   badorm/index
   badorm/quickstart
   badorm/concepts
   badorm/declaring_models
   badorm/connecting_to_a_database
   badorm/crud
   badorm/query
   badorm/advanced_query
   badorm/preloading