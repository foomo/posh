package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/foomo/fender/fend"
	command2 "github.com/foomo/posh/example/posh/command"
	"github.com/foomo/posh/integration/onepassword"
	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/ownbrew"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/foomo/posh/pkg/prompt"
	"github.com/foomo/posh/pkg/prompt/check"
	"github.com/foomo/posh/pkg/prompt/history"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/validate"
	"github.com/spf13/viper"
)

type Plugin struct {
	l        log.Logger
	c        cache.Cache
	commands command.Commands
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(l log.Logger) (plugin.Plugin, error) { //nolint:unparam
	inst := &Plugin{
		l:        l,
		c:        cache.MemoryCache{},
		commands: command.Commands{},
	}

	var (
		err               error
		onePasswordCfg    onepassword.Config
		onePasswordClient *onepassword.OnePassword
	)

	// load configurations
	l.Must(viper.UnmarshalKey("onePassword", &onePasswordCfg))

	// create dependency instances
	onePasswordClient, err = onepassword.New(l, inst.c, onepassword.WithTokenFilename(onePasswordCfg.TokenFilename))
	l.Must(err)

	// add commands
	inst.commands.Add(
		command2.NewGoMod(l, inst.c),
		command2.NewGoGenerate(l, inst.c),
		onepassword.NewCommand(l, onePasswordCfg, onePasswordClient),
		command.NewCache(l, inst.c),
		command.NewExit(l),
		command.NewHelp(l, inst.commands),
	)
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (p *Plugin) Dependencies(ctx context.Context, cfg config.Dependencies) error {
	var fends []fend.Fend
	fends = append(fends, validate.DependenciesEnvs(p.l, cfg.Envs)...)
	fends = append(fends, validate.DependenciesScripts(ctx, p.l, cfg.Scripts)...)
	fends = append(fends, validate.DependenciesPackages(ctx, p.l, cfg.Packages)...)
	fends = append(fends,
		validate.GitUser(ctx, p.l, validate.GitUserName, validate.GitUserEmail(`(.*)@(bestbytes\.com)`)),
	)
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

func (p *Plugin) Prompt(ctx context.Context, cfg config.Prompt) error {
	sh, err := prompt.New(p.l,
		prompt.WithTitle(cfg.Title),
		prompt.WithPrefix(cfg.Prefix),
		prompt.WithContext(ctx),
		prompt.WithCommands(p.commands),
		prompt.WithCheckers(
			func(ctx context.Context, l log.Logger) check.Info {
				return check.Info{
					Name:   "one",
					Note:   "all good",
					Status: check.StatusSuccess,
				}
			},
			func(ctx context.Context, l log.Logger) check.Info {
				return check.Info{
					Name:   "two",
					Note:   "please take some action",
					Status: check.StatusFailure,
				}
			},
		),
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
