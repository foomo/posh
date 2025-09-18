[![Build Status](https://github.com/foomo/posh/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/posh/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/posh)](https://goreportcard.com/report/github.com/foomo/posh)
[![GoDoc](https://godoc.org/github.com/foomo/posh?status.svg)](https://godoc.org/github.com/foomo/posh)

<p align="center">
  <img alt="POSH" src=".github/assets/posh.png"/>
</p>

# Project Oriented SHELL (posh)

> Think of `posh` as an interactive, isolated and hackable Makefile

## Installation

### Download binary

Download a [binary release](https://github.com/foomo/posh/releases)

### Build from source

```
go install github.com/foomo/posh@latest
```

### Homebrew (Linux/macOS)

If you use [Homebrew](https://brew.sh), you can install like this:
```
brew install foomo/tap/posh
```

### Mise

If you use [mise](https://https://mise.jdx.dev), you can install like this:
```
mise use github.com:foomo/posh
```

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

Once initialized, you can start posh through:

```bash
$ make shell
```

## How to Contribute

Please refer to the [CONTRIBUTING](.gihub/CONTRIBUTING.md) details and follow the [CODE_OF_CONDUCT](.gihub/CODE_OF_CONDUCT.md) and [SECURITY](.github/SECURITY.md) guidelines.

## License

Distributed under MIT License, please see license file within the code for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_
