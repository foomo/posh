package prompt

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/foomo/posh/pkg/command"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/check"
	"github.com/foomo/posh/pkg/prompt/flair"
	"github.com/foomo/posh/pkg/prompt/history"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/shell"
)

type (
	Prompt struct {
		l        log.Logger
		ctx      context.Context
		title    string
		flair    flair.Flair
		prefix   string
		check    check.Check
		checkers []check.Checker
		filter   Filter
		readline *readline.Readline
		history  history.History
		commands command.Commands
		// inputRegex - split cmd into args
		promptOptions []prompt.Option
	}
	Option func(*Prompt) error
)

// ------------------------------------------------------------------------------------------------
// ~ Options
// ------------------------------------------------------------------------------------------------

func WithTitle(v string) Option {
	return func(o *Prompt) error {
		o.title = v
		return nil
	}
}

func WithFlair(v flair.Flair) Option {
	return func(o *Prompt) error {
		o.flair = v
		return nil
	}
}

func WithPrefix(v string) Option {
	return func(o *Prompt) error {
		o.prefix = v + " "
		return nil
	}
}

func WithCheck(v check.Check) Option {
	return func(o *Prompt) error {
		o.check = v
		return nil
	}
}

func WithCheckers(v ...check.Checker) Option {
	return func(o *Prompt) error {
		o.checkers = append(o.checkers, v...)
		return nil
	}
}

func WithContext(v context.Context) Option {
	return func(o *Prompt) error {
		o.ctx = v
		return nil
	}
}

func WithHistory(v history.History) Option {
	return func(o *Prompt) error {
		o.history = v
		return nil
	}
}

func WithFileHistory(v ...history.FileOption) Option {
	return func(o *Prompt) error {
		if value, err := history.NewFile(o.l, v...); err != nil {
			return err
		} else {
			o.history = value
			return nil
		}
	}
}

func WithCommands(v command.Commands) Option {
	return func(p *Prompt) error {
		p.commands = v
		return nil
	}
}

func WithFilter(v Filter) Option {
	return func(o *Prompt) error {
		o.filter = v
		return nil
	}
}

