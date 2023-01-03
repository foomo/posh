package log

import (
	"os"

	"github.com/pterm/pterm"
)

type (
	PTerm struct {
		level Level
	}
	PTermOption func(*PTerm) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func PTermWithDisableColor(v bool) PTermOption {
	return func(o *PTerm) error {
		if v {
			pterm.DisableColor()
		}
		return nil
	}
}

func PTermWithLevel(v Level) PTermOption {
	return func(o *PTerm) error {
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
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewPTerm(opts ...PTermOption) (*PTerm, error) {
	inst := &PTerm{
		level: LevelError,
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (l *PTerm) Level() Level {
	return l.level
}

func (l *PTerm) Print(a ...interface{}) {
	pterm.Println(a...)
}

func (l *PTerm) Printf(format string, a ...interface{}) {
	pterm.Printfln(format, a...)
}
func (l *PTerm) Debug(a ...interface{}) {
	if l.level <= LevelDebug {
		pterm.Debug.Println(a...)
	}
}

func (l *PTerm) Debugf(format string, a ...interface{}) {
	if l.level <= LevelDebug {
		pterm.Debug.Printfln(format, a...)
	}
}

func (l *PTerm) Info(a ...interface{}) {
	if l.level <= LevelInfo {
		pterm.Info.Println(a...)
	}
}

func (l *PTerm) Infof(format string, a ...interface{}) {
	if l.level <= LevelInfo {
		pterm.Info.Printfln(format, a...)
	}
}

func (l *PTerm) Warn(a ...interface{}) {
	if l.level <= LevelWarn {
		pterm.Warning.Println(a...)
	}
}

func (l *PTerm) Warnf(format string, a ...interface{}) {
	if l.level <= LevelWarn {
		pterm.Warning.Printfln(format, a...)
	}
}

func (l *PTerm) Error(a ...interface{}) {
	if l.level <= LevelError {
		pterm.Error.Println(a...)
	}
}

func (l *PTerm) Errorf(format string, a ...interface{}) {
	if l.level <= LevelError {
		pterm.Error.Printfln(format, a...)
	}
}

func (l *PTerm) Fatal(a ...interface{}) {
	pterm.Fatal.Println(a...)
}

func (l *PTerm) Fatalf(format string, args ...interface{}) {
	pterm.Fatal.Printfln(format, args...)
}

func (l *PTerm) Must(err error) {
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
}
