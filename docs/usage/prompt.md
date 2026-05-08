# The Interactive Prompt

`make shell` starts the project's REPL. Under the hood it runs `bin/posh prompt`, which calls `Plugin.Prompt(ctx, cfg)` with the `prompt` block from `.posh.yaml`.

## What you see

```text
     ┓
┏┓┏┓┏┣┓
┣┛┗┛┛┛┗
┛
posh ›
```

The banner is configurable (`prompt.title`); the cursor prefix is `prompt.prefix` followed by `›`. With `prompt.prefixGit: true`, the prefix gains the current branch name and tag annotations:

```text
posh (feature/shutdowner) › 
posh (v1.2.0  v1.2-rc1) ›
```

## Tab completion

Press `<Tab>` at any point to see suggestions. The completion engine knows about:

- The list of registered commands (top-level)
- Each command's flags, arguments and pass-through args (if it implements the relevant `*Completer` interface — see [Writing Commands](/plugin/writing-commands#completion))
- File paths in argument position, when the command opts into `goprompt.FilePathCompletion`

There's no `--help` discovery flag like cobra has — instead, type `help <command>` for any command implementing `Helper`.

## History

Every line you submit is appended to `prompt.history.filename` (default `.posh/.history`). On startup, the file is loaded back so `↑`/`↓` walk through previous entries.

The history file is protected by `prompt.history.lockFilename` so that two shells in the same project don't corrupt each other's writes. `prompt.history.limit` (default 100) caps the file size.

To inspect or clear history from inside the prompt:

```text
posh › history
```

The `history` command is added automatically whenever a non-noop history is configured.

## Aliases

Aliases are longest-prefix string substitutions defined in `.posh.yaml`:

```yaml
prompt:
  aliases:
    k: kubectl
    g: git
    gco: git checkout
```

Typing `gco main` is rewritten to `git checkout main` before the line is parsed, so aliases work for the shell-fallback path too. Aliases also surface as suggestions in tab completion.

## Exit and interrupts

| Keystroke | Effect |
| --- | --- |
| `Ctrl+C` | Cancels the **active command** (its `ctx.Done()` fires). The prompt itself stays open. |
| `Ctrl+D` / `exit` | Cleanly exits the prompt. Triggers `Shutdown(ctx)` on every command implementing `Shutdowner`, with a 3-second deadline. |
| `Ctrl+L` | Clears the screen. |
| `Alt+←` / `Alt+→` | Word-wise cursor movement (the macOS bindings are wired for you). |

The shutdown phase is parallel: every `Shutdowner` runs in its own goroutine inside an `errgroup`. If any returns an error, the prompt exits with a non-zero status. If a shutdown exceeds 3 seconds, the parent context is cancelled — your handlers should respect that.

## The fallback to `sh`

If the parsed command is **not** in the registry, the prompt sends the raw line to `pkg/shell` (`sh -c <line>`). This means everything you'd normally do in your shell still works:

```text
posh › ls -la bin
posh › find . -name '*.go' | wc -l
posh › PATH | tr ':' '\n'
```

Add an alias if you find yourself reaching for the same fallback often.

## Pre-flight checks

Before the first prompt is drawn, the configured `prompt.WithCheckers` are run. The seeded plugin includes a stub:

```go
prompt.WithCheckers(
    func(ctx context.Context, l log.Logger) []check.Info {
        return []check.Info{{
            Name:   "example",
            Note:   "all good",
            Status: check.StatusSuccess,
        }}
    },
),
```

Use this to surface anything the operator should know about *before* they start typing — e.g. "you're on a stale branch", "your local kube context is `prod`", "the dev DB is currently down".

## Custom prompt options

The plugin can pass arbitrary `c-bata/go-prompt` options through `prompt.WithPromptOptions(...)`, e.g. for changing the suggestion colour scheme, the maximum suggestion list length, or the keybinding for autocomplete. See `pkg/prompt/prompt.go` for the option list applied by default.
