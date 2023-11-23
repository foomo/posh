package tree_test

import (
	"context"
	"testing"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ErrOK = errors.New("ok")

func TestRoot(t *testing.T) {
	l := log.NewTest(t, log.TestWithLevel(log.LevelInfo))
	ctx := context.TODO()

	var (
		ErrRoot    = errors.New("root")
		ErrFirst   = errors.New("first")
		ErrSecond  = errors.New("second")
		ErrSecond1 = errors.New("second-1")
		ErrSecond2 = errors.New("second-2")
		ErrThird   = errors.New("third")
		ErrThird1  = errors.New("third1")
	)

	r := tree.New(&tree.Node{
		Name:        "root",
		Description: "Root tree",
		Execute: func(ctx context.Context, r *readline.Readline) error {
			return ErrRoot
		},
		Nodes: tree.Nodes{
			{
				Name: "first",
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrFirst
				},
			},
			{
				Name: "second",
				Nodes: tree.Nodes{
					{
						Name: "second-1",
						Execute: func(ctx context.Context, r *readline.Readline) error {
							return ErrSecond1
						},
					},
					{
						Name: "second-2",
						Execute: func(ctx context.Context, r *readline.Readline) error {
							return ErrSecond2
						},
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrSecond
				},
			},
			{
				Name: "third",
				Values: func(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
					return []goprompt.Suggest{
						{Text: "third-a"},
						{Text: "third-b"},
					}
				},
				Nodes: tree.Nodes{
					{
						Name: "third-1",
						Execute: func(ctx context.Context, r *readline.Readline) error {
							return ErrThird1
						},
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrThird
				},
			},
		},
	})

	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "tree",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRoot)
			},
		},
		{
			name: "tree first",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrFirst)
			},
		},
		{
			name: "tree first foo",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrFirst)
			},
		},
		{
			name: "tree second",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSecond)
			},
		},
		{
			name: "tree second second-1",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSecond1)
			},
		},
		{
			name: "tree second second-2",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSecond2)
			},
		},
		{
			name: "tree second second-3",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrSecond)
			},
		},
		{
			name: "tree third-a",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrThird)
			},
		},
		{
			name: "tree third-b",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrThird)
			},
		},
		{
			name: "tree third-c",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrRoot)
			},
		},
		{
			name: "tree third-a third-1",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrThird1)
			},
		},
	}

	rl, err := readline.New(l)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			require.NoError(t1, rl.Parse(tt.name))
			if !tt.wantErr(t1, r.Execute(SetT(ctx, t1), rl)) {
				l.Warn(rl.String())
			} else {
				l.Debug(rl.String())
			}
		})
	}
}

func TestRoot_Node(t *testing.T) {
	l := log.NewTest(t, log.TestWithLevel(log.LevelInfo))
	ctx := context.TODO()

	tests := []struct {
		name    string
		root    tree.Root
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "tree",
			root: tree.New(&tree.Node{}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, tree.ErrNoop)
			},
		},
		{
			name: "tree",
			root: tree.New(&tree.Node{
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrOK)
			},
		},
		{
			name: "tree one",
			root: tree.New(&tree.Node{
				Execute: func(ctx context.Context, r *readline.Readline) error {
					assert.Equal(T(ctx), "one", r.Args().At(0))
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrOK)
			},
		},
	}

	rl, err := readline.New(l)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			require.NoError(t1, rl.Parse(tt.name))
			if !tt.wantErr(t1, tt.root.Execute(SetT(ctx, t1), rl)) {
				l.Warn(rl.String())
			} else {
				l.Debug(rl.String())
			}
		})
	}
}

func TestRoot_NodeArgs(t *testing.T) {
	l := log.NewTest(t, log.TestWithLevel(log.LevelInfo))
	ctx := context.TODO()

	tests := []struct {
		name    string
		root    tree.Root
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "tree",
			root: tree.New(&tree.Node{
				Args: tree.Args{
					{
						Name: "first",
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
		},
		{
			name: "tree one",
			root: tree.New(&tree.Node{
				Args: tree.Args{
					{
						Name: "first",
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					assert.Equal(T(ctx), "one", r.Args().At(0))
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrOK)
			},
		},
		{
			name: "tree",
			root: tree.New(&tree.Node{
				Args: tree.Args{
					{
						Name: "first",
					},
					{
						Name: "second",
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
		},
		{
			name: "tree one",
			root: tree.New(&tree.Node{
				Args: tree.Args{
					{
						Name: "first",
					},
					{
						Name: "second",
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
		},
		{
			name: "tree one two",
			root: tree.New(&tree.Node{
				Args: tree.Args{
					{
						Name: "first",
					},
					{
						Name: "second",
					},
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					assert.Equal(T(ctx), "one", r.Args().At(0))
					assert.Equal(T(ctx), "two", r.Args().At(1))
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrOK)
			},
		},
	}

	rl, err := readline.New(l)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			require.NoError(t1, rl.Parse(tt.name))
			if !tt.wantErr(t1, tt.root.Execute(SetT(ctx, t1), rl)) {
				l.Warn(rl.String())
			} else {
				l.Debug(rl.String())
			}
		})
	}
}

func TestRoot_NodeFlags(t *testing.T) {
	l := log.NewTest(t, log.TestWithLevel(log.LevelDebug))
	ctx := context.TODO()

	tests := []struct {
		name    string
		root    tree.Root
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "tree",
			root: tree.New(&tree.Node{
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					fs.Default().String("first", "first", "first")
					fs.Default().Bool("second", false, "second")
					fs.Default().Int64("third", 0, "third")
					return nil
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					if value, err := r.FlagSets().Default().GetString("first"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), "first", value)
					}
					if value, err := r.FlagSets().Default().GetBool("second"); assert.NoError(T(ctx), err) {
						assert.False(T(ctx), value)
					}
					if value, err := r.FlagSets().Default().GetInt64("third"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), int64(0), value)
					}
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrOK)
			},
		},
		{
			name: "tree --first one --second --third 13",
			root: tree.New(&tree.Node{
				Flags: func(ctx context.Context, r *readline.Readline, fs *readline.FlagSets) error {
					fs.Default().String("first", "first", "first")
					fs.Default().Bool("second", false, "second")
					fs.Default().Int64("third", 0, "third")
					return nil
				},
				Execute: func(ctx context.Context, r *readline.Readline) error {
					if value, err := r.FlagSets().Default().GetString("first"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), "one", value)
					}
					if value, err := r.FlagSets().Default().GetBool("second"); assert.NoError(T(ctx), err) {
						assert.True(T(ctx), value)
					}
					if value, err := r.FlagSets().Default().GetInt64("third"); assert.NoError(T(ctx), err) {
						assert.Equal(T(ctx), int64(13), value)
					}
					return ErrOK
				},
			}),
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrOK)
			},
		},
	}

	rl, err := readline.New(l)
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			require.NoError(t1, rl.Parse(tt.name))
			if !tt.wantErr(t1, tt.root.Execute(SetT(ctx, t1), rl)) {
				l.Warn(rl.String())
			} else {
				l.Debug(rl.String())
			}
		})
	}
}

func T(ctx context.Context) *testing.T {
	return ctx.Value("t").(*testing.T) //nolint:forcetypeassert
}

func SetT(ctx context.Context, t *testing.T) context.Context {
	t.Helper()
	return context.WithValue(ctx, "t", t) //nolint:staticcheck
}
