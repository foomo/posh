package readline_test

import (
	"testing"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
)

func TestReadline(t *testing.T) {
	tests := []struct {
		name string
		want func(t *testing.T, r *readline.Readline)
	}{
		{
			name: "foo bar",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 [bar]
Flags:                []
AdditionalArgs:       []
`,
					r.String(),
				)
			},
		},
		{
			name: "foo bar baz",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 [bar baz]
Flags:                []
AdditionalArgs:       []
`,
					r.String(),
				)
			},
		},
		{
			name: "foo bar --baz",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 [bar]
Flags:                [--baz]
AdditionalArgs:       []
`,
					r.String(),
				)
			},
		},
		{
			name: "foo --baz bar",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 []
Flags:                [--baz bar]
AdditionalArgs:       []
`,
					r.String(),
				)
			},
		},
		{
			name: "foo --baz bar1",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 []
Flags:                [--baz bar1]
AdditionalArgs:       []
`,
					r.String(),
				)
			},
		},
		{
			name: "foo | cat",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 []
Flags:                []
AdditionalArgs:       [| cat]
`,
					r.String(),
				)
			},
		},
		{
			name: "foo --bar1 --bar2 one --bar3 two,three,four",
			want: func(t *testing.T, r *readline.Readline) {
				t.Helper()
				assert.Equal(t, `
Cmd:                  foo
Args:                 []
Flags:                [--bar1 --bar2 one --bar3 two,three,four]
AdditionalArgs:       []
`,
					r.String(),
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := log.NewTest(t, log.TestWithLevel(log.LevelDebug))
			if r, err := readline.New(l); assert.NoError(t, err) {
				assert.NoError(t, r.Parse(tt.name))
				tt.want(t, r)
			}
		})
	}
}
