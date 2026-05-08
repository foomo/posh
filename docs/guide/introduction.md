# Introduction

`posh` — short for **P**roject **O**riented **SH**ell — is a Go-based interactive shell library. The slogan is deliberate:

> Think of `posh` as an interactive, isolated and hackable Makefile.

It exists because most projects accumulate a tangle of helper scripts: a `Makefile` here, a `bin/` of bash there, a `package.json` script that calls another script that calls a Go program. Onboarding becomes a folklore exercise. `posh` replaces that folklore with a single, project-scoped REPL whose commands are real Go code you check into the repo.

## What you get

When you run `posh init` in a project, posh scaffolds a tiny Go program — the project's **own** shell binary — under `.posh/`. After `make shell.build` you can launch it with `make shell` and you land in something like:

```text
posh › help
Help about all available commands.

Usage:
  help [command]

Available Commands:
  exit                exit shell
  help                print help
  history             print history
  welcome             print a welcome message

posh › welcome
Hi, thanks for using POSH!
posh ›
```

Out of the box you get:

- **Tab completion** that's aware of your custom commands, their flags and arguments
- **History** persisted to `.posh/.history`, with file locking so concurrent shells don't trample each other
- **Aliases** declared in `.posh.yaml`
- **Pre-flight checks** (`posh require`) that fail fast when an env var, host package or smoke-test script is missing
- **Tool pinning** via [ownbrew](https://github.com/foomo/ownbrew), so every contributor uses the same `kubectl`, `node`, `gotsrpc`, …
- **A graceful shutdown path** — long-running commands implementing `Shutdown(ctx)` are torn down cleanly when the prompt exits

Anything posh doesn't recognise as a command falls through to `sh -c`, so you don't lose access to the rest of your toolchain.

## Why not just use a Makefile?

Makefiles are excellent for *describing* tasks but mediocre at *running* them interactively:

| | Makefile | posh |
| --- | --- | --- |
| Interactive REPL | ❌ | ✅ |
| Tab completion of arguments | ❌ | ✅ (per-command) |
| Static-typed command implementations | ❌ | ✅ (Go) |
| Per-project tool versions | external | built-in (ownbrew) |
| Pre-flight environment checks | ad-hoc | declarative (`posh require`) |
| Shared command library | copy-paste | Go imports |
| Aliases & history | shell-level | project-level |

A common pattern is to keep the Makefile for CI-style entry points (`make test`, `make build`) and use posh for everything humans do day-to-day.

## Why not just write a CLI in Go?

You could. `posh` is what's left after several teams did exactly that and converged on the same pieces: a cobra root, a prompt loop, a config schema, version-pinned tools, prerequisite checks. Rather than reinventing those, you import `github.com/foomo/posh/pkg/...`, implement the [`Plugin`](/plugin/overview) interface, and write only the parts that are unique to your project.

## What it is *not*

- **Not a shell replacement.** It runs *inside* your existing shell, manages a single project, and exits cleanly.
- **Not a daemon.** No background process, no socket, nothing to garbage-collect when you're done.
- **Not magic.** Your shell binary is a few hundred lines of Go that you can read in one sitting.

## Where to next

- New here? Start with [Installation](./installation) and the [Quick Start](./quick-start).
- Want to understand the moving parts? Read [Concepts](./concepts).
- Ready to add your own commands? Jump to [Plugin Authoring](/plugin/overview).
