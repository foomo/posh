package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/foomo/posh/pkg/log"
)

// Shell struct
type Shell struct {
	l      log.Logger
	cmd    *exec.Cmd
	quiet  bool
	args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func New(ctx context.Context, l log.Logger, inputs ...interface{}) *Shell {
	var args []string
	for _, input := range inputs {
		if value, ok := input.(string); ok {
			args = append(args, value)
		} else if value, ok := input.([]string); ok {
			args = append(args, value...)
		} else {
			args = append(args, fmt.Sprintf("%v", args))
		}
	}
	cmd := exec.CommandContext(ctx, "sh", "-c")
	cmd.Env = os.Environ()
	return &Shell{
		l:      l.Named("shell"),
		cmd:    cmd,
		args:   args,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (s *Shell) Args(args ...string) *Shell {
	s.args = append(s.args, args...)
	return s
}

func (s *Shell) Env(env ...string) *Shell {
	s.cmd.Env = append(s.cmd.Env, env...)
	return s
}

func (s *Shell) Dir(dir string) *Shell {
	s.cmd.Dir = dir
	return s
}

func (s *Shell) Quiet() *Shell {
	s.quiet = true
	return s
}

func (s *Shell) Run() error {
	args := s.args
	s.cmd.Args = append(s.cmd.Args, strings.Join(args, " "))
	s.cmd.Stdin = s.Stdin
	if !s.quiet {
		s.cmd.Stdout = s.Stdout
		s.cmd.Stderr = s.Stderr
	}
	s.trace()
	return s.cmd.Run()
}

func (s *Shell) Output() ([]byte, error) {
	args := s.args
	s.cmd.Args = append(s.cmd.Args, strings.Join(args, " "))
	if !s.quiet {
		s.cmd.Stdin = s.Stdin
		s.cmd.Stderr = s.Stderr
	}
	s.trace()
	return s.cmd.Output()
}

func (s *Shell) CombinedOutput() ([]byte, error) {
	args := s.args
	s.cmd.Args = append(s.cmd.Args, strings.Join(args, " "))
	if !s.quiet {
		s.cmd.Stdin = s.Stdin
	}
	s.trace()
	return s.cmd.CombinedOutput()
}

func (s *Shell) Wait() error {
	args := s.args
	s.cmd.Args = append(s.cmd.Args, strings.Join(args, " "))
	s.cmd.Stdin = s.Stdin
	if !s.quiet {
		s.cmd.Stdout = s.Stdout
		s.cmd.Stderr = s.Stderr
	}
	s.trace()
	// start the process and wait till it's finished
	if err := s.cmd.Start(); err != nil {
		return err
	} else if err := s.cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (s *Shell) trace() {
	s.l.Tracef(`"Executing:
$ %s

Directory: %s

%s
`,
		s.cmd.String(),
		s.cmd.Dir,
		strings.Join(s.cmd.Environ(), "\n"),
	)
}
