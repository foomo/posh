# Project Layout

After `posh init`, your project gains the following structure. None of it is hidden — it's all checked into git so every contributor gets the same shell.

```text
your-project/
├── .gitignore                    # adds bin/, .posh/.history*, .posh/tmp/
├── .makerc                       # shared make include (POSH_VERSION, etc.)
├── .posh.yaml                    # configuration (see Configuration)
├── .posh/                        # the project shell binary's source
│   ├── .gitignore                # ignores compiled artifacts
│   ├── go.mod                    # depends on github.com/foomo/posh
│   ├── main.go                   # cmd.Init(internal.New); cmd.Execute()
│   ├── internal/
│   │   ├── plugin.go             # implements pkg/plugin.Plugin
│   │   ├── command/              # your custom commands
│   │   │   └── welcome.go        # seed example
│   │   └── config/
│   │       └── welcome.go        # typed structs for .posh.yaml subtrees
│   └── scripts/
│       └── ownbrew/
│           └── example.sh        # local ownbrew package script
├── bin/                          # ownbrew-managed symlinks; in $PATH
│   └── posh                      # built by `make shell.build`
└── Makefile                      # gains shell, shell.build targets
```

## Anatomy

### `.posh.yaml`

The configuration file. See [Configuration](./configuration) for every key.

### `.posh/`

Your project's *own* posh binary, as Go source. It's a regular Go module with its own `go.mod` so it can be built and tested independently of your main project's dependencies.

Two principles to keep:

1. **Don't reach across modules.** The `.posh/` module imports your project's main module only when it must (e.g. to share types). Most teams keep these decoupled.
2. **Treat it like any other Go code.** Run `gofmt`, lint it, write tests for non-trivial commands. The framework is opinionated about plumbing, not about your code quality.

### `bin/`

Where compiled binaries land. The seed config prepends this to `$PATH`, so anything in `bin/` is callable from inside the shell:

- `bin/posh` — the project shell, built from `.posh/`
- ownbrew-managed binaries (`bin/kubectl`, `bin/gotsrpc`, …) symlinked from `.posh/bin/<name>/<version>/<bin>`

### `Makefile` and `.makerc`

The seeded Makefile gains:

```make
.PHONY: shell.build
shell.build:
    cd .posh && go build -o ../bin/posh

.PHONY: shell
shell: shell.build
    bin/posh prompt
```

`.makerc` exists so multiple Makefiles in the repo can share variables (POSH version pin, project name, etc.) without duplication. It's `include`d at the top of the Makefile.

### `bin/posh execute …`

Anything you can do interactively, you can also do non-interactively from CI or another script:

```shell
$ bin/posh execute welcome
$ bin/posh execute deploy --env=staging
```

This dispatches to `Plugin.Execute(ctx, args)`, which is implemented by the scaffold to look up the command and call it directly. Convenient for "the same code path is the source of truth in the prompt and in CI".

## What stays out of git

The seeded `.gitignore` excludes:

```text
.posh/.history
.posh/.history.lock
.posh/tmp/
.posh/bin/
bin/
```

`bin/` is excluded because everything in it is reproducible — either built from your source (`bin/posh`) or managed by ownbrew (everything else). A fresh clone runs `make shell.build` to repopulate it.

## Re-rendering after upgrades

When the framework's scaffold gains new files (a new `.gitignore` entry, a new template file, a new option), pull them with:

```shell
$ posh init --override
```

This re-renders **every** scaffold file, overwriting your local copies. You almost certainly don't want that — instead, run with `--dry` to see what *would* change:

```shell
$ posh init --dry --override
```

…and copy in only the bits you want. There's no merge tooling yet; it's plain git diff.
