package require

import (
	"context"
	"os/exec"
	"regexp"
	"strings"

	"github.com/foomo/fender/fend"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type GitRule func(l log.Logger) fend.Fend

func GitUser(l log.Logger, rules ...GitRule) fend.Fends {
	fends := make(fend.Fends, len(rules))
	for i, r := range rules {
		fends[i] = r(l)
	}

	return fends
}

func GitUserName(l log.Logger) fend.Fend {
	return fend.Var("", func(ctx context.Context, v string) error {
		l.Debug("validate git user.name")

		if output, err := exec.CommandContext(ctx, "git", "config", "user.name").CombinedOutput(); err != nil {
			return errors.Wrap(err, string(output))
		} else if parts := strings.Split(trim(string(output)), " "); len(parts) < 2 {
			return errors.New(`
Please configure a human readable git name e.g. "Max Mustermann" instead of "` + parts[0] + `".

$ git config user.name "Max Mustermann"
`)
		}

		return nil
	})
}

func GitUserEmail(pattern string) GitRule {
	return func(l log.Logger) fend.Fend {
		reg := regexp.MustCompile(pattern)

		return fend.Var("", func(ctx context.Context, v string) error {
			l.Debug("validate git user.email")

			if output, err := exec.CommandContext(ctx, "git", "config", "user.email").CombinedOutput(); err != nil {
				return errors.Wrap(err, string(output))
			} else if output := trim(string(output)); !reg.MatchString(output) {
				return errors.New(`
Please configure your github email to match the pattern "` + pattern + ` instead of "` + output + `".

$ git config user.email "max.muster@dev.null"
`)
			}

			return nil
		})
	}
}

func trim(s string) string {
	return strings.TrimSuffix(strings.TrimSpace(s), "\n")
}
