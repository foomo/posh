package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/foomo/fender/fend"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/ownbrew"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/foomo/posh/pkg/prompt"
	"github.com/foomo/posh/pkg/prompt/history"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/validate"
)

type Plugin struct {
	l        log.Logger
	commands command.Commands
}

func New(l log.Logger) (plugin.Plugin, error) { //nolint:unparam
	inst := &Plugin{
		l:        l,
		commands: command.Commands{},
	}

	// add commands
	inst.commands.Add(
		command.NewExit(l),
		command.NewHelp(l, inst.commands),
	)

	return inst, nil
}

func (p *Plugin) Prompt(ctx context.Context, cfg config.Prompt) error {
	sh, err := prompt.New(p.l,
		prompt.WithTitle(cfg.Title),
		prompt.WithPrefix(cfg.Prefix),
		prompt.WithContext(ctx),
		prompt.WithCommands(p.commands),
		prompt.WithFileHistory(
			history.FileWithLimit(cfg.History.Limit),
			history.FileWithFilename(cfg.History.Filename),
			history.FileWithLockFilename(cfg.History.LockFilename),
		),
	)
	if err != nil {
		return err
	}
	return sh.Run()
}

func (p *Plugin) Execute(ctx context.Context, args []string) error {
	r, err := readline.New(p.l)
	if err != nil {
		return err
	}

	if err := r.Parse(strings.Join(args, " ")); err != nil {
		return err
	}

	if c := p.commands.Get(r.Cmd()); c == nil {
		return fmt.Errorf("invalid [cmd] argument: %s", r.Cmd())
	} else if err := c.Execute(ctx, r); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) Dependencies(ctx context.Context, cfg config.Dependencies) error {
	var fends []fend.Fend
	fends = append(fends, validate.DependenciesEnvs(p.l, cfg.Envs)...)
	fends = append(fends, validate.DependenciesScripts(ctx, p.l, cfg.Scripts)...)
	fends = append(fends, validate.DependenciesPackages(ctx, p.l, cfg.Packages)...)
	if fendErr, err := fend.First(fends...); err != nil {
		return err
	} else if fendErr != nil {
		return fendErr
	}
	return nil
}

func (p *Plugin) Packages(ctx context.Context, cfg []config.Package) error {
	brew, err := ownbrew.New(p.l,
		ownbrew.WithPackages(cfg...),
	)
	if err != nil {
		return err
	}
	return brew.Install(ctx)
}
