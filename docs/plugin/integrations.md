# Integrations

Beyond the command interface, posh ships a handful of focused packages you'll reach for repeatedly when authoring a plugin. This page covers the four most important: `pkg/exec`, `pkg/require`, ownbrew, and `pkg/log`.

## The exec package

`pkg/exec` wraps `os/exec.Cmd` with a middleware chain. Use it instead of calling `exec.Cmd.Run()` directly when you want cross-cutting concerns like logging, env injection or dry-run.

### Direct use

```go
import "github.com/foomo/posh/pkg/exec"

err := exec.NewCommand(ctx, "kubectl", "get", "pods").
    Dir(env.ProjectRoot()).
    Env("KUBECONFIG=" + cfgPath).
    Run()
```

`NewCommand` defaults to inheriting `os.Environ()`, with `Stdout`/`Stderr` wired to the parent process. The fluent setters mirror `exec.Cmd`'s fields. `Run()` is a wrapper around `exec.Run(ctx, cmd, middlewares...)`.

### Middleware

A middleware is `func(next Handler) Handler` where `Handler = func(ctx context.Context, cmd *exec.Cmd) error`. The framework ships a few in `pkg/exec/middleware`:

| Middleware | What it does |
| --- | --- |
| `middleware.WithEnv(vars...)` | Appends env vars to `cmd.Env` |
| `middleware.CaptureStdout(buf)` | Redirects stdout to a buffer |
| `middleware.CaptureStderr(buf)` | Redirects stderr to a buffer |

Compose them with `Middleware(...)`:

```go
var stdout bytes.Buffer

err := exec.NewCommand(ctx, "kubectl", "version", "--client", "-o", "json").
    Middleware(
        middleware.WithEnv("KUBECONFIG=" + cfgPath),
        middleware.CaptureStdout(&stdout),
    ).
    Run()
```

Middlewares are applied right-to-left around `cmd.Run()`, so the first one passed is the outermost. Author your own:

```go
func WithTimeout(d time.Duration) exec.Middleware {
    return func(next exec.Handler) exec.Handler {
        return func(ctx context.Context, cmd *exec.Cmd) error {
            ctx, cancel := context.WithTimeout(ctx, d)
            defer cancel()
            return next(ctx, cmd)
        }
    }
}
```

Test commands by replacing the middleware chain with a fake that asserts on `cmd.Args` and returns canned output.

### Why not just `pkg/shell`?

`pkg/shell` is for the prompt's *fallback*: it runs `sh -c <line>`. That's the right tool when the input is a free-form shell line. For *programmatic* invocations of specific binaries, `pkg/exec` is safer — no quoting bugs, no PATH ambiguity, and middleware composes.

## Require checks

