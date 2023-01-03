package onepassword

import (
	"context"

	argparse2 "github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/shell"
)

type Command struct {
	l      log.Logger
	cfg    Config
	parser *argparse2.Root
	client *OnePassword
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewCommand(l log.Logger, cfg Config, client *OnePassword) *Command {
	inst := &Command{
		l:      l,
		cfg:    cfg,
		client: client,
	}

	inst.parser = &argparse2.Root{
		Name: "op",
		Nodes: []*argparse2.Node{
			{
				Name:        "get",
				Description: "retrieve item",
				Args: []*argparse2.Arg{
					{
						Name: "id",
					},
				},
				Execute: inst.get,
			},
			{
				Name:        "signin",
				Description: "sign into your account",
				Execute:     inst.signin,
			},
			{
				Name:        "register",
				Description: "register an account",
				Args: []*argparse2.Arg{
					{
						Name: "email",
					},
				},
				Execute: inst.register,
			},
		},
	}
	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Command) Name() string {
	return c.parser.Name
}

func (c *Command) Description() string {
	return "run go mod"
}

func (c *Command) Complete(ctx context.Context, r *readline.Readline, d prompt.Document) (suggests []prompt.Suggest) {
	return c.parser.Complete(ctx, r)
}

func (c *Command) Execute(ctx context.Context, r *readline.Readline) error {
	return c.parser.Execute(ctx, r)
}

func (c *Command) Help() string {
	return `1Password session helper.

Usage:
  op [command]

Available commands:
  get [id]          Retrieve an entry from your account
  signin            Sign into your 1Password account for the session
  register [email]  Add your 1Password account
`
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Command) get(ctx context.Context, r *readline.Readline) error {
	return shell.New(ctx, c.l,
		"op",
		"--account", c.cfg.Account,
		"item", "get", r.Args().At(1),
		"--format", "json",
	).
		Args(r.AdditionalArgs()...).
		Run()
}

func (c *Command) register(ctx context.Context, r *readline.Readline) error {
	return shell.New(ctx, c.l,
		"op", "account", "add",
		"--address", c.cfg.Account+".1password.eu",
		"--email", r.Args().At(1),
	).
		Args(r.AdditionalArgs()...).
		Wait()
}

func (c *Command) signin(ctx context.Context, r *readline.Readline) error {
	if ok, _ := c.client.Session(c.cfg.Account); ok {
		c.l.Info("Already signed in")
		return nil
	} else if err := c.client.SignIn(ctx, c.cfg.Account); err != nil {
		return err
	}
	return nil
}
