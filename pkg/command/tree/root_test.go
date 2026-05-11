package tree_test

import (
	"context"
	"testing"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	errOK   = errors.New("ok")
	errBoom = errors.New("boom")
)

// ------------------------------------------------------------------------------------------------
// ~ New / Node accessor
// ------------------------------------------------------------------------------------------------

func TestNew(t *testing.T) {
	n := &tree.Node{Name: "root"}
	r := tree.New(n)
	assert.Same(t, n, r.Node(), "Node() should return the node passed to New")

	assert.NotPanics(t, func() {
		_ = tree.New(nil).Node()
	})
}

// ------------------------------------------------------------------------------------------------
// ~ Execute
// ------------------------------------------------------------------------------------------------

func TestRoot_Execute_TreeTraversal(t *testing.T) {
	var (
		errRoot    = errors.New("root")
		errFirst   = errors.New("first")
		errSecond  = errors.New("second")
		errSecond1 = errors.New("second-1")
		errSecond2 = errors.New("second-2")
		errThird   = errors.New("third")
		errThird1  = errors.New("third-1")
	)

	r := tree.New(&tree.Node{
		Name:        "root",
		Description: "Root tree",
		Execute:     func(ctx context.Context, r *readline.Readline) error { return errRoot },
		Nodes: tree.Nodes{
			{
				Name:    "first",
				Execute: func(ctx context.Context, r *readline.Readline) error { return errFirst },
			},
			{
				Name: "second",
				Nodes: tree.Nodes{
					{Name: "second-1", Execute: func(ctx context.Context, r *readline.Readline) error { return errSecond1 }},
					{Name: "second-2", Execute: func(ctx context.Context, r *readline.Readline) error { return errSecond2 }},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errSecond },
			},
			{
				Name: "third",
				Values: func(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
					return []goprompt.Suggest{{Text: "third-a"}, {Text: "third-b"}}
				},
				Nodes: tree.Nodes{
					{Name: "third-1", Execute: func(ctx context.Context, r *readline.Readline) error { return errThird1 }},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errThird },
			},
		},
	})

	tests := []struct {
		input string
		want  error
	}{
		{"tree", errRoot},
		{"tree first", errFirst},
		{"tree first unknown-extra", errFirst},
		{"tree second", errSecond},
		{"tree second second-1", errSecond1},
		{"tree second second-2", errSecond2},
		{"tree second second-3", errSecond}, // unknown sub falls back to second
		{"tree third-a", errThird},          // dynamic Values() match
		{"tree third-b", errThird},
		{"tree third-c", errRoot}, // dynamic miss falls back to root
		{"tree third-a third-1", errThird1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			rl := newReadline(t, tt.input)
			assert.ErrorIs(t, r.Execute(t.Context(), rl), tt.want)
		})
	}
}

func TestRoot_Execute_Node(t *testing.T) {
	tests := []struct {
		name    string
		root    tree.Root
		input   string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil node returns ErrNoop",
			root:    tree.New(nil),
			input:   "anything",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, tree.ErrNoop) },
		},
		{
			name:    "empty node returns ErrNoop",
			root:    tree.New(&tree.Node{}),
			input:   "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, tree.ErrNoop) },
		},
		{
			name: "node with execute returns its error",
			root: tree.New(&tree.Node{
				Execute: func(ctx context.Context, r *readline.Readline) error { return errOK },
			}),
			input:   "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
		{
			name: "execute sees parsed args",
			root: tree.New(&tree.Node{
				Execute: func(ctx context.Context, r *readline.Readline) error {
					assert.Equal(T(ctx), "one", r.Args().At(0))
					return errOK
				},
			}),
			input:   "tree one",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
		{
			name: "node with children but no input returns ErrMissingCommand",
			root: tree.New(&tree.Node{
				Nodes: tree.Nodes{
					{Name: "child", Execute: func(ctx context.Context, r *readline.Readline) error { return errOK }},
				},
			}),
			input: "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingCommand)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newReadline(t, tt.input)
			tt.wantErr(t, tt.root.Execute(SetT(t.Context(), t), rl))
		})
	}
}

func TestRoot_Execute_NodeArgs(t *testing.T) {
	tests := []struct {
		name    string
		root    tree.Root
		input   string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "required missing",
			root: tree.New(&tree.Node{
				Args:    tree.Args{{Name: "first"}},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errOK },
			}),
			input: "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
		},
		{
			name: "required supplied",
			root: tree.New(&tree.Node{
				Args: tree.Args{{Name: "first"}},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					assert.Equal(T(ctx), "one", r.Args().At(0))
					return errOK
				},
			}),
			input:   "tree one",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
		{
			name: "two required, none supplied",
			root: tree.New(&tree.Node{
				Args:    tree.Args{{Name: "first"}, {Name: "second"}},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errOK },
			}),
			input: "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
		},
		{
			name: "two required, one supplied",
			root: tree.New(&tree.Node{
				Args:    tree.Args{{Name: "first"}, {Name: "second"}},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errOK },
			}),
			input: "tree one",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
		},
		{
			name: "two required, both supplied",
			root: tree.New(&tree.Node{
				Args: tree.Args{{Name: "first"}, {Name: "second"}},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					assert.Equal(T(ctx), "one", r.Args().At(0))
					assert.Equal(T(ctx), "two", r.Args().At(1))

					return errOK
				},
			}),
			input:   "tree one two",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
		{
			name: "optional missing is ok",
			root: tree.New(&tree.Node{
				Args:    tree.Args{{Name: "first", Optional: true}},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errOK },
			}),
			input:   "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newReadline(t, tt.input)
			tt.wantErr(t, tt.root.Execute(SetT(t.Context(), t), rl))
		})
	}
}

func TestRoot_Execute_NodeFlags(t *testing.T) {
	tests := []struct {
		name    string
		root    tree.Root
		input   string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "defaults when no flags supplied",
			root: tree.New(&tree.Node{
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					fs.Default().String("first", "first", "first")
					fs.Default().Bool("second", false, "second")
					fs.Default().Int64("third", 0, "third")

					return nil
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					if v, err := r.FlagSets().Default().GetString("first"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), "first", v)
					}

					if v, err := r.FlagSets().Default().GetBool("second"); assert.NoError(T(ctx), err) {
						assert.False(T(ctx), v)
					}

					if v, err := r.FlagSets().Default().GetInt64("third"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), int64(0), v)
					}

					return errOK
				},
			}),
			input:   "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
		{
			name: "supplied flags parsed",
			root: tree.New(&tree.Node{
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					fs.Default().String("first", "first", "first")
					fs.Default().Bool("second", false, "second")
					fs.Default().Int64("third", 0, "third")

					return nil
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					if v, err := r.FlagSets().Default().GetString("first"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), "one", v)
					}

					if v, err := r.FlagSets().Default().GetBool("second"); assert.NoError(T(ctx), err) {
						assert.True(T(ctx), v)
					}

					if v, err := r.FlagSets().Default().GetInt64("third"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), int64(13), v)
					}

					return errOK
				},
			}),
			input:   "tree --first one --second --third 13",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errOK) },
		},
		{
			name: "flags callback error propagates",
			root: tree.New(&tree.Node{
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					return errBoom
				},
				Execute: func(ctx context.Context, r *readline.Readline) error { return errOK },
			}),
			input:   "tree",
			wantErr: func(t assert.TestingT, err error, _ ...any) bool { return assert.ErrorIs(t, err, errBoom) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newReadline(t, tt.input)
			tt.wantErr(t, tt.root.Execute(SetT(t.Context(), t), rl))
		})
	}
}
