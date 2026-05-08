package log

import (
	"log/slog"
)

type Logger interface {
	Level() Level
	IsLevel(level Level) bool
	Named(name string) Logger
	Print(a ...any)
	Printf(format string, a ...any)
	Success(a ...any)
	Successf(format string, a ...any)
	Trace(a ...any)
	Tracef(format string, a ...any)
	Debug(a ...any)
	Debugf(format string, a ...any)
	Info(a ...any)
	Infof(format string, a ...any)
	Warn(a ...any)
	Warnf(format string, a ...any)
	Error(a ...any)
	Errorf(format string, a ...any)
	Fatal(a ...any)
	Fatalf(format string, a ...any)
	Must(err error)
	SlogHandler() slog.Handler
}
