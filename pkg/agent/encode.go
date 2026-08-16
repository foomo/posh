package agent

import (
	"encoding/json"
	"os"
)

// Encode writes v to stdout as a single JSON value followed by a newline.
//
// Commands use this in agent mode to emit one machine-readable value instead of
// the human-formatted PTerm tables and trees they render otherwise. Results are
// emitted bare: a caller invoked a specific command and already knows what it
// asked for. Log lines written by pkg/log.AgentJSON keep a "type" field, which
// is what lets a consumer tell an interleaved log line from the result value.
//
// Encode does not itself check agent mode - callers branch on IsAgentMode, and
// `posh agent catalog` uses it unconditionally.
//
// HTML escaping is disabled to match pkg/log.AgentJSON: environment values and
// cache contents routinely contain characters like & and < that would otherwise
// be mangled into escape sequences.
func Encode(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	return enc.Encode(v)
}
