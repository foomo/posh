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

func (l *Fmt) Print(a ...any) {
	fmt.Println(l.prefix("", a)...)
}

func (l *Fmt) Printf(format string, a ...any) {
	l.Print(fmt.Sprintf(format, a...))
}

func (l *Fmt) Success(a ...any) {
	fmt.Println(l.prefix("success", a)...)
}

func (l *Fmt) Successf(format string, a ...any) {
	l.Success(fmt.Sprintf(format, a...))
}

func (l *Fmt) Trace(a ...any) {
	if l.IsLevel(LevelTrace) {
		fmt.Println(l.prefix("trace", a)...)
	}
}

func (l *Fmt) Tracef(format string, a ...any) {
	l.Trace(fmt.Sprintf(format, a...))
}

func (l *Fmt) Debug(a ...any) {
	if l.IsLevel(LevelDebug) {
		fmt.Println(l.prefix("debug", a)...)
	}
}

func (l *Fmt) Debugf(format string, a ...any) {
	l.Debug(fmt.Sprintf(format, a...))
}

func (l *Fmt) Info(a ...any) {
	if l.IsLevel(LevelInfo) {
		fmt.Println(l.prefix("info", a)...)
	}
}

func (l *Fmt) Infof(format string, a ...any) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Fmt) Warn(a ...any) {
	if l.IsLevel(LevelWarn) {
		fmt.Println(l.prefix("warn", a)...)
	}
}

func (l *Fmt) Warnf(format string, a ...any) {
	l.Warn(fmt.Sprintf(format, a...))
}

func (l *Fmt) Error(a ...any) {
	if l.IsLevel(LevelError) {
		fmt.Println(l.prefix("error", a)...)
	}
}

func (l *Fmt) Errorf(format string, a ...any) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *Fmt) Fatal(a ...any) {
	fmt.Println(l.prefix("fatal", a)...)
	os.Exit(1)
}

func (l *Fmt) Fatalf(format string, a ...any) {
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
	var ret []any
	if level != "" {
		ret = append(ret, level+":")
	}

	if l.name != "" && l.IsLevel(LevelDebug) {
		ret = append(ret, fmt.Sprintf("[%s]", l.name))
	}

	return append(ret, a...)
}
