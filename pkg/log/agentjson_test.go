package log_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/foomo/posh/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentJSON_Info(t *testing.T) {
	var buf bytes.Buffer

	l := log.NewAgentJSON(log.AgentJSONWithWriter(&buf))

	l.Info("hello", "world")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 1)

	var entry map[string]any
	require.NoError(t, json.Unmarshal([]byte(lines[0]), &entry))

	assert.Equal(t, "posh.log", entry["type"])
	assert.Equal(t, "1", entry["schema_version"])
	assert.Equal(t, "info", entry["level"])
	assert.Equal(t, "hello world", entry["message"])
}

func TestAgentJSON_LevelGating(t *testing.T) {
	var buf bytes.Buffer

	l := log.NewAgentJSON(log.AgentJSONWithWriter(&buf), log.AgentJSONWithLevel(log.LevelWarn))

	l.Info("suppressed")
	l.Warn("kept")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	require.Len(t, lines, 1)
	assert.Contains(t, lines[0], "kept")
}

func TestAgentJSON_EachLineIsValidJSON(t *testing.T) {
	var buf bytes.Buffer

	l := log.NewAgentJSON(log.AgentJSONWithWriter(&buf))

	l.Info("first")
	l.Error("second")

	for line := range strings.SplitSeq(strings.TrimSpace(buf.String()), "\n") {
		var entry map[string]any
		assert.NoError(t, json.Unmarshal([]byte(line), &entry))
	}
}
