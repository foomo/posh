package log

import (
	"fmt"
	"log/slog"

	ptermx "github.com/foomo/posh/pkg/pterm"
	"github.com/pterm/pterm"
)

type (
	PTerm struct {
		name  string
		level Level
	}
	PTermOption func(*PTerm)
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func PTermWithDisableColor(v bool) PTermOption {
	return func(o *PTerm) {
		if v {
			pterm.DisableColor()
		}
	}
}

func PTermWithLevel(v Level) PTermOption {
	return func(o *PTerm) {
		o.level = v
		switch {
		case v <= LevelTrace:
			pterm.Debug.ShowLineNumber = true
			fallthrough
		case v <= LevelDebug:
			pterm.EnableDebugMessages()
			fallthrough
		default:
			pterm.Debug.LineNumberOffset = 1
		}
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewPTerm(opts ...PTermOption) *PTerm {
	inst := &PTerm{
		level: LevelError,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(inst)
		}
	}

	pterm.Info.Prefix.Text = "⎈"
	pterm.Debug.Prefix.Text = "⛏︎"
	pterm.Fatal.Prefix.Text = "⛔︎"
	pterm.Error.Prefix.Text = "⛌"
	pterm.Warning.Prefix.Text = "⚠"
	pterm.Success.Prefix.Text = "✓"
	pterm.DefaultLogger.ShowTime = false

	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (l *PTerm) Level() Level {
	return l.level
}

func (l *PTerm) IsLevel(v Level) bool {
	return l.level <= v
}

func (l *PTerm) Named(name string) Logger {
	clone := *l
	clone.name = name
	return &clone
}

func (l *PTerm) Print(a ...interface{}) {
	pterm.Println(l.prefix(a)...)
}

func (l *PTerm) Printf(format string, a ...interface{}) {
	l.Print(fmt.Sprintf(format, a...))
}

func (l *PTerm) Success(a ...interface{}) {
	pterm.Success.Println(l.prefix(a)...)
}

func (l *PTerm) Successf(format string, a ...interface{}) {
	l.Success(fmt.Sprintf(format, a...))
}

func (l *PTerm) Trace(a ...interface{}) {
	if l.IsLevel(LevelTrace) {
		pterm.Debug.Println(l.prefix(a)...)
	}
}

func (l *PTerm) Tracef(format string, a ...interface{}) {
	l.Trace(fmt.Sprintf(format, a...))
}

func (l *PTerm) Debug(a ...interface{}) {
	if l.IsLevel(LevelDebug) {
		pterm.Debug.Println(l.prefix(a)...)
	}
}

func (l *PTerm) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *PTerm) Info(a ...interface{}) {
	if l.IsLevel(LevelInfo) {
		pterm.Info.Println(l.prefix(a)...)
	}
}

func (l *PTerm) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *PTerm) Warn(a ...interface{}) {
	if l.IsLevel(LevelWarn) {
		pterm.Warning.Println(l.prefix(a)...)
	}
}

func (l *PTerm) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *PTerm) Error(a ...interface{}) {
	if l.IsLevel(LevelError) {
		pterm.Error.Println(l.prefix(a)...)
	}
}

func (l *PTerm) Errorf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *PTerm) Fatal(a ...interface{}) {
	pterm.Fatal.Println(l.prefix(a)...)
}

func (l *PTerm) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l *PTerm) Must(err error) {
	if err != nil {
		l.Fatal(err.Error())
	}
}

func (l *PTerm) SlogHandler() slog.Handler {
	return ptermx.NewSlogHandler()
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (l *PTerm) prefix(a []any) []any {
	var ret []any
	if l.name != "" && l.IsLevel(LevelDebug) {
		ret = append(ret, fmt.Sprintf("[%s]", l.name))
	}
	return append(ret, a...)
}
