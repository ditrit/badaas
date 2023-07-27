# badctl: the BadAas controller

`badctl` is the command line tool that makes it possible to configure and run your BadAas applications easily.

- [badctl: the BadAas controller](#badctl-the-badaas-controller)
  - [Install with go install](#install-with-go-install)
  - [Build from sources](#build-from-sources)
  - [Commands](#commands)
    - [gen](#gen)
  - [Contributing](#contributing)

## Install with go install

For simply installing it, use:

```bash
go install github.com/ditrit/badaas/tools/badctl@249d3c0
```

Or you can build it from sources.

## Build from sources

Get the sources of the project, either by visiting the [releases](https://github.com/ditrit/badaas/releases) page and downloading an archive or clone the main branch (please be aware that is it not a stable version).

To build the project:

- [Install go](https://go.dev/dl/#go1.18.4) v1.18
- `cd tools/badctl`
- Install project dependencies

    ```bash
    go get
    ```

- Run build command

    ```bash
    go build .
    ```

Well done, you have a binary `badctl` at the root of the project.

## Commands

You can see the available commands by running:

```bash
badctl help
```

For more information about the functionality provided and how to use each command use:

```bash
badctl help [command]
```

### gen

gen is the command you can use to generate the files and configurations necessary for your project to use BadAss in a simple way.

`gen` will generate the docker and configuration files needed to run the application in the `badaas/docker` and `badaas/config` folders respectively.

All these files can be modified in case you need different values than those provided by default. For more information about the configuration head to [configuration docs](../../configuration.md).

A Makefile will be generated for the execution of a badaas server, with the command:

```bash
make badaas_run
```

## Contributing

You can make modifications to the badctl source code and compile it locally with:

```bash
go build .
```

You can then run the badctl executable directly or add a link in your $GOPATH to run it from a project:

```bash
ln -sf badctl $GOPATH/bin/badctl
```
