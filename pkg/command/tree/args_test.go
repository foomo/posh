package tree_test

import (
	"testing"

	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
)

func TestArgs_Last(t *testing.T) {
	first := &tree.Arg{Name: "first"}
	second := &tree.Arg{Name: "second"}

	tests := []struct {
		name string
		args tree.Args
		want *tree.Arg
	}{
		{name: "empty", args: tree.Args{}, want: nil},
		{name: "nil", args: nil, want: nil},
		{name: "one", args: tree.Args{first}, want: first},
		{name: "many", args: tree.Args{first, second}, want: second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Same(t, tt.want, tt.args.Last())
		})
	}
}

func TestArgs_Validate(t *testing.T) {
	tests := []struct {
		name    string
		args    tree.Args
		input   readline.Args
		wantErr assert.ErrorAssertionFunc
		errMsg  string
	}{
		{
			name:    "no args no input",
			args:    tree.Args{},
			input:   readline.Args{},
			wantErr: assert.NoError,
		},
		{
			name:  "required missing",
			args:  tree.Args{{Name: "first"}},
			input: readline.Args{},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
			errMsg: "first",
		},
		{
			name:    "required supplied",
			args:    tree.Args{{Name: "first"}},
			input:   readline.Args{"value"},
			wantErr: assert.NoError,
		},
		{
			name:    "optional missing",
			args:    tree.Args{{Name: "first", Optional: true}},
			input:   readline.Args{},
			wantErr: assert.NoError,
		},
		{
			name:  "second required missing",
			args:  tree.Args{{Name: "first"}, {Name: "second"}},
			input: readline.Args{"v1"},
			wantErr: func(t assert.TestingT, err error, _ ...any) bool {
				return assert.ErrorIs(t, err, tree.ErrMissingArgument)
			},
			errMsg: "second",
		},
		{
			name:    "required ok optional missing",
			args:    tree.Args{{Name: "first"}, {Name: "second", Optional: true}},
			input:   readline.Args{"v1"},
			wantErr: assert.NoError,
		},
		{
			name:    "all supplied",
			args:    tree.Args{{Name: "first"}, {Name: "second"}},
			input:   readline.Args{"v1", "v2"},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.Validate(tt.input)
			tt.wantErr(t, err)

			if tt.errMsg != "" && err != nil {
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}
