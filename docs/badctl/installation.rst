================
Installation
================

Install with go install
-----------------------------------

For simply installing it, use::

    go install github.com/ditrit/badaas/tools/badctl

Or you can build it from sources.

Build from sources
-----------------------------------

Get the sources of the project, either by visiting the `releases <https://github.com/ditrit/badaas/releases>`_ 
page and downloading an archive or clone the main branch (please be aware that is it not a stable version).

To build the project:

- Install `go <https://go.dev/doc/install>`_
- :code:`cd tools/badctl`
- Install project dependencies: :code:`go get`
- Run build command: :code:`go build .`

Well done, you have a binary `badctl` at the root of the project.

