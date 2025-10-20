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
		l              log.Logger
		mu             sync.RWMutex
		cmd            string
		mode           Mode
		args           Args
		flags          Args
		flagSets       *FlagSets
		additionalArgs Args
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

func (a *Readline) FlagSets() *FlagSets {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.flagSets
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

	for i, part := range parts {
		if a.mode == ModeArgs && Arg(part).IsFlag() {
			a.mode = ModeFlags
		}

		if Arg(part).IsAdditional() && i < len(parts)-1 {
			a.mode = ModeAdditionalArgs
		}

		switch a.mode {
		case ModeArgs:
			a.args = append(a.args, part)
		case ModeFlags:
			a.flags = append(a.flags, part)
		case ModeAdditionalArgs:
			a.additionalArgs = append(a.additionalArgs, part)
		}
	}

	return nil
}

func (a *Readline) SetFlagSets(fs *FlagSets) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.flagSets = fs
}

func (a *Readline) ParseFlagSets() error {
	if fs := a.FlagSets(); fs != nil {
		if err := fs.Parse(a.flags); err != nil {
			return err
		}
	}

	return nil
}

func (a *Readline) String() string {
	return fmt.Sprintf(`
Cmd:                  %s
Args:                 %s
Flags:                %s
AdditionalArgs:       %s
`, a.Cmd(), a.Args(), a.Flags(), a.AdditionalArgs())
}

func (a *Readline) IsModeDefault() bool {
	return a.Mode() == ModeArgs
}

func (a *Readline) IsModeAdditional() bool {
	return a.Mode() == ModeAdditionalArgs
}

func (a *Readline) AllFlags() []*pflag.Flag {
	var ret []*pflag.Flag
	if fs := a.FlagSets(); fs != nil {
		fs.All().VisitAll(func(f *pflag.Flag) {
			ret = append(ret, f)
		})
	}

	return ret
}

func (a *Readline) VisitedFlags() Flags {
	var ret Flags
	if fs := a.FlagSets(); fs != nil {
		ret = fs.Visited()
	}

	return ret
}

func (a *Readline) AdditionalFlags() Args {
	ret := append(Args{}, a.flags...)
	if fs := a.FlagSets(); fs != nil {
		fs.VisitAll(func(f *pflag.Flag) {
			if i := ret.IndexOf("--" + f.Name); i >= 0 {
				switch f.Value.Type() {
				case "bool":
					ret = ret.Splice(ret.IndexOf("--"+f.Name), 1)
				default:
					ret = ret.Splice(ret.IndexOf("--"+f.Name), 2)
				}
			}
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
	a.flagSets = nil
	a.additionalArgs = nil
}
