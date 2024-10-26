[![Build Status](https://github.com/foomo/posh/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/posh/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/posh)](https://goreportcard.com/report/github.com/foomo/posh)
[![Coverage Status](https://coveralls.io/repos/github/foomo/posh/badge.svg?branch=main&)](https://coveralls.io/github/foomo/posh?branch=main)
[![GoDoc](https://godoc.org/github.com/foomo/posh?status.svg)](https://godoc.org/github.com/foomo/posh)

<p align="center">
  <img alt="POSH" src=".github/assets/posh.png"/>
</p>

# Project Oriented SHELL (posh)

> think of `posh` as an interactive and hackable Makefile

## Installing

Install the latest release of the cli:

````bash
$ brew update
$ brew install foomo/tap/posh
````

## Usage

```shell
$ posh help
Project Oriented Shell (posh)

Usage:
  posh [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Print loaded configuration
  help        Help about any command
  init        Initialize a Project Oriented Shell
  version     Print the version

Flags:
  -h, --help           help for posh
      --level string   set log level (default: info) (default "info")
      --no-color       disabled colors (default: false)

Use "posh [command] --help" for more information about a command.
```

To start using posh, go into your project and run:

```bash
$ cd your/project
$ posh init
```

This will generate the standard layout for posh which can be changed as required through `.posh.yml`.

Once initialized you can start posh through:

```bash
$ make shell
```

## How to Contribute

Make a pull request...

## License

Distributed under MIT License, please see license file within the code for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_
