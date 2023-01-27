package log

import (
	"fmt"
	"testing"
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

func (l *Test) Print(a ...interface{}) {
	l.t.Log(l.prefix("", a)...)
}

func (l *Test) Printf(format string, a ...interface{}) {
	l.Print(fmt.Sprintf(format, a...))
}

func (l *Test) Success(a ...interface{}) {
	l.t.Log(l.prefix("success", a)...)
}

func (l *Test) Successf(format string, a ...interface{}) {
	l.Success(fmt.Sprintf(format, a...))
}

func (l *Test) Trace(a ...interface{}) {
	if l.IsLevel(LevelTrace) {
		l.t.Log(l.prefix("trace", a)...)
	}
}

func (l *Test) Tracef(format string, a ...interface{}) {
	l.Trace(fmt.Sprintf(format, a...))
}

func (l *Test) Debug(a ...interface{}) {
	if l.IsLevel(LevelDebug) {
		l.t.Log(l.prefix("debug", a)...)
	}
}

func (l *Test) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *Test) Info(a ...interface{}) {
	if l.IsLevel(LevelInfo) {
		l.t.Log(l.prefix("info", a)...)
	}
}

func (l *Test) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Test) Warn(a ...interface{}) {
	if l.IsLevel(LevelWarn) {
		l.t.Log(l.prefix("warn", a)...)
	}
}

func (l *Test) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *Test) Error(a ...interface{}) {
	if l.IsLevel(LevelError) {
		l.t.Error(l.prefix("error", a)...)
	}
}

func (l *Test) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *Test) Fatal(a ...interface{}) {
	l.t.Fatal(l.prefix("fatal", a)...)
}

func (l *Test) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l *Test) Must(err error) {
	if err != nil {
		l.Fatal(err.Error())
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (l *Test) prefix(level string, a []any) []any {
	var ret []interface{}
	if level != "" {
		ret = append(ret, level+":")
	}
	if l.name != "" && l.IsLevel(LevelDebug) {
		ret = append(ret, fmt.Sprintf("[%s]", l.name))
	}
	return append(ret, a...)
}
