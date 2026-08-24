# Built-in Commands

posh ships a small set of ready-made commands you can register inside your Plugin. All of them live in `pkg/command` and
are constructed with `command.New<Name>(l, opts...)`.

## `exit`

Exits the prompt loop. Equivalent to typing `Ctrl+D`. Use this so users can quit without remembering shell shortcuts.

```go
inst.commands.Add(command.NewExit(l))
```

No options. No flags. Implementing `Shutdowner` on your own commands is what runs *during* the exit.

## `help`

Lists every command with a `Description()`, and prints the long help for any command that implements the [
`Helper`](/plugin/writing-commands#optional-behaviours) interface.

```go
inst.commands.Add(command.NewHelp(l, inst.commands))
```

Note that `help` takes the entire `commands` registry — it discovers commands by reflection, so it sees whatever you add
later in the constructor.

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

posh › help welcome
Print a welcome message

Usage:
welcome
```

## `history`

Added **automatically** by `pkg/prompt` whenever a non-noop history backend is configured — you don't register it
manually. It prints the persisted history file.

## `env`

Manages the shell's internal environment variables. Register it manually:

```go
inst.commands.Add(command.NewEnv(l))
```

| Subcommand              | Purpose                                             |
|-------------------------|-----------------------------------------------------|
| `env list`              | Table of all environment variables and their values |
| `env get <Key>`         | Print the value of a single variable                |
| `env set <Key> [Value]` | Set a variable (empty value if omitted)             |
| `env unset <Key>`       | Remove a variable                                   |

`env get` completes `<Key>` from the current environment:

```text
posh › env get P<Tab>
PROJECT_ROOT
POSH_ROOT_CONFIG_PATH
```

## `welcome` (scaffolded)

Not technically a built-in — it's the seed example in `embed/scaffold/init/$.posh/internal/command/welcome.go.gotext`.
It demonstrates the smallest plausible command:

- Reads `welcome.message` from `.posh.yaml` via `viper.UnmarshalKey`
- Implements `Name`, `Description`, `Execute`, `Help`
- Uses `WelcomeWithConfigKey("welcome")` as a functional option

Treat it as a template for your own commands and replace it once you've internalised the pattern.

## Other helpers in `pkg/command`

These aren't always-on built-ins, but the package ships a few more constructors for common needs:

| Constructor             | Purpose                                                |
|-------------------------|--------------------------------------------------------|
| `command.NewCache(...)` | A scoped cache around expensive completions or lookups |
| `command.NewCheck(...)` | Wraps a `check.Checker` as a callable command          |

See the godoc for current signatures: <https://pkg.go.dev/github.com/foomo/posh/pkg/command>.

## Agent-facing output

Every built-in above describes itself to [`posh agent catalog`](/usage/agents#posh-agent-catalog). Most also contribute
prose to their [generated skill](/usage/agents#layout) — `exit` and `help` deliberately don't, since theirs would only
steer an agent away from a command it had to load the skill to learn about. `cache` and `env` additionally print JSON instead of tables when
[agent mode](/usage/agents#agent-mode) is active — so the same command serves a human at the prompt and an AI agent in
CI. Your own commands opt in the same way, via
[`Describer`](/plugin/writing-commands#describer-agent-catalog) and
[`Skiller`](/plugin/writing-commands#skiller-agent-skill-prose).
