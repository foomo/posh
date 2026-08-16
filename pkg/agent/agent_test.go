package agent_test

import (
	"testing"

	"github.com/foomo/posh/pkg/agent"
	"github.com/stretchr/testify/assert"
)

func TestDetectFrom(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{name: "no signal", env: map[string]string{}, want: false},
		{name: "harness env var set", env: map[string]string{"CLAUDECODE": "1"}, want: true},
		{name: "unrelated env var set", env: map[string]string{"HOME": "/root"}, want: false},
		{name: "explicit override true", env: map[string]string{"POSH_AGENT_MODE": "true"}, want: true},
		{
			name: "explicit override false wins over harness var",
			env:  map[string]string{"POSH_AGENT_MODE": "0", "CLAUDECODE": "1"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getenv := func(key string) string { return tt.env[key] }

			assert.Equal(t, tt.want, agent.DetectFrom(getenv))
		})
	}
}

func TestIsAgentMode(t *testing.T) {
	t.Run("explicit flag forces true", func(t *testing.T) {
		agent.SetFlag(true)
		defer agent.SetFlag(false)

		assert.True(t, agent.IsAgentMode())
	})
}
