---
title: CLI Reference
---

# CLI Reference

This section is auto-generated from the cobra command tree by `cmd/docgen`. Regenerate it with:

```shell
make docs.cli
```

The output is committed to git, so the docs site builds without Go installed. Re-run `make docs.cli` and commit the result whenever you add, remove or reflag a cobra command.

## Two flavours

| Command | Available in `posh` (global) | Available in `bin/posh` (project) |
| --- | --- | --- |
| [`posh init`](./posh_init) | ✅ | — |
| [`posh config`](./posh_config) | ✅ | ✅ |
| [`posh version`](./posh_version) | ✅ | ✅ |
| [`posh agent`](./posh_agent) | — | ✅ |
| [`posh prompt`](./posh_prompt) | — | ✅ |
| [`posh execute`](./posh_execute) | — | ✅ |
| [`posh require`](./posh_require) | — | ✅ |
| [`posh brew`](./posh_brew) | — | ✅ |

The split exists because `agent`, `prompt`, `execute`, `require` and `brew` delegate to the [`Plugin`](/plugin/overview) compiled into the project binary. The global `posh` has no plugin to delegate to.

## All commands

- [posh](./posh)
- [posh init](./posh_init)
- [posh config](./posh_config)
- [posh version](./posh_version)
- [posh agent](./posh_agent)
  - [posh agent catalog](./posh_agent_catalog)
  - [posh agent skill](./posh_agent_skill)
    - [posh agent skill get](./posh_agent_skill_get)
    - [posh agent skill install](./posh_agent_skill_install)
    - [posh agent skill update](./posh_agent_skill_update)
    - [posh agent skill uninstall](./posh_agent_skill_uninstall)
- [posh prompt](./posh_prompt)
- [posh execute](./posh_execute)
- [posh require](./posh_require)
- [posh brew](./posh_brew)
