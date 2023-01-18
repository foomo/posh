package readline

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/pflag"
)

type (
	Readline struct {
		l                  log.Logger
		mu                 sync.RWMutex
		cmd                string
		mode               Mode
		args               Args
		flags              Args
		flagSet            *FlagSet
		passThroughFlags   Args
		passThroughFlagSet *FlagSet
		additionalArgs     Args
		// regex - split cmd into args (https://regex101.com/r/EgiOzv/1)
		regex *regexp.Regexp
	}
	Option func(*Readline) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithRegex(v *regexp.Regexp) Option {
	return func(o *Readline) error {
		o.regex = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(l log.Logger, opts ...Option) (*Readline, error) {
	inst := &Readline{
		l:     l.Named("readline"),
		regex: regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'|(\s$)`),
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
// ~ Getter
// ------------------------------------------------------------------------------------------------

func (a *Readline) Mode() Mode {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.mode
}

func (a *Readline) Cmd() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.cmd
}

func (a *Readline) Args() Args {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.args
}

func (a *Readline) Flags() Args {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.flags
}

func (a *Readline) FlagSet() *FlagSet {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.flagSet
}

func (a *Readline) PassThroughFlags() Args {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.passThroughFlags
}

func (a *Readline) PassThroughFlagSet() *FlagSet {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.passThroughFlagSet
}

func (a *Readline) AdditionalArgs() Args {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.additionalArgs
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (a *Readline) Parse(input string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.reset()
	parts := a.regex.FindAllString(input, -1)

	if len(parts) == 0 {
		return nil
	} else {
		a.cmd, parts = parts[0], parts[1:]
	}

	last := len(parts) - 1
	for i, part := range parts {
		if a.mode == ModeArgs && Arg(part).IsFlag() {
			a.mode = ModeFlags
		}
		if i != last && (a.mode == ModeArgs || a.mode == ModeFlags) && Arg(part).IsPass() {
			a.mode = ModePassThroughFlags
		}
		if Arg(part).IsAdditional() && i < len(parts)-1 {
			a.mode = ModeAdditionalArgs
		}

		switch a.mode {
		case ModeArgs:
			a.args = append(a.args, part)
		case ModeFlags:
			a.flags = append(a.flags, part)
		case ModePassThroughFlags:
			a.passThroughFlags = append(a.passThroughFlags, part)
		case ModeAdditionalArgs:
			a.additionalArgs = append(a.additionalArgs, part)
		}
	}

	return nil
}

func (a *Readline) SetFlags(fs *FlagSet) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.flagSet = fs
}

func (a *Readline) ParseFlags() error {
	if fs := a.FlagSet(); fs == nil {
		return nil
	} else if err := fs.Parse(a.flags); err != nil {
		return err
	}
	return nil
}

func (a *Readline) SetParsePassThroughFlags(fs *FlagSet) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.passThroughFlagSet = fs
}

func (a *Readline) ParsePassThroughFlags() error {
	if fs := a.PassThroughFlagSet(); fs == nil {
		return nil
	} else if err := fs.Parse(a.passThroughFlags); err != nil {
		return err
	}
	return nil
}

func (a *Readline) String() string {
	return fmt.Sprintf(`
Cmd:              %s
Mode              %s
Args:             %s
Flags:            %s
PassThroughFlags: %s
AdditionalArgs    %s
`, a.Cmd(), a.Mode(), a.Args(), a.Flags(), a.PassThroughFlags(), a.AdditionalArgs())
}

func (a *Readline) IsModeDefault() bool {
	return a.Mode() == ModeArgs
}

func (a *Readline) IsModePassThrough() bool {
	return a.Mode() == ModePassThroughFlags
}

func (a *Readline) IsModeAdditional() bool {
	return a.Mode() == ModeAdditionalArgs
}

func (a *Readline) AllFlags() []*pflag.Flag {
	var ret []*pflag.Flag
	if fs := a.FlagSet(); fs != nil {
		fs.VisitAll(func(f *pflag.Flag) {
			ret = append(ret, f)
		})
	}
	return ret
}

func (a *Readline) VisitedFlags() []*pflag.Flag {
	var ret []*pflag.Flag
	if fs := a.FlagSet(); fs != nil {
		fs.Visit(func(f *pflag.Flag) {
			ret = append(ret, f)
		})
	}
	return ret
}

func (a *Readline) AllPassThroughFlags() []*pflag.Flag {
	var ret []*pflag.Flag
	if fs := a.PassThroughFlagSet(); fs != nil {
		fs.VisitAll(func(f *pflag.Flag) {
			ret = append(ret, f)
		})
	}
	return ret
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (a *Readline) reset() {
	a.mode = ModeArgs
	a.cmd = ""
	a.args = nil
	a.flags = nil
	a.flagSet = nil
	a.passThroughFlags = nil
	a.passThroughFlagSet = nil
	a.additionalArgs = nil
}
