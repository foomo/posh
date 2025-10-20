package log

import (
	"fmt"
	"log/slog"
	"os"
)

type (
	Fmt struct {
		name  string
		level Level
	}
	FmtOption func(*Fmt)
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func FmtWithLevel(v Level) FmtOption {
	return func(o *Fmt) {
		o.level = v
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewFmt(opts ...FmtOption) *Fmt {
	inst := &Fmt{
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

func (l *Fmt) Level() Level {
	return l.level
}

func (l *Fmt) IsLevel(v Level) bool {
	return l.level <= v
}

func (l *Fmt) Named(name string) Logger {
	clone := *l
	clone.name = name

	return &clone
}

func (l *Fmt) Print(a ...interface{}) {
	fmt.Println(l.prefix("", a)...)
}

func (l *Fmt) Printf(format string, a ...interface{}) {
	l.Print(fmt.Sprintf(format, a...))
}

func (l *Fmt) Success(a ...interface{}) {
	fmt.Println(l.prefix("success", a)...)
}

func (l *Fmt) Successf(format string, a ...interface{}) {
	l.Success(fmt.Sprintf(format, a...))
}

func (l *Fmt) Trace(a ...interface{}) {
	if l.IsLevel(LevelTrace) {
		fmt.Println(l.prefix("trace", a)...)
	}
}

func (l *Fmt) Tracef(format string, a ...interface{}) {
	l.Trace(fmt.Sprintf(format, a...))
}

func (l *Fmt) Debug(a ...interface{}) {
	if l.IsLevel(LevelDebug) {
		fmt.Println(l.prefix("debug", a)...)
	}
}

func (l *Fmt) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *Fmt) Info(a ...interface{}) {
	if l.IsLevel(LevelInfo) {
		fmt.Println(l.prefix("info", a)...)
	}
}

func (l *Fmt) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Fmt) Warn(a ...interface{}) {
	if l.IsLevel(LevelWarn) {
		fmt.Println(l.prefix("warn", a)...)
	}
}

func (l *Fmt) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *Fmt) Error(a ...interface{}) {
	if l.IsLevel(LevelError) {
		fmt.Println(l.prefix("error", a)...)
	}
}

func (l *Fmt) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *Fmt) Fatal(a ...interface{}) {
	fmt.Println(l.prefix("fatal", a)...)
	os.Exit(1)
}

func (l *Fmt) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}

func (l *Fmt) Must(err error) {
	if err != nil {
		l.Fatal(err.Error())
	}
}

func (l *Fmt) SlogHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, nil)
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (l *Fmt) prefix(level string, a []any) []any {
	var ret []interface{}
	if level != "" {
		ret = append(ret, level+":")
	}

	if l.name != "" && l.IsLevel(LevelDebug) {
		ret = append(ret, fmt.Sprintf("[%s]", l.name))
	}

	return append(ret, a...)
}
