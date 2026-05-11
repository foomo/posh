package tree_test

import (
	"context"
	"testing"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
)

func TestRoot_Help(t *testing.T) {
	tests := []struct {
		name        string
		root        tree.Root
		input       string
		contains    []string
		notContains []string
		equals      string
	}{
		{
			name:   "nil node",
			root:   tree.New(nil),
			input:  "anything",
			equals: "command not found",
		},
		{
			name:   "node with no children and no args, no input -> not found",
			root:   tree.New(&tree.Node{Name: "root", Description: "desc"}),
			input:  "root",
			equals: "command not found",
		},
		{
			name: "LenIs(1) renders root help",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Nodes: tree.Nodes{
					{Name: "alpha", Description: "alpha desc"},
					{Name: "bravo", Description: "bravo desc"},
				},
			}),
			input:    "root self",
			contains: []string{"root desc", "Available Commands:", "alpha", "alpha desc", "bravo", "bravo desc", "Usage:", "root [command]"},
		},
		{
			name: "leaf with required arg renders [name] bracket",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Args:        tree.Args{{Name: "x", Description: "x desc"}},
			}),
			input:       "root self",
			contains:    []string{"Usage:", "root", "[x]", "Arguments:", "x desc"},
			notContains: []string{"<x>"},
		},
		{
			name: "leaf with optional arg renders <name> bracket",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Args:        tree.Args{{Name: "x", Description: "x desc", Optional: true}},
			}),
			input:       "root self",
			contains:    []string{"Usage:", "root", "<x>", "Arguments:", "x desc"},
			notContains: []string{"[x]"},
		},
		{
			name: "leaf with flags renders Flags section",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					fs.Default().String("foo", "", "foo desc")
					return nil
				},
			}),
			input:    "root self",
			contains: []string{"Flags:", "--foo", "foo desc"},
		},
		{
			name: "leaf with flags callback error omits Flags section",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					return errBoom
				},
			}),
			input:       "root self",
			contains:    []string{"Usage:", "root"},
			notContains: []string{"Flags:"},
		},
		{
			name: "descends to matched child",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Nodes: tree.Nodes{
					{Name: "child", Description: "child desc"},
				},
			}),
			input:       "root self child",
			contains:    []string{"child desc", "Usage:", "child"},
			notContains: []string{"root desc"},
		},
		{
			name: "unknown child falls back to root",
			root: tree.New(&tree.Node{
				Name:        "root",
				Description: "root desc",
				Nodes: tree.Nodes{
					{Name: "child", Description: "child desc"},
				},
			}),
			input:       "root self unknown",
			contains:    []string{"root desc", "Available Commands:", "child"},
			notContains: []string{"child desc\n\nUsage:"},
		},
		{
			name:   "more than one arg with no children -> not found",
			root:   tree.New(&tree.Node{Name: "root", Description: "root desc"}),
			input:  "root self extra",
			equals: "command not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newReadline(t, tt.input)

			got := tt.root.Help(t.Context(), rl)
			if tt.equals != "" {
				assert.Equal(t, tt.equals, got)
				return
			}

			for _, s := range tt.contains {
				assert.Contains(t, got, s, "expected substring %q in help output", s)
			}

			for _, s := range tt.notContains {
				assert.NotContains(t, got, s, "did not expect substring %q in help output", s)
			}
		})
	}
}
