# Configuration

`.posh.yaml` is the single source of truth for a project shell. It is loaded by [viper](https://github.com/spf13/viper), which means:

- Keys are case-insensitive
- Nested keys can be supplied as environment variables (`POSH_PROMPT_TITLE=...`)
- Anything *not* covered by the framework is yours — your Plugin can `viper.UnmarshalKey(...)` arbitrary sub-trees

This page documents every framework-owned key. The schema lives at [`posh.schema.json`](https://raw.githubusercontent.com/foomo/posh/refs/heads/main/posh.schema.json) — point your editor at it for autocomplete:

```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/foomo/posh/refs/heads/main/posh.schema.json
version: v1.0
```

## `prompt`

Consumed by `Plugin.Prompt`.

```yaml
prompt:
  title: "Posh"               # window title and banner text
  prefix: "posh >"            # cursor prefix; framework appends ' › '
  prefixGit: false            # show git branch & tags in prefix
  history:
    limit: 100                # max lines kept in the history file
    filename: .posh/.history
    lockFilename: .posh/.history.lock
  historySearch: true         # enable Ctrl+R fuzzy history search
  aliases:
    k: kubectl
    gco: git checkout
```

| Field | Type | Notes |
| --- | --- | --- |
| `title` | string | Window/banner title. Pass-through to `prompt.OptionTitle`. |
| `prefix` | string | Static prefix. Framework appends ` › `. |
| `prefixGit` | bool | If true, the prefix becomes `prefix (branch  tag) ›`. |
| `history.limit` | int | Lines retained on disk. |
| `history.filename` | string | Relative to `${PROJECT_ROOT}`. |
| `history.lockFilename` | string | Lock file for concurrent shells. |
| `historySearch` | bool | If true, `Ctrl+R` opens fuzzy history search (see [Prompt → History search](/usage/prompt#history-search)). |
| `aliases` | map[string]string | Longest-prefix substitution. Suggested on `<Tab>`. |

## `env`

Name/value pairs prepended to the shell's environment on startup. Variables in values are expanded by the host shell **at config-load time**.

```yaml
env:
  - name: PATH
    value: "${PROJECT_ROOT}/bin:${PATH}"
  - name: KUBECONFIG
    value: "${PROJECT_ROOT}/.posh/.kubeconfig"
  - name: AWS_PROFILE
    value: "myproject-dev"
```

The first entry is what makes ownbrew-installed binaries take precedence over host-installed ones — keep it.

## `ownbrew`

Consumed by `Plugin.Brew` (and exposed as `posh brew`). See [ownbrew docs](https://github.com/foomo/ownbrew) for the full schema; the relevant keys for posh:

```yaml
ownbrew:
  binDir: bin                       # symlinks land here (in $PATH)
  tapDir: .posh/scripts/ownbrew     # local install scripts
  tempDir: .posh/tmp                # download staging
  cellarDir: .posh/bin              # versioned install root
  packages:
    - name: gotsrpc
      tap: foomo/tap/foomo/gotsrpc  # remote package via ownbrew-tap
      version: 2.6.2
    - name: example                 # local package (script in tapDir)
      version: 0.0.0
```

`make shell.build` (or the seeded `make shell` target) runs `bin/posh brew` first to make sure tools are installed before the prompt opens.

## `require`

Consumed by `Plugin.Require` (and exposed as `posh require`). Three groups, each independently optional:

```yaml
require:
  envs:
    - name: VOLTA_HOME
      help: |
        Missing $VOLTA_HOME. Run `volta setup`.

  scripts:
    - name: git
      command: git status >/dev/null 2>&1
      help: This is not a git repo. Please clone the repository.

  packages:
    - name: go
      version: '>=1.26'
      command: go env GOVERSION | cut -c3-
      help: |
        Please install go ≥ %s:
          $ brew install go
```

Behaviour:

- `envs[*]` — fails if the env var is unset
- `scripts[*]` — runs the script; non-zero exit fails with `help`
- `packages[*]` — runs `command`, parses stdout as semver, and compares against `version` (`~`, `>=`, etc.). The `%s` in `help` is replaced with the required version range.

By default, `require.First` short-circuits at the first failure.

## `version`

```yaml
version: v1.0
```

Reserved for future schema migrations. Always set `v1.0` today.

## Custom keys

Anything else under `.posh.yaml` is yours. The seeded scaffold uses this for `welcome.message`:

```yaml
welcome:
  message: Hi, thanks for using POSH!
```

…and reads it from inside the `Welcome` command:

```go
func WelcomeWithConfigKey(v string) WelcomeOption {
    return func(o *Welcome) error {
        return viper.UnmarshalKey(v, &o.cfg)
    }
}
```

This is the recommended pattern: typed config structs in `internal/config/`, populated via `WithConfigKey` options on each command. It keeps the YAML diff-friendly and the Go side type-safe.

## Inspecting the loaded config

```text
posh › config
```

(Equivalent to `bin/posh config` from outside the prompt.) Prints the merged viper view as syntax-highlighted YAML — useful for debugging which env-var overrides won.
