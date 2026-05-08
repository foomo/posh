---
title: CLI Reference
---

# CLI Reference

This section is auto-generated from the cobra command tree by `cmd/docgen`. Regenerate it with:

```shell
make docs.cli
```

The output is committed to git, so the docs site builds without Go installed. CI verifies the generated tree matches the current command definitions.

## Two flavours

| Command | Available in `posh` (global) | Available in `bin/posh` (project) |
| --- | --- | --- |
| [`posh init`](./posh_init) | ✅ | — |
| [`posh config`](./posh_config) | ✅ | ✅ |
| [`posh version`](./posh_version) | ✅ | ✅ |
| [`posh prompt`](./posh_prompt) | — | ✅ |
| [`posh execute`](./posh_execute) | — | ✅ |
| [`posh require`](./posh_require) | — | ✅ |
| [`posh brew`](./posh_brew) | — | ✅ |

The split exists because `prompt`, `execute`, `require` and `brew` delegate to the [`Plugin`](/plugin/overview) compiled into the project binary. The global `posh` has no plugin to delegate to.

## All commands

- [posh](./posh)
- [posh init](./posh_init)
- [posh config](./posh_config)
- [posh version](./posh_version)
- [posh prompt](./posh_prompt)
- [posh execute](./posh_execute)
- [posh require](./posh_require)
- [posh brew](./posh_brew)