func WithPromptOptions(v ...prompt.Option) Option {
	return func(o *Prompt) error {
		o.promptOptions = v
		return nil
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(l log.Logger, opts ...Option) (*Prompt, error) {
	inst := &Prompt{
		l:        l,
		ctx:      context.Background(),
		title:    "posh",
		prefix:   "> ",
		flair:    flair.DefaultFlair,
		filter:   FilterFuzzy,
		check:    check.DefaultCheck,
		history:  history.NewNoop(l),
		commands: command.Commands{},
	}
	for _, opt := range opts {
		if opt != nil {
			if err := opt(inst); err != nil {
				return nil, err
			}
		}
	}
	if value, err := readline.New(l); err != nil {
		return nil, err
	} else {
		inst.readline = value
	}
	return inst, nil
}

// ------------------------------------------------------------------------------------------------
// ~ Getter
// ------------------------------------------------------------------------------------------------

func (s *Prompt) Commands() command.Commands {
	return s.commands
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (s *Prompt) Run() error {
	histories, err := s.history.Load(s.ctx)
	if err != nil {
		return err
	}

	p := prompt.New(
		s.execute,
		s.complete,
		append(
			[]prompt.Option{
				prompt.OptionTitle(s.title),
				prompt.OptionPrefix(s.prefix),
				prompt.OptionPrefixTextColor(prompt.Cyan),
				prompt.OptionInputTextColor(prompt.DefaultColor),
				prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
				prompt.OptionHistoryIgnoreDuplicates(),
				prompt.OptionHistory(histories),
				// macos alt+left fix
				prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
					ASCIICode: []byte{0x1b, 0x62},
					Fn:        prompt.GoLeftWord,
				}),
				// macos alt+right fix
				prompt.OptionAddASCIICodeBind(prompt.ASCIICodeBind{
					ASCIICode: []byte{0x1b, 0x66},
					Fn:        prompt.GoRightWord,
				}),
			},
			s.promptOptions...,
		)...,
	)

	if err := s.flair(s.title); err != nil {
		return err
	}

	if err := s.check(s.ctx, s.l, s.checkers); err != nil {
		return err
	}

	p.Run()

	return nil
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (s *Prompt) execute(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	s.history.Persist(s.ctx, input)

	if err := s.readline.Parse(input); err != nil {
		s.l.Error("failed to parse line:", err.Error())
		return
	}

	ctx := s.context()
	if cmd := s.Commands().Get(s.readline.Cmd()); cmd != nil {
		s.l.Debugf(`executing command:

> %s

%s
`, input, s.readline.String())
		if value, ok := cmd.(command.Validator); ok {
			if err := value.Validate(ctx, s.readline); err != nil {
				s.l.Error(err.Error())
			}
		}
		if err := cmd.Execute(ctx, s.readline); err != nil {
			s.l.Error(err.Error())
		}
	} else if err := shell.New(ctx, s.l, input).Run(); err != nil {
		s.l.Error(err.Error())
	}
}

func (s *Prompt) complete(d prompt.Document) []prompt.Suggest {
	input := d.TextBeforeCursor()
	if input == "" {
		return nil
	}

	if err := s.readline.Parse(input); err != nil {
		s.l.Debug("failed to parse line:", err.Error())
		return nil
	}
	word := d.GetWordBeforeCursor()

	// return root completion
	if s.readline.IsModeDefault() && s.readline.Args().LenIs(0) {
		var suggests []prompt.Suggest
		for _, inst := range s.Commands().List() {
			suggests = append(suggests, prompt.Suggest{Text: inst.Name(), Description: inst.Description()})
		}
		return s.filter(suggests, word, true)
	}

	// retrieve command instance
	cmd := s.commands.Get(s.readline.Cmd())
	if cmd == nil {
		return nil
	}

	ctx := s.context()

	switch s.readline.Mode() {
	case readline.ModeArgs:
		if value, ok := cmd.(command.ArgumentCompleter); ok {
			return s.filter(value.CompleteArguments(ctx, s.readline, d), word, true)
		} else if value, ok := cmd.(command.Completer); ok {
			return s.filter(value.Complete(ctx, s.readline, d), word, true)
		}
	case readline.ModeFlags:
		if value, ok := cmd.(command.FlagCompleter); ok {
			return s.filter(value.CompleteFlags(ctx, s.readline, d), word, true)
		} else if value, ok := cmd.(command.Completer); ok {
			return s.filter(value.Complete(ctx, s.readline, d), word, true)
		}
	case readline.ModePassThroughArgs:
		if value, ok := cmd.(command.PassThroughArgsCompleter); ok {
			return s.filter(value.CompletePassTroughArgs(ctx, s.readline, d), word, true)
		} else if value, ok := cmd.(command.Completer); ok {
			return s.filter(value.Complete(ctx, s.readline, d), word, true)
		}
	case readline.ModePassThroughFlags:
		if value, ok := cmd.(command.PassThroughFlagsCompleter); ok {
			return s.filter(value.CompletePassTroughFlags(ctx, s.readline, d), word, true)
		} else if value, ok := cmd.(command.Completer); ok {
			return s.filter(value.Complete(ctx, s.readline, d), word, true)
		}
	case readline.ModeAdditionalArgs:
		if value, ok := cmd.(command.AdditionalArgsCompleter); ok {
			return s.filter(value.CompleteAdditionalArgs(ctx, s.readline, d), word, true)
		} else if value, ok := cmd.(command.Completer); ok {
			return s.filter(value.Complete(ctx, s.readline, d), word, true)
		}
	}
	return nil
}

// context returns and watches over a new context
func (s *Prompt) context() context.Context {
	ctx, cancel := context.WithCancel(s.ctx)
	go func(ctx context.Context, cancel context.CancelFunc) {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		select {
		case <-sigChan:
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}(ctx, cancel)
	return ctx
}
