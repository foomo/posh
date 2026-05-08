# Concepts

A small mental model goes a long way with posh. This page is the one you should read once and refer back to.

## Two binaries, one library

The repository ships **two distinct things**:

1. **The `posh` CLI** — a small static binary you install once. It has only three meaningful subcommands: `init`, `config`, `version`. Its job is to scaffold a project.
2. **A Go library** under `github.com/foomo/posh/pkg/...`. Downstream projects import this and compile their **own** shell binary at `bin/posh`.

```text
                 ┌────────────────────┐
                 │  posh (global CLI) │
                 │  brew install posh │
                 └──────────┬─────────┘
                            │ posh init
                            ▼
                 ┌────────────────────┐
                 │   .posh/main.go    │   compiled with…
                 │   .posh/internal/  │   github.com/foomo/posh/pkg/...
                 │   .posh.yaml       │
                 └──────────┬─────────┘
                            │ make shell.build
                            ▼
                 ┌────────────────────┐
                 │   bin/posh         │  ← the project's own shell
                 │   prompt / execute │
                 │   brew / require   │
                 └────────────────────┘
```

When you type `posh prompt` you're running **the project's** binary, not the global one. The same is true of `execute`, `brew`, `require` — those subcommands only exist when a [`Plugin`](/plugin/overview) is wired in.

## The plugin seam

Every project shell is a thin `main.go` that calls into the framework with one argument: a constructor for your Plugin.

```go
// .posh/main.go
package main

import (
    "your-module/posh/internal"
    "github.com/foomo/posh/cmd"
)

func init()   { cmd.Init(internal.New) }   // internal.New returns plugin.Plugin
func main()  { cmd.Execute() }
```

The framework wires cobra subcommands and delegates to your Plugin:

```go
// pkg/plugin/plugin.go
type Plugin interface {
    Prompt(ctx context.Context, cfg config.Prompt) error
    Execute(ctx context.Context, args []string) error
    Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error
    Require(ctx context.Context, cfg config.Require) error
}
```

Each method is one cobra subcommand. The scaffolded `internal/plugin.go` is the canonical implementation — it composes `pkg/prompt`, `pkg/require`, ownbrew and your custom commands. Most projects only ever change the **command registration** block.

## Commands and optional behaviours

A command is anything that satisfies:

```go
type Command interface {
    Name() string
    Description() string
    Execute(ctx context.Context, r *readline.Readline) error
}
```

The base interface is intentionally minimal. Layer in extra capabilities by also implementing optional interfaces:

| Interface | When the framework calls it |
| --- | --- |
| `Helper` | `help <cmd>` from the prompt |
| `Validator` | Before `Execute`, lets you reject bad input |
| `Shutdowner` | After the prompt exits, with a 3s deadline |
| `Completer` | Generic completion for Args/Flags/AdditionalArgs |
| `ArgumentCompleter` | Positional argument completion (preferred over `Completer`) |
| `FlagCompleter` | Completion when typing `--something` |
| `AdditionalArgsCompleter` | Anything after `--` |
| `PassThroughFlagsCompleter` | Flags passed through to a wrapped tool |

The prompt uses Go type assertions to detect what each command supports. There is no super-interface to implement — pick what you need. See [Writing Commands](/plugin/writing-commands) for examples.

## The prompt loop

`pkg/prompt.Run()` is the heart of the interactive shell. Each line you type is:

1. Trimmed; if empty, ignored
2. Appended to history (file-locked)
3. Alias-expanded (longest-prefix match against `prompt.aliases` in config)
4. Parsed by `pkg/readline` into `cmd`, `args`, `flags`, `additionalArgs`
5. Looked up in the `command.Commands` registry
6. If matched: `Validator.Validate` → `Command.Execute`
7. If not matched: handed to `pkg/shell` (`sh -c <input>`)

`Ctrl+C` cancels the active command's context but does **not** exit the prompt — typing `exit` (or `Ctrl+D`) does that. On exit, the prompt loop calls `Shutdown(ctx)` on every registered command implementing `Shutdowner`, in parallel, with a 3-second timeout.

Tab completion runs the same parser and dispatches to the appropriate optional `*Completer` interface based on the parser's mode (Args / Flags / AdditionalArgs).

## Configuration & isolation

`.posh.yaml` is loaded by viper and exposes typed structs from `pkg/config`:

- `prompt` — title, prefix, history, aliases (consumed by `Plugin.Prompt`)
- `env` — name/value pairs prepended to the process environment when the shell starts
- `ownbrew` — package list and target dirs (consumed by `Plugin.Brew`)
- `require` — env vars, host packages and smoke-test scripts (consumed by `Plugin.Require`)

Anything else in the file is yours — your Plugin can `viper.UnmarshalKey("yourkey", &yourStruct)` to read custom sections (the seeded `welcome.message` works exactly this way).

The "isolated" in the slogan means: the shell's environment is built fresh on launch from `.posh.yaml#env`. Most projects prepend `${PROJECT_ROOT}/bin` to `$PATH` so that ownbrew-installed tools win over the host versions.

## The exec package

`pkg/exec` (newer addition) wraps `os/exec.Cmd` with a middleware chain. Prefer it over calling `exec.Cmd.Run()` directly when you want cross-cutting concerns (logging, env injection, dry-run). See [Integrations](/plugin/integrations) for examples.

## Where the lines are drawn

- `internal/` in the posh repo is **for the `posh` CLI itself**. Don't import it.
- `pkg/` is the public library surface. Stable, documented, importable.
- The scaffold under `embed/scaffold/init/` is what your project gets. Updates land there first; re-run `posh init --override` to pull them in.

That's the entire model. The rest of the documentation expands on each piece.
