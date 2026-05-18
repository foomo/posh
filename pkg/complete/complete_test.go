package complete_test

import (
	"context"
	"testing"

	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/complete"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
)

type testCmd struct {
	name        string
	description string
	args        []goprompt.Suggest
	flags       []goprompt.Suggest
	addl        []goprompt.Suggest
	generic     []goprompt.Suggest
}

func (c *testCmd) Name() string        { return c.name }
func (c *testCmd) Description() string { return c.description }
func (c *testCmd) Execute(ctx context.Context, r *readline.Readline) error {
	return nil
}

type withArgs struct{ *testCmd }

func (c *withArgs) CompleteArguments(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.args
}

type withFlags struct{ *testCmd }

func (c *withFlags) CompleteFlags(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.flags
}

type withAdditional struct{ *testCmd }

func (c *withAdditional) CompleteAdditionalArgs(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.addl
}

type withGeneric struct{ *testCmd }

func (c *withGeneric) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.generic
}

func newLogger(t *testing.T) log.Logger {
	t.Helper()
	return log.NewTest(t, log.TestWithLevel(log.LevelInfo))
}

func TestSuggest_RootListsCommands(t *testing.T) {
	cmds := command.Commands{}
	cmds.Add(&testCmd{name: "alpha", description: "first"})
	cmds.Add(&testCmd{name: "beta", description: "second"})

	got := complete.Suggest(t.Context(), newLogger(t), cmds, nil, "")

	assert.Equal(t, []string{"alpha\tfirst", "beta\tsecond"}, got)
}

func TestSuggest_RootPrefixFilter(t *testing.T) {
	cmds := command.Commands{}
	cmds.Add(&testCmd{name: "alpha", description: "first"})
	cmds.Add(&testCmd{name: "alphabet", description: "longer"})
	cmds.Add(&testCmd{name: "beta", description: "second"})

	got := complete.Suggest(t.Context(), newLogger(t), cmds, nil, "alph")

	assert.Equal(t, []string{"alpha\tfirst", "alphabet\tlonger"}, got)
}

func TestSuggest_RootOmitsDescriptionWhenEmpty(t *testing.T) {
	cmds := command.Commands{}
	cmds.Add(&testCmd{name: "bare"})

	got := complete.Suggest(t.Context(), newLogger(t), cmds, nil, "")

	assert.Equal(t, []string{"bare"}, got)
}

func TestSuggest_ArgumentCompleterDispatch(t *testing.T) {
	cmd := &withArgs{testCmd: &testCmd{
		name: "foo",
		args: []goprompt.Suggest{{Text: "one", Description: "1st"}, {Text: "two"}},
	}}
	cmds := command.Commands{}
	cmds.Add(cmd)

	got := complete.Suggest(t.Context(), newLogger(t), cmds, []string{"foo"}, "")

	assert.Equal(t, []string{"one\t1st", "two"}, got)
}

func TestSuggest_ArgumentCompleterPrefixFilter(t *testing.T) {
	cmd := &withArgs{testCmd: &testCmd{
		name: "foo",
		args: []goprompt.Suggest{{Text: "one"}, {Text: "onyx"}, {Text: "two"}},
	}}
	cmds := command.Commands{}
	cmds.Add(cmd)

	got := complete.Suggest(t.Context(), newLogger(t), cmds, []string{"foo"}, "on")

	assert.Equal(t, []string{"one", "onyx"}, got)
}

func TestSuggest_GenericCompleterFallback(t *testing.T) {
	cmd := &withGeneric{testCmd: &testCmd{
		name:    "foo",
		generic: []goprompt.Suggest{{Text: "g1"}},
	}}
	cmds := command.Commands{}
	cmds.Add(cmd)

	got := complete.Suggest(t.Context(), newLogger(t), cmds, []string{"foo"}, "")

	assert.Equal(t, []string{"g1"}, got)
}

func TestSuggest_FlagCompleterDispatch(t *testing.T) {
	cmd := &withFlags{testCmd: &testCmd{
		name:  "foo",
		flags: []goprompt.Suggest{{Text: "--bar"}, {Text: "--baz"}},
	}}
	cmds := command.Commands{}
	cmds.Add(cmd)

	got := complete.Suggest(t.Context(), newLogger(t), cmds, []string{"foo"}, "--")

	assert.Equal(t, []string{"--bar", "--baz"}, got)
}

func TestSuggest_UnknownCommandReturnsEmpty(t *testing.T) {
	cmds := command.Commands{}
	cmds.Add(&testCmd{name: "foo"})

	got := complete.Suggest(t.Context(), newLogger(t), cmds, []string{"nope"}, "")

	assert.Empty(t, got)
}
