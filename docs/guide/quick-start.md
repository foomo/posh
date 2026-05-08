# Quick Start

This page walks from a fresh repo to a running posh shell in under a minute.

## Prerequisites

- `posh` on your `$PATH` — see [Installation](./installation)
- **Go 1.26+** for building the project shell binary
- **Git** — `posh init` infers your Go module path from the git remote

## 1. Initialise

```shell
$ cd path/to/your/project
$ posh init
```

This renders the scaffold into your repo. The notable additions:

```text
.
├── .posh.yaml          # configuration
├── .posh/
│   ├── main.go         # your shell's entrypoint
│   ├── go.mod
│   ├── internal/
│   │   ├── plugin.go   # the Plugin implementation
│   │   ├── command/    # custom commands (welcome.go is the seed)
│   │   └── config/     # typed config snippets
│   └── scripts/ownbrew/example.sh
└── Makefile            # gains shell / shell.build targets if absent
```

You can re-run `posh init --override` to re-render after upgrading; `--dry` prints the file list without writing anything.

::: tip
`posh init` must run inside a Go module. If you don't have one yet:

```shell
git init && git remote add origin <YOUR_REMOTE>
go mod init github.com/you/your-project
posh init
```
:::

## 2. Build the project shell

```shell
$ make shell.build
```

This compiles `.posh/main.go` into `bin/posh`. From here on, **`bin/posh`** is your project's shell binary — distinct from the global `posh` you ran for `init`.

## 3. Open the shell

```shell
$ make shell
```

You should see a banner and prompt:

```text
     ┓
┏┓┏┓┏┣┓
┣┛┗┛┛┛┗
┛
posh ›
```

Try the built-ins:

```text
posh › help
posh › welcome
Hi, thanks for using POSH!
posh › exit
```

## 4. Add your first command

Edit `.posh/internal/command/welcome.go` (rename it to taste) or create a new file alongside. The minimum a command needs is `Name`, `Description` and `Execute`. Wire it up in `.posh/internal/plugin.go`:

```go
inst.commands.MustAdd(
    icommand.NewWelcome(l, icommand.WelcomeWithConfigKey("welcome")),
)

// add yours
inst.commands.MustAdd(
    icommand.NewMyCommand(l),
)
```

Rebuild and run:

```shell
$ make shell.build
$ make shell
posh › my-command
```

For the full command authoring guide, see [Writing Commands](/plugin/writing-commands).

## 5. Pin your tools

`.posh.yaml` has an `ownbrew.packages` section. Add an entry for each external tool your project depends on:

```yaml
ownbrew:
  binDir: bin
  tapDir: .posh/scripts/ownbrew
  tempDir: .posh/tmp
  cellarDir: .posh/bin
  packages:
    - name: gotsrpc
      tap: foomo/tap/foomo/gotsrpc
      version: 2.6.2
```

Then in your shell:

```text
posh › brew
```

ownbrew downloads, version-checks and symlinks the binary into `bin/`. The shell's PATH already includes `${PROJECT_ROOT}/bin`, so the tool is immediately available.

## Where next

- [The Interactive Prompt](/usage/prompt) — completion, history, aliases
- [Configuration](/usage/configuration) — every key in `.posh.yaml`
- [Plugin Overview](/plugin/overview) — how the pieces connect once you start customising
