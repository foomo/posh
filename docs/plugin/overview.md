# Plugin Overview

A *plugin* is the Go code your project writes to turn `github.com/foomo/posh/pkg/...` into a shell tailored to your repo. It's a single Go interface — minus the boilerplate, ~20 lines.

## The contract

```go
// github.com/foomo/posh/pkg/plugin
type Plugin interface {
    Prompt(ctx context.Context, cfg config.Prompt) error
    Execute(ctx context.Context, args []string) error
    Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error
    Require(ctx context.Context, cfg config.Require) error
}
```

Each method backs one cobra subcommand on `bin/posh`:

| Method | Subcommand | Called by user as |
| --- | --- | --- |
| `Prompt` | `prompt` | `make shell` |
| `Execute` | `execute` | `bin/posh execute <cmd> [args]` (CI, scripts) |
| `Brew` | `brew` | `bin/posh brew` (often `make shell.build`) |
| `Require` | `require` | `bin/posh require` (preflight) |

The framework wires each cobra command, parses flags, loads config, and hands you a typed struct. Your job is to wire the implementation — usually **just composition** of helpers from `pkg/...`.

## The scaffolded plugin

`posh init` writes `.posh/internal/plugin.go`, a near-canonical implementation:

```go
type Plugin struct {
    l        log.Logger
    commands command.Commands
}

func New(l log.Logger) (plugin.Plugin, error) {
    inst := &Plugin{
        l:        l,
        commands: command.Commands{},
    }

    inst.commands.Add(command.NewExit(l))
    inst.commands.Add(command.NewHelp(l, inst.commands))

    inst.commands.MustAdd(
        icommand.NewWelcome(l,
            icommand.WelcomeWithConfigKey("welcome"),
        ),
    )

    return inst, nil
}
```

The constructor:

1. Builds an empty registry
2. Registers always-on built-ins (`exit`, `help`)
3. Registers your custom commands

`Prompt`, `Execute`, `Brew`, `Require` are all implemented in the same file — about 70 LOC total. Read it once, then come back here to extend.

### `Prompt`

```go
func (p *Plugin) Prompt(ctx context.Context, cfg config.Prompt) error {
    sh, err := prompt.New(p.l,
        prompt.WithContext(ctx),
        prompt.WithTitle(cfg.Title),
        prompt.WithPrefix(cfg.Prefix),
        prompt.WithAliases(cfg.Aliases),
        prompt.WithCommands(p.commands),
        prompt.WithCheckers(myChecker),
        prompt.WithFileHistory(
            history.FileWithLimit(cfg.History.Limit),
            history.FileWithFilename(cfg.History.Filename),
            history.FileWithLockFilename(cfg.History.LockFilename),
        ),
    )
    if err != nil {
        return err
    }
    return sh.Run()
}
```

Functional options for everything. Add `prompt.WithFlair(...)`, `prompt.WithPrefixGit(true)`, custom `prompt.WithFilter(...)` etc. as needed.

### `Execute`

```go
func (p *Plugin) Execute(ctx context.Context, args []string) error {
    r, err := readline.New(p.l)
    if err != nil { return err }
    if err := r.Parse(strings.Join(args, " ")); err != nil { return err }

    cmd := p.commands.Get(r.Cmd())
    if cmd == nil {
        return fmt.Errorf("invalid [cmd] argument: %s", r.Cmd())
    }
    if v, ok := cmd.(command.Validator); ok {
        if err := v.Validate(ctx, r); err != nil { return err }
    }
    return cmd.Execute(ctx, r)
}
```

Parses argv with the same readline parser the prompt uses, then dispatches to the same command instance. **Same code path interactively and in CI** is the design goal here.

### `Brew` and `Require`

Both are thin pass-throughs:

```go
func (p *Plugin) Brew(ctx context.Context, cfg ownbrewconfig.Config, tags []string, dry bool) error {
    brew, err := ownbrew.New(slog.New(p.l.SlogHandler()),
        ownbrew.WithDry(dry),
        ownbrew.WithBinDir(cfg.BinDir),
        ownbrew.WithTapDir(cfg.TapDir),
        ownbrew.WithTempDir(cfg.TempDir),
        ownbrew.WithCellarDir(cfg.CellarDir),
        ownbrew.WithPackages(cfg.Packages...),
    )
    if err != nil { return err }
    return brew.Install(ctx, tags...)
}

func (p *Plugin) Require(ctx context.Context, cfg config.Require) error {
    return require.First(ctx, p.l,
        require.Envs(p.l, cfg.Envs),
        require.Packages(p.l, cfg.Packages),
        require.Scripts(p.l, cfg.Scripts),
    )
}
```

You rarely change these. When you do, it's usually to add a custom checker (`require.First(ctx, p.l, builtins, myCheck(...))`).

## What you actually customise

In practice, 90 % of plugin authoring is:

1. Add commands to the registry in `New()`
2. Tweak `prompt.With*` options to taste
3. Add custom checkers to `Require`/`Prompt`

For everything else, see:

- [Writing Commands](./writing-commands) — the Command interface and its optional siblings
- [Integrations](./integrations) — `pkg/exec` middleware, custom `require.Fend`s, ownbrew taps, logging

## A word on dependencies

The `.posh/` module pulls in:

- `github.com/foomo/posh` (this library)
- `github.com/foomo/ownbrew` (transitively, for the brew config types)
- `github.com/spf13/viper` (for config decoding in commands)
- whatever your commands need

Keep it lean. The shell binary is rebuilt every time you change a command — fast builds matter when you're iterating.
