# Project Oriented SHELL (posh)

[![Build Status](https://github.com/foomo/posh/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/foomo/posh/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/foomo/posh)](https://goreportcard.com/report/github.com/foomo/posh)
[![Coverage Status](https://coveralls.io/repos/github/foomo/posh/badge.svg?branch=main&)](https://coveralls.io/github/foomo/posh?branch=main)
[![GoDoc](https://godoc.org/github.com/foomo/posh?status.svg)](https://godoc.org/github.com/foomo/posh)

> think of `posh` as an interactive and hackable Makefile

## Installing

Install the latest release of the cli:

````bash
$ brew update
$ brew install foomo/tap/posh
````

## Usage

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
