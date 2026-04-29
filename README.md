[![Build Status](https://github.com/foomo/posh/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/posh/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/posh)](https://goreportcard.com/report/github.com/foomo/posh)
[![GoDoc](https://godoc.org/github.com/foomo/posh?status.svg)](https://godoc.org/github.com/foomo/posh)

<p align="center">
  <img alt="posh-providers" src="docs/public/logo.png" width="400" height="400"/>
</p>

# Project Oriented SHELL (posh)

> Think of `posh` as an interactive, isolated and hackable Makefile

## Installation

<details>
<summary><b>Homebrew</b> (macOS / Linux)</summary>

```shell
brew install foomo/tap/posh
```

See the [foomo/homebrew-tap](https://github.com/foomo/homebrew-tap) repository.

</details>

<details>
<summary><b>Docker</b></summary>

```shell
docker run --rm foomo/posh:latest --help
```

Multi-arch images (`amd64`, `arm64`) are published to [Docker Hub](https://hub.docker.com/r/foomo/posh).

</details>

<details>
<summary><b>mise</b></summary>

```shell
mise use github:foomo/posh
```

or run directly:

```shell
mise x github:foomo/posh -- --help
```

See [mise.jdx.dev](https://mise.jdx.dev).

</details>

<details>
<summary><b>Binary release</b></summary>

Download the archive for your OS/arch from the [releases page](https://github.com/foomo/posh/releases) and extract `posh` into your `$PATH`.

</details>

<details>
<summary><b>go install</b></summary>

```shell
go install github.com/foomo/posh/cmd/posh@latest
```

Requires Go 1.26+.

</details>

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

```shell
$ cd your/project
$ posh init
```

This will generate the standard layout for posh which can be changed as required through `.posh.yml`.

```yaml
version: v1.0

## Prompt settings
prompt:
  title: "Posh"
  prefix: "posh >"
  history:
    limit: 100
    filename: .posh/.history
    lockFilename: .posh/.history.lock

## Environment variables
env:
  - name: PATH
    value: "${PROJECT_ROOT}/bin:${PATH}"

## Ownbrew settings
ownbrew:
  binDir: "bin"
  tapDir: ".posh/scripts/ownbrew"
  tempDir: ".posh/tmp"
  cellarDir: ".posh/bin"
  packages: []
    ## Remote package
    ## See `https://github.com/foomo/ownbrew-tap`
    #- name: gotsrpc
    #  tap: foomo/tap/foomo/gotsrpc
    #  version: 2.6.2
    ## Local package `.posh/scripts/ownbrew`
    #- name: example
    #  version: 0.0.0

## Requirement settings
require:
  ## Required environment variables
  envs: []
    ## Example: require VOLTA_HOME
    #- name: VOLTA_HOME
    #  help: |
    #    Missing required $VOLTA_HOME env var.
    #
    #    Please initialize volta and ensure $VOLTA_HOME is set:
    #
    #      $ volta setup

  ## Required scripts that need to succeed
  scripts: []
    ## Example: git
    #- name: git
    #  command: |
    #    git status && exit 0 || exit 1
    #  help: |
    #    This is not a git repo. Please clone the repository

    ## Example: npm
    #- name: npm
    #  command: npm whoami --registry=https://npm.pkg.github.com > /dev/null 2>&1
    #  help: |
    #    You're not yet logged into the github npm registry!
    #
    #      $ npm login --scope=@<SCOPE> --registry=https://npm.pkg.github.com
    #      Username: [GITHUB_USERNAME]
    #      Password: [GITHUB_TOKEN]
    #      Email: [EMAIL]

  ## Required packages to be installed on the host
  packages: []
    ## Example: git
    #- name: git
    #  version: '~2'
    #  command: git version | awk '{print $3}'
    #  help: |
    #    Please ensure you have 'git' installed in the required version: %s!
    #
    #      $ brew update
    #      $ brew install git

    ## Example: go
    #- name: go
    #  version: '>=1.23'
    #  command: go env GOVERSION | cut -c3-
    #  help: |
    #    Please ensure you have 'go' installed in the required version: %s!
    #
    #      $ brew update
    #      $ brew install go

    #- name: volta
    #  version: '>=2'
    #  command: volta --version
    #  help: |
    #    Please ensure you have 'volta' installed in a recent version: %s!
    #
    #      $ curl https://get.volta.sh | bash
    #
    #    Or see the documentation: https://docs.volta.sh/guide/getting-started


## Integrations

## Example: Custom
welcome:
  message: Hi, thanks for using POSH!

```

Once initialized, you can start posh through:

```shell
$ make shell
```

## How to Contribute

Contributions are welcome! Please read the [contributing guide](docs/CONTRIBUTING.md).

![Contributors](https://contributors-table.vercel.app/image?repo=foomo/posh&width=50&columns=15)

## License

Distributed under MIT License, please see the [license](LICENSE) file within the code for more details.

_Made with ♥ [foomo](https://www.foomo.org) by [bestbytes](https://www.bestbytes.com)_
