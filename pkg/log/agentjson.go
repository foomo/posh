package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// AgentJSON is a Logger implementation that writes one JSON object per line
// (JSONL) to stdout instead of PTerm's human-formatted, colored output. It is
// selected in agent mode (see pkg/agent) so an AI coding agent driving
// `posh execute` gets machine-parseable output instead of prose.
//
// Every entry carries a fixed "type"/"schema_version" discriminator so
// consumers can dispatch on a marker rather than guessing from shape.
type (
	AgentJSON struct {
		name  string
		level Level
		w     io.Writer
		out   *json.Encoder
	}
	AgentJSONOption func(*AgentJSON)
)

type agentLogEntry struct {
	Type          string `json:"type"`
	SchemaVersion string `json:"schema_version"`
	Level         string `json:"level"`
	Name          string `json:"name,omitempty"`
	Message       string `json:"message"`
}

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func AgentJSONWithWriter(w io.Writer) AgentJSONOption {
	return func(o *AgentJSON) {
		o.w = w
	}
}

func AgentJSONWithLevel(v Level) AgentJSONOption {
	return func(o *AgentJSON) {
		o.level = v
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewAgentJSON(opts ...AgentJSONOption) *AgentJSON {
	inst := &AgentJSON{
		level: LevelInfo,
		w:     os.Stdout,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(inst)
		}
	}

	enc := json.NewEncoder(inst.w)
	enc.SetEscapeHTML(false)
	inst.out = enc

	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (l *AgentJSON) Level() Level {
	return l.level
}

func (l *AgentJSON) IsLevel(v Level) bool {
	return l.level <= v
}

func (l *AgentJSON) Named(name string) Logger {
	clone := *l
	clone.name = name

	return &clone
}

func (l *AgentJSON) Print(a ...any) {
	l.write("info", sprintSpaced(a))
}

func (l *AgentJSON) Printf(format string, a ...any) {
	l.write("info", fmt.Sprintf(format, a...))
}

func (l *AgentJSON) Success(a ...any) {
	l.write("success", sprintSpaced(a))
}

func (l *AgentJSON) Successf(format string, a ...any) {
	l.write("success", fmt.Sprintf(format, a...))
}

func (l *AgentJSON) Trace(a ...any) {
	if l.IsLevel(LevelTrace) {
		l.write("trace", sprintSpaced(a))
	}
}

func (l *AgentJSON) Tracef(format string, a ...any) {
	if l.IsLevel(LevelTrace) {
		l.write("trace", fmt.Sprintf(format, a...))
	}
}

func (l *AgentJSON) Debug(a ...any) {
	if l.IsLevel(LevelDebug) {
		l.write("debug", sprintSpaced(a))
	}
}

func (l *AgentJSON) Debugf(format string, a ...any) {
	if l.IsLevel(LevelDebug) {
		l.write("debug", fmt.Sprintf(format, a...))
	}
}

func (l *AgentJSON) Info(a ...any) {
	if l.IsLevel(LevelInfo) {
		l.write("info", sprintSpaced(a))
	}
}

func (l *AgentJSON) Infof(format string, a ...any) {
	if l.IsLevel(LevelInfo) {
		l.write("info", fmt.Sprintf(format, a...))
	}
}

func (l *AgentJSON) Warn(a ...any) {
	if l.IsLevel(LevelWarn) {
		l.write("warn", sprintSpaced(a))
	}
}

func (l *AgentJSON) Warnf(format string, a ...any) {
	if l.IsLevel(LevelWarn) {
		l.write("warn", fmt.Sprintf(format, a...))
	}
}

func (l *AgentJSON) Error(a ...any) {
	if l.IsLevel(LevelError) {
		l.write("error", sprintSpaced(a))
	}
}

func (l *AgentJSON) Errorf(format string, a ...any) {
	if l.IsLevel(LevelError) {
		l.write("error", fmt.Sprintf(format, a...))
	}
}

func (l *AgentJSON) Fatal(a ...any) {
	l.write("fatal", sprintSpaced(a))
	os.Exit(1)
}

func (l *AgentJSON) Fatalf(format string, a ...any) {
	l.write("fatal", fmt.Sprintf(format, a...))
	os.Exit(1)
}

func (l *AgentJSON) Must(err error) {
	if err != nil {
		l.Fatal(err.Error())
	}
}

func (l *AgentJSON) SlogHandler() slog.Handler {
	return slog.NewJSONHandler(l.w, nil)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

// sprintSpaced joins a like fmt.Println does (always space-separated),
// matching PTerm's Print/Info/etc. behavior (which delegates to
// pterm.Println), unlike fmt.Sprint which only spaces non-string operands.
func sprintSpaced(a []any) string {
	return strings.TrimSuffix(fmt.Sprintln(a...), "\n")
}

func (l *AgentJSON) write(level, message string) {
	if err := l.out.Encode(agentLogEntry{
		Type:          "posh.log",
		SchemaVersion: "1",
		Level:         level,
		Name:          l.name,
		Message:       message,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
