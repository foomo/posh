package tree_test

import (
	"context"
	"testing"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// staticSuggest returns a fixed list of suggestions for use in tests.
func staticSuggest(texts ...string) func(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return func(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
		out := make([]goprompt.Suggest, len(texts))
		for i, t := range texts {
			out[i] = goprompt.Suggest{Text: t}
		}

		return out
	}
}

func argSuggest(texts ...string) func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
	return func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
		out := make([]goprompt.Suggest, len(texts))
		for i, x := range texts {
			out[i] = goprompt.Suggest{Text: x}
		}

		return out
	}
}

func suggestTexts(s []goprompt.Suggest) []string {
	out := make([]string, len(s))
	for i, v := range s {
		out[i] = v.Text
	}

	return out
}

// ------------------------------------------------------------------------------------------------
// ~ ModeArgs
// ------------------------------------------------------------------------------------------------

func TestRoot_Complete_ModeArgs(t *testing.T) {
	tests := []struct {
		name  string
		root  tree.Root
		input string
		want  []string
	}{
		{
			name: "root lists static children sorted",
			root: tree.New(&tree.Node{
				Nodes: tree.Nodes{
					{Name: "bravo", Description: "b"},
					{Name: "alpha", Description: "a"},
				},
			}),
			input: "tree",
			want:  []string{"alpha", "bravo"},
		},
		{
			name: "root lists dynamic Values from one child",
			root: tree.New(&tree.Node{
				Nodes: tree.Nodes{
					{Name: "dynamic", Values: staticSuggest("zeta", "alpha")},
				},
			}),
			input: "tree",
			want:  []string{"alpha", "zeta"},
		},
		{
			name: "root no children uses Args[0].Suggest",
			root: tree.New(&tree.Node{
				Args: tree.Args{{Name: "x", Suggest: argSuggest("one", "two")}},
			}),
			input: "tree",
			want:  []string{"one", "two"},
		},
		{
			name: "nested arg suggestion after matched child",
			root: tree.New(&tree.Node{
				Nodes: tree.Nodes{
					{
						Name: "first",
						Args: tree.Args{{Name: "x", Suggest: argSuggest("kappa", "iota")}},
					},
				},
			}),
			input: "tree first iota",
			want:  []string{"iota", "kappa"},
		},
		{
			name: "repeat arg uses Last().Suggest",
			root: tree.New(&tree.Node{
				Args: tree.Args{{Name: "x", Repeat: true, Suggest: argSuggest("rho", "pi")}},
			}),
			input: "tree a b",
			want:  []string{"pi", "rho"},
		},
		{
			name: "flags callback error returns nil",
			root: tree.New(&tree.Node{
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					return errBoom
				},
				Args: tree.Args{{Name: "x", Suggest: argSuggest("never")}},
			}),
			input: "tree",
			want:  nil,
		},
		{
			name:  "no match no args returns empty",
			root:  tree.New(&tree.Node{}),
			input: "tree x",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newReadline(t, tt.input)
			require.Equal(t, readline.ModeArgs, rl.Mode())

			got := tt.root.Complete(t.Context(), rl)
			if len(tt.want) == 0 {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.want, suggestTexts(got))
			}
		})
	}
}

// ------------------------------------------------------------------------------------------------
// ~ ModeFlags
// ------------------------------------------------------------------------------------------------

func TestRoot_Complete_ModeFlags(t *testing.T) {
	t.Run("lists registered flags sorted", func(t *testing.T) {
		r := tree.New(&tree.Node{
			Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
				fs.Default().String("foo", "", "foo desc")
				fs.Default().Bool("bar", false, "bar desc")

				return nil
			},
		})

		rl := newReadline(t, "tree --")
		require.Equal(t, readline.ModeFlags, rl.Mode())

		got := r.Complete(t.Context(), rl)
		assert.Equal(t, []string{"--bar", "--foo"}, suggestTexts(got))
	})

	t.Run("flag value completion via GetValues", func(t *testing.T) {
		r := tree.New(&tree.Node{
			Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
				fs.Default().String("foo", "", "foo desc")
				return fs.Default().SetValues("foo", "beta", "alpha")
			},
		})

		rl := newReadline(t, "tree --foo something")
		require.Equal(t, readline.ModeFlags, rl.Mode())

		got := r.Complete(t.Context(), rl)
		// Complete sorts the final slice by Text.
		assert.Equal(t, []string{"alpha", "beta"}, suggestTexts(got))
	})

	t.Run("flags callback error returns nil", func(t *testing.T) {
		r := tree.New(&tree.Node{
			Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
				return errBoom
			},
		})

		rl := newReadline(t, "tree --")
		require.Equal(t, readline.ModeFlags, rl.Mode())

		assert.Nil(t, r.Complete(t.Context(), rl))
	})

	t.Run("matched child flags are listed", func(t *testing.T) {
		r := tree.New(&tree.Node{
			Nodes: tree.Nodes{
				{
					Name: "first",
					Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
						fs.Default().String("only-on-first", "", "")
						return nil
					},
				},
			},
		})

		rl := newReadline(t, "tree first --")
		require.Equal(t, readline.ModeFlags, rl.Mode())

		got := r.Complete(t.Context(), rl)
		assert.Equal(t, []string{"--only-on-first"}, suggestTexts(got))
	})
}

// ------------------------------------------------------------------------------------------------
// ~ ModeAdditionalArgs
// ------------------------------------------------------------------------------------------------

func TestRoot_Complete_ModeAdditionalArgs(t *testing.T) {
	r := tree.New(&tree.Node{
		Args: tree.Args{{Name: "x", Suggest: argSuggest("should-not-show")}},
	})

	rl := newReadline(t, "tree -- additional")
	require.Equal(t, readline.ModeAdditionalArgs, rl.Mode())

	assert.Nil(t, r.Complete(t.Context(), rl))
}
