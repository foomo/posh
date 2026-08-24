package command

import (
	"context"

	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/check"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/util"
)

type (
	Check struct {
		l        log.Logger
		tree     tree.Root
		check    check.Check
		checkers check.Checkers
	}
	Option func(*Check)
)

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewCheck(l log.Logger, checkers ...check.Checker) *Check {
	inst := &Check{
		l: l,
		// AgentCheck emits the same results as JSON; DefaultCheck renders a table.
		check:    util.Pick(agent.IsAgentMode(), check.AgentCheck, check.DefaultCheck),
		checkers: checkers,
	}
	inst.tree = tree.New(&tree.Node{
		Name:        "check",
		Description: "Print all system checks",
		Execute:     inst.run,
	})

	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Check) Name() string {
	return c.tree.Node().Name
}

func (c *Check) Description() string {
	return c.tree.Node().Description
}

// Describe implements the optional Describer interface, letting
// `posh agent catalog` describe this command's subtree.
func (c *Check) Describe(ctx context.Context) CommandInfo {
	return c.tree.Describe(ctx)
}

// Skill implements the optional Skiller interface. The command takes no
// arguments, so its structure says nothing - what an agent needs is that this
// is the diagnostic to reach for first, and that the checks are project
// defined.
func (c *Check) Skill(ctx context.Context, name string) string {
	return "#### Notes\n\n" +
		"Runs this project's environment checks - the same ones the interactive\n" +
		"shell prints on startup - reporting each as `success`, `failure`,\n" +
		"`warning` or `note`. Which checks exist is defined by this project, so\n" +
		"they cover whatever it actually depends on: required binaries,\n" +
		"credentials, cluster reachability.\n\n" +
		"Run this first when another command fails for an environment-shaped\n" +
		"reason: a missing tool, an expired login, an unreachable service. It is\n" +
		"read-only and changes nothing, so it is always safe to run.\n\n" +
		"A failing check is a report, not an error: the command still exits 0. Read\n" +
		"the results rather than relying on the exit code. Fixing a failure usually\n" +
		"means acting outside posh - installing the tool, re-authenticating - and\n" +
		"posh cannot do that for you."
}

func (c *Check) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.tree.Complete(ctx, r)
}

func (c *Check) Execute(ctx context.Context, r *readline.Readline) error {
	return c.tree.Execute(ctx, r)
}

func (c *Check) Help(ctx context.Context, r *readline.Readline) string {
	return c.tree.Help(ctx, r)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Check) run(ctx context.Context, r *readline.Readline) error {
	return c.check(ctx, c.l, c.checkers)
}