`pkg/require` wraps the [`fender`](https://github.com/foomo/fender) validation library to express preflight checks. The framework already uses it for `envs`, `scripts` and `packages` from `.posh.yaml`. You can extend it.

### Built-ins

```go
require.Envs(l, cfg.Envs)         // checks env vars are set
require.Packages(l, cfg.Packages) // host package + version checks
require.Scripts(l, cfg.Scripts)   // run a script; non-zero fails
require.GitUser(l, rules...)      // git user.name / user.email rules
```

`require.GitUser` accepts `GitUserName` and `GitUserEmail("regex")` rules — useful for enforcing identity conventions in monorepos.

### Composing your own

`require.First(ctx, l, fends...)` accepts mixed types — single `Fend`, slice of `Fend`, or `fend.Fends`. Append your own check:

```go
func (p *Plugin) Require(ctx context.Context, cfg config.Require) error {
    return require.First(ctx, p.l,
        require.Envs(p.l, cfg.Envs),
        require.Packages(p.l, cfg.Packages),
        require.Scripts(p.l, cfg.Scripts),
        myDockerCheck(p.l),
    )
}

func myDockerCheck(l log.Logger) fend.Fend {
    return fend.Var("", func(ctx context.Context, _ string) error {
        if err := exec.NewCommand(ctx, "docker", "info").
            Stdout(io.Discard).
            Run(); err != nil {
            return errors.New("Docker daemon is not running. Please start Docker Desktop.")
        }
        return nil
    })
}
```

`require.First` short-circuits on the first failure, so order checks from cheapest to most expensive. Use `fend.Var(...)` to wrap arbitrary functions as `fend.Fend` values.

You can also surface a `require` instance as a *runtime* checker via `prompt.WithCheckers(...)` — it'll run before the prompt opens and surface its results inline, without exiting the shell.

## Ownbrew

[`foomo/ownbrew`](https://github.com/foomo/ownbrew) is the version-pinned package manager that ships with posh. The default `Plugin.Brew` implementation just constructs an `ownbrew.Brew` from `.posh.yaml#ownbrew` and calls `Install`.

### Local packages

Create a script in `.posh/scripts/ownbrew/<name>.sh`:

```bash
#!/usr/bin/env bash
set -euo pipefail

# Available env vars:
#   $OWNBREW_NAME      package name
#   $OWNBREW_VERSION   requested version
#   $OWNBREW_OS        darwin / linux
#   $OWNBREW_ARCH      amd64 / arm64
#   $OWNBREW_TEMP_DIR  scratch space
#   $OWNBREW_CELLAR_DIR target install dir

curl -L "https://example.com/${OWNBREW_NAME}/${OWNBREW_VERSION}/${OWNBREW_OS}_${OWNBREW_ARCH}.tar.gz" \
  | tar -xz -C "$OWNBREW_CELLAR_DIR"
```

Reference it from `.posh.yaml`:

```yaml
ownbrew:
  packages:
    - name: my-tool
      version: 1.4.2
```

### Remote packages

Hosted in [`foomo/ownbrew-tap`](https://github.com/foomo/ownbrew-tap):

```yaml
ownbrew:
  packages:
    - name: gotsrpc
      tap: foomo/tap/foomo/gotsrpc
      version: 2.6.2
```

### Tag filtering

`posh brew --tags ci` only installs packages whose tag list matches. Use `--tags=-ci` to *exclude*. Wire this into your CI image build to skip dev-only tools:

```yaml
ownbrew:
  packages:
    - name: gotsrpc
      tap: foomo/tap/foomo/gotsrpc
      version: 2.6.2
      tags: [ci, dev]
    - name: docs-generator
      tap: ...
      version: ...
      tags: [dev]   # skipped on `--tags ci`
```

## Logging

`pkg/log` defines a single `Logger` interface used by every framework component and handed to every command constructor. The default implementation prints with [pterm](https://github.com/pterm/pterm) colour styling.

### Levels

```go
l.Trace("very verbose")
l.Debug("debug detail")
l.Info("informational")
l.Success("operation succeeded")  // green check
l.Warn("warning")
l.Error("error")
l.Fatal("error and exit")          // calls os.Exit(1)
```

The `Successf`/`Warnf`/etc. variants take a printf format string. `Print`/`Printf` write without a level prefix — use them for command output that the user is supposed to read as data.

### Named loggers

```go
l := p.l.Named("kube") // -> "[kube] info: …"
```

Use `Named` to scope log output for a sub-component. The framework already does this for built-ins (`prompt`, `history`, etc.).

### slog interop

`Logger.SlogHandler() slog.Handler` returns a `log/slog` handler so you can pass the logger to libraries that expect `slog`. Ownbrew uses this internally.

### Test logger

```go
import "github.com/foomo/posh/pkg/log"

func TestSomething(t *testing.T) {
    l := log.NewTest(t, log.TestWithLevel(log.LevelDebug))
    // ... pass l to your code under test
}
```

`log.NewTest(t)` forwards every log entry to `t.Log` (and `t.Error`/`t.Fatal` for `Error`/`Fatal` levels) so failed-test output shows the full sequence inline. Pass `log.TestWithLevel(...)` to surface lower-level entries.

## Putting it together

A real-world command often looks like this — config-driven, exec-wrapped, log-aware:

```go
func (c *Deploy) Execute(ctx context.Context, r *readline.Readline) error {
    env := r.Args().At(0)
    target, ok := c.cfg.Environments[env]
    if !ok {
        return fmt.Errorf("unknown environment %q", env)
    }

    c.l.Infof("deploying to %s", env)
    return exec.NewCommand(ctx, "kubectl", "apply", "-f", target.Manifest).
        Middleware(
            middleware.WithEnv("KUBECONFIG=" + target.Kubeconfig),
        ).
        Run()
}
```

A few lines of business logic; the rest is structural. That's the shape posh is going for.
