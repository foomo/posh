package tree_test

import (
	"context"
	"sort"
	"testing"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescribe_NilNode(t *testing.T) {
	assert.Empty(t, tree.New(nil).Describe(t.Context()))
}

func TestDescribe_Nested(t *testing.T) {
	actual := tree.New(&tree.Node{
		Name:        "cache",
		Description: "Manage the internal cache",
		Nodes: tree.Nodes{
			{
				Name:        "clear",
				Description: "Clear caches",
				Args: tree.Args{
					{
						Name:        "Namespace",
						Description: "Name of namespace to clear.",
						Optional:    true,
						Repeat:      true,
					},
				},
			},
			{
				Name:        "list",
				Description: "List all caches",
			},
		},
	}).Describe(t.Context())

	assert.Equal(t, "cache", actual.FullPath)
	assert.Equal(t, "Manage the internal cache", actual.Description)

	require.Len(t, actual.Subcommands, 2)

	clearCmd := actual.Subcommands[0]
	assert.Equal(t, "cache clear", clearCmd.FullPath)
	assert.Equal(t, "Clear caches", clearCmd.Description)

	require.Len(t, clearCmd.Arguments, 1)
	assert.Equal(t, tree.ArgInfo{
		Name:        "Namespace",
		Description: "Name of namespace to clear.",
		Optional:    true,
		Repeat:      true,
	}, clearCmd.Arguments[0])
	assert.Equal(t, "<Namespace>...", clearCmd.Usage())

	assert.Equal(t, "cache list", actual.Subcommands[1].FullPath)
	assert.Empty(t, actual.Subcommands[1].Arguments)
	assert.Empty(t, actual.Subcommands[1].Usage())
}

func TestDescribe_NestedDepth(t *testing.T) {
	actual := tree.New(&tree.Node{
		Name: "a",
		Nodes: tree.Nodes{
			{
				Name:  "b",
				Nodes: tree.Nodes{{Name: "c"}},
			},
		},
	}).Describe(t.Context())

	require.Len(t, actual.Subcommands, 1)
	require.Len(t, actual.Subcommands[0].Subcommands, 1)
	assert.Equal(t, "a b c", actual.Subcommands[0].Subcommands[0].FullPath)
}

func TestDescribe_Args(t *testing.T) {
	actual := tree.New(&tree.Node{
		Name: "run",
		Args: tree.Args{
			{Name: "Cluster", Description: "the target cluster"},
			{Name: "Target", Optional: true},
		},
	}).Describe(t.Context())

	// Order is significant: arguments are positional.
	require.Len(t, actual.Arguments, 2)
	assert.Equal(t, "Cluster", actual.Arguments[0].Name)
	assert.Equal(t, "the target cluster", actual.Arguments[0].Description)
	assert.False(t, actual.Arguments[0].Optional)
	assert.False(t, actual.Arguments[0].Repeat)

	assert.Equal(t, "Target", actual.Arguments[1].Name)
	assert.Empty(t, actual.Arguments[1].Description)
	assert.True(t, actual.Arguments[1].Optional)

	assert.Equal(t, "[Cluster] <Target>", actual.Usage())
}

func TestCommandInfo_Usage(t *testing.T) {
	tests := []struct {
		name     string
		args     []tree.ArgInfo
		expected string
	}{
		{
			name:     "none",
			expected: "",
		},
		{
			name:     "required",
			args:     []tree.ArgInfo{{Name: "Key"}},
			expected: "[Key]",
		},
		{
			name:     "optional repeat",
			args:     []tree.ArgInfo{{Name: "Namespace", Optional: true, Repeat: true}},
			expected: "<Namespace>...",
		},
		{
			name:     "mixed",
			args:     []tree.ArgInfo{{Name: "Key"}, {Name: "Value", Optional: true}},
			expected: "[Key] <Value>",
		},
		{
			name:     "required repeat",
			args:     []tree.ArgInfo{{Name: "Path", Repeat: true}},
			expected: "[Path]...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tree.CommandInfo{Arguments: tt.args}.Usage())
		})
	}
}

// TestDescribe_FlagsSorted pins the ordering guarantee callers rely on:
// flags come out sorted by name regardless of which set they were registered
// in or in what order. This holds because FlagSets.All merges everything into
// one pflag.FlagSet whose VisitAll sorts - the test exists so that stops being
// an implementation detail nobody notices changing.
func TestDescribe_FlagsSorted(t *testing.T) {
	names := []string{"zeta", "yankee", "xray", "whiskey", "victor", "uniform", "tango", "sierra"}

	actual := tree.New(&tree.Node{
		Name: "run",
		Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
			// Registered in reverse alphabetical order, alternating sets.
			for i, name := range names {
				if i%2 == 0 {
					fs.Default().String(name, "", name+" flag")
				} else {
					fs.Internal().String(name, "", name+" flag")
				}
			}

			return nil
		},
	}).Describe(t.Context())

	require.Len(t, actual.Flags, len(names))
	assert.True(t, sort.SliceIsSorted(actual.Flags, func(i, j int) bool {
		return actual.Flags[i].Name < actual.Flags[j].Name
	}), "flags must be sorted by name, got %v", actual.Flags)
	assert.Equal(t, "sierra", actual.Flags[0].Name)
	assert.Equal(t, "zeta", actual.Flags[len(names)-1].Name)
}

func TestDescribe_FlagFields(t *testing.T) {
	actual := tree.New(&tree.Node{
		Name: "run",
		Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
			fs.Default().StringP("mid", "m", "fallback", "mid flag")
			fs.Default().Bool("dry-run", false, "dry run flag")

			return nil
		},
	}).Describe(t.Context())

	require.Len(t, actual.Flags, 2)

	assert.Equal(t, "dry-run", actual.Flags[0].Name)
	assert.Equal(t, "bool", actual.Flags[0].Type)
	assert.Equal(t, "false", actual.Flags[0].Default)
	assert.Empty(t, actual.Flags[0].Shorthand)

	assert.Equal(t, "mid", actual.Flags[1].Name)
	assert.Equal(t, "m", actual.Flags[1].Shorthand)
	assert.Equal(t, "string", actual.Flags[1].Type)
	assert.Equal(t, "fallback", actual.Flags[1].Default)
	assert.Equal(t, "mid flag", actual.Flags[1].Description)
}

func TestDescribe_FlagsError(t *testing.T) {
	actual := tree.New(&tree.Node{
		Name:        "run",
		Description: "run it",
		Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
			return errors.New("nope")
		},
	}).Describe(t.Context())

	assert.Equal(t, "run", actual.FullPath)
	assert.Equal(t, "run it", actual.Description)
	assert.Empty(t, actual.Flags)
}

// TestDescribe_Dynamic asserts the walk never invokes a Values callback: doing
// so would resolve live state (kubeconfig, cloud APIs) at catalog time.
func TestDescribe_Dynamic(t *testing.T) {
	actual := tree.New(&tree.Node{
		Name: "kubectl",
		Nodes: tree.Nodes{
			{
				Name:        "cluster",
				Description: "target cluster",
				Values: func(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
					panic("Values must not be invoked while building the catalog")
				},
			},
		},
	}).Describe(t.Context())

	require.Len(t, actual.Subcommands, 1)
	assert.Equal(t, "kubectl <cluster>", actual.Subcommands[0].FullPath)
	assert.True(t, actual.Subcommands[0].Dynamic)
	assert.False(t, actual.Dynamic)
}
