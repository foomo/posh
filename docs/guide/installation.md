# Installation

`posh` ships as a single static binary. Pick whichever channel fits your environment.

## Homebrew (macOS / Linux)

```shell
brew install foomo/tap/posh
```

The formula lives in [`foomo/homebrew-tap`](https://github.com/foomo/homebrew-tap).

## mise

```shell
mise use github:foomo/posh
```

Or run directly without installing:

```shell
mise x github:foomo/posh -- --help
```

See [mise.jdx.dev](https://mise.jdx.dev) for project-level pinning.

## Docker

Multi-arch images (`amd64`, `arm64`) are published to [Docker Hub](https://hub.docker.com/r/foomo/posh):

```shell
docker run --rm foomo/posh:latest --help
```

For project use, mount your repo and pass through the working directory:

```shell
docker run --rm -it -v "$PWD:/work" -w /work foomo/posh:latest init
```

## Binary release

Download the archive for your OS/arch from the [releases page](https://github.com/foomo/posh/releases) and extract `posh` somewhere on your `$PATH`:

```shell
curl -L https://github.com/foomo/posh/releases/latest/download/posh_$(uname -s)_$(uname -m).tar.gz \
  | tar -xz -C /usr/local/bin posh
```

## go install

```shell
go install github.com/foomo/posh@latest
```

Requires **Go 1.26+**. Note: `go install` builds without the linker flags used by release builds, so `posh version` will print `unknown`.

## Verify

```shell
$ posh version
v1.x.y

$ posh --help
Project Oriented Shell (posh)

Usage:
  posh [command]
…
```

You're ready for the [Quick Start](./quick-start).

## What gets installed

The `posh` binary is intentionally small — `init`, `config` and `version` are the only meaningful subcommands. The interesting commands (`prompt`, `execute`, `require`, `brew`) only appear inside a *scaffolded* project shell, because they delegate to the [Plugin](/plugin/overview) you write. See [Concepts](./concepts) for the full picture.
