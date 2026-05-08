# Local Dev Servers

A common posh use-case: a single command that brings up the project's full local stack — backend, frontend, database, message bus — and shuts it all down cleanly when the prompt exits.

This recipe builds a `dev` command that:

- Starts processes in parallel
- Streams their interleaved output to the prompt with a colour-coded prefix
- Implements `Shutdowner` so `Ctrl+D` (or `exit`) tears everything down within 3 seconds

## The shape

```go
// .posh/internal/command/dev.go
package command

import (
    "context"
    "os/exec"
    "sync"

    "github.com/foomo/posh/pkg/log"
    "github.com/foomo/posh/pkg/readline"
    "golang.org/x/sync/errgroup"
)

type Dev struct {
    l    log.Logger
    name string

    mu      sync.Mutex
    running []*exec.Cmd
}

func NewDev(l log.Logger) *Dev {
    return &Dev{l: l, name: "dev"}
}

func (c *Dev) Name() string        { return c.name }
func (c *Dev) Description() string { return "start local dev servers" }
func (c *Dev) Help(ctx context.Context, r *readline.Readline) string {
    return `Start the local dev stack.

Usage:
  dev
`
}
```

## Starting the stack

```go
func (c *Dev) Execute(ctx context.Context, r *readline.Readline) error {
    eg, ctx := errgroup.WithContext(ctx)

    for _, srv := range c.servers() {
        srv := srv
        eg.Go(func() error { return c.start(ctx, srv) })
    }

    return eg.Wait()
}

type server struct {
    name string
    bin  string
    args []string
}

func (c *Dev) servers() []server {
    return []server{
        {"api",  "go",   []string{"run", "./cmd/api"}},
        {"web",  "bun",  []string{"run", "dev"}},
    }
}

func (c *Dev) start(ctx context.Context, s server) error {
    cmd := exec.CommandContext(ctx, s.bin, s.args...)
    cmd.Stdout = c.l.Named(s.name)  // log.Logger satisfies io.Writer via pterm
    cmd.Stderr = c.l.Named(s.name)

    c.mu.Lock()
    c.running = append(c.running, cmd)
    c.mu.Unlock()

    if err := cmd.Run(); err != nil && ctx.Err() == nil {
        return err
    }
    return nil
}
```

A few things worth highlighting:

- `errgroup.WithContext` cancels every other goroutine when one returns an error — failure of any process tears the rest down.
- `exec.CommandContext` propagates the cancellation: when the prompt's context is cancelled (Ctrl+C or shutdown), each subprocess gets `SIGKILL`.
- `c.l.Named(s.name)` is reused as `io.Writer` — output is automatically prefixed with the server name and shown with the command's log level.

## Graceful shutdown

```go
func (c *Dev) Shutdown(ctx context.Context) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    var firstErr error
    for _, cmd := range c.running {
        if cmd.Process == nil {
            continue
        }
        // Ask politely first.
        if err := cmd.Process.Signal(os.Interrupt); err != nil && firstErr == nil {
            firstErr = err
        }
    }
    return firstErr
}
```

The framework calls `Shutdown(ctx)` on every `Shutdowner` in parallel, with a **3-second deadline** rooted in the parent context. Sending SIGINT here gives processes a chance to clean up — if any hold past the deadline, `exec.CommandContext` cancels them via SIGKILL anyway.

## Wire it up

```go
// .posh/internal/plugin.go
inst.commands.MustAdd(icommand.NewDev(l))
```

Rebuild and test:

```text
$ make shell.build && make shell
posh › dev
[api] info: listening on :8080
[web] info: VITE v6.0.0  ready in 412 ms
[web] info:   ➜  Local:   http://localhost:5173/
[api] info: GET / 200 1.2ms
…
^C
posh › exit
```

`Ctrl+C` cancels the *active command* — the dev stack shuts down but the prompt stays open. `exit` triggers the same path via `Shutdown`.

## Next steps

- **Health checks before signalling ready.** Wrap each server in a small "wait until port X is open" helper before declaring start-up done.
- **Per-server flags.** Use `tree.Args` to accept `dev api`, `dev web`, `dev all` (default).
- **Restart on file change.** Combine with [`fsnotify`](https://github.com/fsnotify/fsnotify) inside the command. Posh doesn't dictate the watcher — bring your own.

A more elaborate version of this pattern, plus a "service registry" config block for `.posh.yaml`, lives in the maintainer's [posh-providers](https://github.com/foomo/posh-providers) repo.
