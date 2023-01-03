package log

import (
	"fmt"
	"os"
)

type (
	Fmt struct {
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

func (l *Fmt) Print(a ...interface{}) {
	fmt.Println(append([]interface{}{"Info:"}, a...)...)
}

func (l *Fmt) Printf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func (l *Fmt) Debug(a ...interface{}) {
	if l.level <= LevelDebug {
		fmt.Println(append([]interface{}{"Debug:"}, a...)...)
	}
}

func (l *Fmt) Debugf(format string, a ...interface{}) {
	if l.level <= LevelDebug {
		fmt.Printf("Debug: "+format+"\n", a...)
	}
}

func (l *Fmt) Info(a ...interface{}) {
	if l.level <= LevelInfo {
		fmt.Println(append([]interface{}{"Info:"}, a...)...)
	}
}

func (l *Fmt) Infof(format string, a ...interface{}) {
	if l.level <= LevelInfo {
		fmt.Printf("Info: "+format+"\n", a...)
	}
}

func (l *Fmt) Warn(a ...interface{}) {
	if l.level <= LevelWarn {
		fmt.Println(append([]interface{}{"Warn:"}, a...)...)
	}
}

func (l *Fmt) Warnf(format string, a ...interface{}) {
	if l.level <= LevelWarn {
		fmt.Printf("Warn: "+format+"\n", a...)
	}
}

func (l *Fmt) Error(a ...interface{}) {
	if l.level <= LevelError {
		fmt.Println(append([]interface{}{"Error:"}, a...)...)
	}
}

func (l *Fmt) Errorf(format string, a ...interface{}) {
	if l.level <= LevelError {
		fmt.Printf("Error: "+format+"\n", a...)
	}
}

func (l *Fmt) Fatal(a ...interface{}) {
	fmt.Println(append([]interface{}{"Fatal:"}, a...)...)
	os.Exit(1)
}

func (l *Fmt) Fatalf(format string, args ...interface{}) {
	fmt.Printf("Fatal: "+format+"\n", args...)
	os.Exit(1)
}

func (l *Fmt) Must(err error) {
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
}
