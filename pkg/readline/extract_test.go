package readline_test

import (
	"testing"

	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
)

func TestExtractFlags(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantRest  []string
		wantCalls []string
	}{
		{
			name:     "no args",
			args:     nil,
			wantRest: nil,
		},
		{
			name:     "plain pass-through command",
			args:     []string{"kubectl", "staging", "get", "pods"},
			wantRest: []string{"kubectl", "staging", "get", "pods"},
		},
		{
			name:      "leading flag is consumed",
			args:      []string{"--agent", "kubectl", "staging", "apply"},
			wantRest:  []string{"kubectl", "staging", "apply"},
			wantCalls: []string{"--agent"},
		},
		{
			name:      "flag only",
			args:      []string{"--agent"},
			wantRest:  []string{},
			wantCalls: []string{"--agent"},
		},
		{
			name:     "flags after the command name are not consumed",
			args:     []string{"kubectl", "--agent", "get", "pods"},
			wantRest: []string{"kubectl", "--agent", "get", "pods"},
		},
		{
			name:      "a leading run of several flags is consumed",
			args:      []string{"--agent", "--quiet", "kubectl", "get"},
			wantRest:  []string{"kubectl", "get"},
			wantCalls: []string{"--agent", "--quiet"},
		},
		{
			name:      "scanning stops at the first unknown flag",
			args:      []string{"--agent", "--unknown", "--quiet"},
			wantRest:  []string{"--unknown", "--quiet"},
			wantCalls: []string{"--agent"},
		},
		{
			name:      "a repeated flag applies each time",
			args:      []string{"--agent", "--agent", "welcome"},
			wantRest:  []string{"welcome"},
			wantCalls: []string{"--agent", "--agent"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls []string

			flags := map[string]func(){
				"--agent": func() { calls = append(calls, "--agent") },
				"--quiet": func() { calls = append(calls, "--quiet") },
			}

			assert.Equal(t, tt.wantRest, readline.ExtractFlags(tt.args, flags))
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestExtractFlags_NoFlags(t *testing.T) {
	args := []string{"--agent", "welcome"}

	// With nothing to match, every token passes through untouched.
	assert.Equal(t, args, readline.ExtractFlags(args, nil))
}
