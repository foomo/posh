package require

import (
	"context"
	"os/exec"
	"regexp"
	"strings"

	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type GitRule func(ctx context.Context, l log.Logger) rule.Rule

func GitUser(ctx context.Context, l log.Logger, rules ...GitRule) fend.Fend {
	return func() []rule.Rule {
		ret := make([]rule.Rule, len(rules))
		for i, r := range rules {
			ret[i] = r(ctx, l)
		}
		return ret
	}
}

func GitUserName(ctx context.Context, l log.Logger) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate git user.name")
		if output, err := exec.CommandContext(ctx, "git", "config", "user.name").CombinedOutput(); err != nil {
			return nil, errors.Wrap(err, string(output))
		} else if parts := strings.Split(trim(string(output)), " "); len(parts) < 2 {
			return nil, errors.New(`
Please configure a human readable git name e.g. "Max Mustermann" instead of "` + parts[0] + `".

$ git config user.name "Max Mustermann"
`)
		}
		return nil, nil
	}
}

func GitUserEmail(pattern string) GitRule {
	reg := regexp.MustCompile(pattern)
	return func(ctx context.Context, l log.Logger) rule.Rule {
		return func() (*rule.Error, error) {
			l.Debug("validate git user.email")
			if output, err := exec.CommandContext(ctx, "git", "config", "user.email").CombinedOutput(); err != nil {
				return nil, errors.Wrap(err, string(output))
			} else if output := trim(string(output)); !reg.MatchString(output) {
				return nil, errors.New(`
Please configure your github email to match the pattern "` + pattern + ` instead of "` + output + `".

$ git config user.email "max.muster@dev.null"
`)
			}
			return nil, nil
		}
	}
}

func trim(s string) string {
	return strings.TrimSuffix(strings.TrimSpace(s), "\n")
}
