package agent_test

import (
	"testing"

	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/readline"
	"github.com/stretchr/testify/assert"
)

// TestFlag_SetsAgentMode covers the wiring the posh CLI relies on: a leading
// agent.Flag scanned off argv turns agent mode on. The scan itself is tested in
// pkg/readline.
func TestFlag_SetsAgentMode(t *testing.T) {
	defer agent.SetFlag(false)

	restore := agent.SetDetected(false)
	defer restore()

	agent.SetFlag(false)

	rest := readline.ExtractFlags([]string{agent.Flag, "welcome"}, map[string]func(){
		agent.Flag: func() { agent.SetFlag(true) },
	})

	assert.Equal(t, []string{"welcome"}, rest)
	assert.True(t, agent.IsAgentMode())
}
