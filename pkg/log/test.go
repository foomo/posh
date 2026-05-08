package log

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/neilotoole/slogt"
)

type (
	Test struct {
		t     *testing.T
		name  string
		level Level
	}
	TestOption func(*Test)
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func TestWithLevel(v Level) TestOption {
	return func(o *Test) {
		o.level = v
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewTest(t *testing.T, opts ...TestOption) *Test {
	t.Helper()

	inst := &Test{
		t:     t,
		level: LevelError,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(inst)
		}
	}

	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (l *Test) Level() Level {
	return l.level
}

func (l *Test) IsLevel(v Level) bool {
	return l.level <= v
}

func (l *Test) Named(name string) Logger {
	clone := *l
	clone.name = name

	return &clone
}

func (l *Test) Print(a ...any) {
	l.t.Log(l.prefix("", a)...)
}

func (l *Test) Printf(format string, a ...any) {
	l.Print(fmt.Sprintf(format, a...))
}

func (l *Test) Success(a ...any) {
	l.t.Log(l.prefix("success", a)...)
}

func (l *Test) Successf(format string, a ...any) {
	l.Success(fmt.Sprintf(format, a...))
}

func (l *Test) Trace(a ...any) {
	if l.IsLevel(LevelTrace) {
		l.t.Log(l.prefix("trace", a)...)
	}
}

func (l *Test) Tracef(format string, a ...any) {
	l.Trace(fmt.Sprintf(format, a...))
}

func (l *Test) Debug(a ...any) {
	if l.IsLevel(LevelDebug) {
		l.t.Log(l.prefix("debug", a)...)
	}
}

func (l *Test) Debugf(format string, a ...any) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *Test) Info(a ...any) {
	if l.IsLevel(LevelInfo) {
		l.t.Log(l.prefix("info", a)...)
	}
}

func (l *Test) Infof(format string, a ...any) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Test) Warn(a ...any) {
	if l.IsLevel(LevelWarn) {
		l.t.Log(l.prefix("warn", a)...)
	}
}

func (l *Test) Warnf(format string, a ...any) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *Test) Error(a ...any) {
	if l.IsLevel(LevelError) {
		l.t.Error(l.prefix("error", a)...)
	}
}

func (l *Test) Errorf(format string, a ...any) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *Test) Fatal(a ...any) {
	l.t.Fatal(l.prefix("fatal", a)...)
}

func (l *Test) Fatalf(format string, a ...any) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l *Test) Must(err error) {
	if err != nil {
		l.Fatal(err.Error())
	}
}

func (l *Test) SlogHandler() slog.Handler {
	return slogt.New(l.t).Handler()
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (l *Test) prefix(level string, a []any) []any {
	var ret []any
	if level != "" {
		ret = append(ret, level+":")
	}

	if l.name != "" && l.IsLevel(LevelDebug) {
		ret = append(ret, fmt.Sprintf("[%s]", l.name))
	}

	return append(ret, a...)
}
