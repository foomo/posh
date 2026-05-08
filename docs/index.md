---
layout: home

hero:
  name: posh
  text: Project Oriented Shell
  tagline: An interactive, isolated and hackable Makefile.
  image:
    src: /logo.png
    alt: posh
  actions:
    - theme: brand
      text: Get Started
      link: /guide/introduction
    - theme: alt
      text: Quick Start
      link: /guide/quick-start
    - theme: alt
      text: View on GitHub
      link: https://github.com/foomo/posh

features:
  - icon: 🐚
    title: Interactive REPL
    details: A scoped, project-aware shell with completion, history, aliases and a built-in command registry — not a thin wrapper around `bash`.
  - icon: 🧱
    title: Hackable in Go
    details: Every shell is a small Go program you own. Add commands by implementing an interface; compose with completion, validation and graceful shutdown as needed.
  - icon: 📦
    title: Reproducible by default
    details: Pin tool versions with ownbrew, declare environment variables and prerequisite checks in `.posh.yaml`, and ship the same shell to every contributor.
  - icon: 🧩
    title: Plays well with Make
    details: posh is started from `make shell` and executes shell commands as a fallback. Keep your Makefile; build the interactive surface on top.
---
