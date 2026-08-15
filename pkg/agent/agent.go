// Package agent detects whether posh is being driven by an AI coding agent
// rather than a human, so callers (e.g. cmd/execute.go) can pick agent-friendly
// defaults such as JSON output instead of human-formatted prose.
package agent

import (
	"os"
	"strings"
)

// harnessEnvVars lists environment variables set by known agent harnesses.
// Presence of any of these (non-empty) is treated as a signal that posh is
// running under an agent, absent an explicit POSH_AGENT_MODE override.
var harnessEnvVars = []string{
	"CLAUDECODE",
	"CURSOR_AGENT",
	"GITHUB_COPILOT",
	"AMAZON_Q",
	"OPENCODE",
}

var (
	detected bool
	flagSet  bool
)

func init() {
	detected = DetectFrom(os.Getenv)
}

// DetectFrom resolves agent mode using the given environment lookup. An
// explicit POSH_AGENT_MODE wins over harness detection in both directions.
// Exported so detection logic is testable without mutating the real process
// environment.
func DetectFrom(getenv func(string) string) bool {
	switch strings.ToLower(getenv("POSH_AGENT_MODE")) {
	case "0", "false", "no":
		return false
	case "1", "true", "yes":
		return true
	}

	for _, name := range harnessEnvVars {
		if getenv(name) != "" {
			return true
		}
	}

	return false
}

// SetFlag records an explicit --agent flag, e.g. for harnesses that set none
// of the recognized environment variables.
func SetFlag(v bool) {
	flagSet = v
}

// IsAgentMode reports whether posh should behave as if driven by an agent.
func IsAgentMode() bool {
	return detected || flagSet
}

// SetDetected overrides the environment-based detection performed in init and
// returns a func restoring the previous value.
//
// This exists for tests: inside an agent harness the real environment sets
// CLAUDECODE, so IsAgentMode reports true for the whole test binary and the
// human-formatted output path would be unreachable. Production code should let
// init do the detection and use POSH_AGENT_MODE to override it.
func SetDetected(v bool) func() {
	prev := detected
	detected = v

	return func() { detected = prev }
}
