# Theming

The interactive prompt can be colored with one of the four
[Catppuccin](https://catppuccin.com) palettes. It is **opt-in**: without the
`prompt.theme` key, posh uses the same colors it always has.

```yaml
prompt:
  theme: mocha
```

| Flavor | Background |
| --- | --- |
| `mocha` | dark |
| `macchiato` | dark |
| `frappe` | dark |
| `latte` | light |

There is no flag and no environment variable — the flavor is a project
decision, so it lives in `.posh.yaml` alongside the rest of the prompt
configuration. An unknown name fails at startup rather than falling back
silently, so a typo surfaces immediately.

## Scope

Theming affects the **prompt only**: the cursor prefix, the input line and the
completion dropdown (suggestions, descriptions and the scrollbar). Log output,
tables, trees, the startup banner and pre-flight checks are unchanged.

The dropdown is the part worth setting: go-prompt's stock colors are white on
cyan with black-on-turquoise descriptions, which is hard to read on a light
background. `theme: latte` is the fix.

## Fidelity

The prompt is built on a fork of `c-bata/go-prompt`, whose color type carries
only the 16 ANSI slots — it has no 256-color or truecolor path, so it cannot
render Catppuccin hex values. Colors are therefore mapped to the nearest ANSI
slot.

In practice this usually still looks right: if you run a Catppuccin *terminal*
theme, slots 0–15 already hold Catppuccin values. When the terminal reports no
color support at all, the theme yields the terminal's defaults.

## For plugin authors

The scaffolded plugin passes the key through:

```go
prompt.New(p.l,
    prompt.WithTitle(cfg.Title),
    prompt.WithTheme(cfg.Theme),
    // ...
)
```

An existing shell that predates this option keeps its current colors until
`prompt.WithTheme(cfg.Theme)` is added to its own `plugin.go`.

Options passed via `prompt.WithPromptOptions(...)` are applied after the theme,
so they still win.

The palette is also readable directly, keyed by semantic role rather than color
name:

```go
import "github.com/foomo/posh/pkg/theme"

t, err := theme.Resolve("latte")   // nil when the string is empty
if t != nil {
    hex := t.Hex(theme.RoleWarning)      // "#df8e1d"
    col := t.Prompt(theme.RoleWarning)   // nearest go-prompt color
}
```

Roles: `RoleText`, `RoleMuted`, `RolePrimary`, `RoleSecondary`, `RoleSuccess`,
`RoleWarning`, `RoleError`, `RoleInfo`, `RoleDebug`.

Custom palettes are not supported — the four Catppuccin flavors are the whole
set.
