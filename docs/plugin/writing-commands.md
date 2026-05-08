# Writing Commands

A command is a Go type satisfying `pkg/command.Command`. The base interface is three methods; everything else is opt-in.

## Minimum viable command

```go
package command

import (
    "context"

    "github.com/foomo/posh/pkg/log"
    "github.com/foomo/posh/pkg/readline"
)

type Hello struct {
    l log.Logger
}

func NewHello(l log.Logger) *Hello {
    return &Hello{l: l}
}

func (c *Hello) Name() string        { return "hello" }
func (c *Hello) Description() string { return "say hello" }
func (c *Hello) Execute(ctx context.Context, r *readline.Readline) error {
    name := "world"
    if r.Args().LenIs(1) {
        name = r.Args().At(0)
    }
    c.l.Successf("hello, %s!", name)
    return nil
}
```

Register it in `internal/plugin.go`:

```go
inst.commands.Add(NewHello(l))
```

Rebuild (`make shell.build`) and:

```text
posh › hello kevin
hello, kevin!
```

## The base interface

```go
type Command interface {
    Name() string
    Description() string
    Execute(ctx context.Context, r *readline.Readline) error
}
```

`r` (`*readline.Readline`) gives you parsed structure:

| Method | Returns |
| --- | --- |
| `r.Cmd()` | The command name as typed |
| `r.Args()` | Positional args (`Len`, `At`, `LenIs`, `LenLte`) |
| `r.Flags()` | Parsed flags |
| `r.AdditionalArgs()` | Anything after `--` |
| `r.Mode()` | Current parser mode (used for completion) |

The convention is: don't write your own argv parser, use what `readline` already gave you.

## Optional behaviours

Add capabilities by also implementing one or more of these. The framework uses Go type assertions; there's no super-interface.

### `Helper` — long help

```go
func (c *Hello) Help(ctx context.Context, r *readline.Readline) string {
    return `Print a greeting.

Usage:
  hello [name]
`
}
```

`help <command>` prints whatever you return.

### `Validator` — pre-execute check

```go
func (c *Hello) Validate(ctx context.Context, r *readline.Readline) error {
    if r.Args().Len() > 1 {
        return errors.New("hello takes 0 or 1 arguments")
    }
    return nil
}
```

Returning an error skips `Execute` and prints the error to the prompt. Use this for argument-shape validation; do business-logic errors inside `Execute`.

### `Shutdowner` — graceful teardown

```go
func (c *Server) Shutdown(ctx context.Context) error {
    return c.srv.Shutdown(ctx)
}
```

Called when the prompt loop exits. Every `Shutdowner` runs in parallel with a 3-second deadline. If you hold long-running resources (HTTP servers, database connections, watchers), implement this.

### Completion

Three flavours, picked by the parser's mode:

| Interface | Triggered when… |
| --- | --- |
| `Completer` | Generic; picked if no more specific completer exists |
| `ArgumentCompleter` | The cursor is in an argument position |
| `FlagCompleter` | The cursor is in a `--flag` position |
| `AdditionalArgsCompleter` | The cursor is after `--` |
| `PassThroughFlagsCompleter` | Flags forwarded to a wrapped CLI |

```go
import "github.com/foomo/posh/pkg/prompt/goprompt"

func (c *Hello) CompleteArguments(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
    return []goprompt.Suggest{
        {Text: "kevin", Description: "the maintainer"},
        {Text: "world", Description: "the default"},
    }
}
```

Suggestions are filtered by the prompt's filter (`goprompt.FilterCombined` by default) — you don't need to do prefix matching yourself, just return everything plausible.

For dynamic completions (e.g. "list of running pods"), keep them cheap or memoise. The prompt calls your completer on every keystroke; a 200ms RPC will feel awful.

## Functional options pattern

The repo is consistent on this — every constructor in `pkg/` follows it, and your commands should too:

```go
type (
    Hello struct {
        l       log.Logger
        cfg     config.Hello
        nameFmt string
    }
    HelloOption func(*Hello) error
)

func HelloWithConfig(v config.Hello) HelloOption {
    return func(o *Hello) error { o.cfg = v; return nil }
}

func HelloWithConfigKey(v string) HelloOption {
    return func(o *Hello) error {
        return viper.UnmarshalKey(v, &o.cfg)
    }
}

func NewHello(l log.Logger, opts ...HelloOption) (*Hello, error) {
    inst := &Hello{l: l, nameFmt: "%s"}
    for _, opt := range opts {
        if opt != nil {
            if err := opt(inst); err != nil { return nil, err }
        }
    }
    return inst, nil
}
```

The `WithConfigKey` variant is the one you'll reach for most often — it lets the registration site pick which `.posh.yaml` subtree to read:

```go
inst.commands.MustAdd(
    NewHello(l, HelloWithConfigKey("hello")),
)
```

## Tree (subcommand) commands

For commands with their own command tree (`docker run …`, `kubectl get pods`), use `pkg/command/tree`:

```go
import (
    "github.com/foomo/posh/pkg/command/tree"
)

func NewKube(l log.Logger) *Kube {
    return &Kube{
        l: l,
        tree: tree.New(&tree.Node{
            Name:        "kube",
            Description: "kubernetes helpers",
            Nodes: []*tree.Node{
                {
                    Name:        "ctx",
                    Description: "switch context",
                    Args: tree.Args{{
                        Name: "context",
                        Suggest: func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
                            return listKubeContexts(ctx)
                        },
                    }},
                    Execute: func(ctx context.Context, r *readline.Readline) error {
                        return switchContext(ctx, r.Args().At(1))
                    },
                },
                {
                    Name:        "logs",
                    Description: "tail pod logs",
                    Args: tree.Args{
                        {Name: "pod", Description: "target pod"},
                    },
                    Execute: func(ctx context.Context, r *readline.Readline) error {
                        return tailLogs(ctx, r.Args().At(1))
                    },
                },
            },
        }),
    }
}

func (c *Kube) Name() string                                    { return "kube" }
func (c *Kube) Description() string                             { return "kubernetes helpers" }
func (c *Kube) Help(ctx context.Context, r *readline.Readline) string {
    return c.tree.Help(ctx, r)
}
func (c *Kube) Execute(ctx context.Context, r *readline.Readline) error {
    return c.tree.Execute(ctx, r)
}
func (c *Kube) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
    return c.tree.Complete(ctx, r)
}
```

The tree handles completion, dispatch and help for you. Each `Node` can have nested `Nodes`, typed `Args`, `Flags`, and an `Execute` callback. See `pkg/command/tree/root_test.go` for a worked example.

## Testing commands

Commands are plain Go types — test them as such. The framework provides `pkg/log.NewTest(t)` for a logger that forwards messages to `t.Log` (and `t.Error`/`t.Fatal` at the corresponding levels):

```go
func TestHello_Execute(t *testing.T) {
    l := log.NewTest(t, log.TestWithLevel(log.LevelDebug))
    cmd := NewHello(l)

    r, err := readline.New(l)
    require.NoError(t, err)
    require.NoError(t, r.Parse("hello kevin"))

    require.NoError(t, cmd.Execute(t.Context(), r))
}
```

For commands that wrap external processes, see [Integrations § The exec package](./integrations#the-exec-package) — `middleware.CaptureStdout(&buf)` makes them straightforward to assert against in tests.
