package require

import (
	"context"
	"os/exec"

	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type ScriptRule func(ctx context.Context, l log.Logger, v config.RequireScript) rule.Rule

func Scripts(ctx context.Context, l log.Logger, v []config.RequireScript) []fend.Fend {
	ret := make([]fend.Fend, len(v))
	for i, vv := range v {
		ret[i] = Script(ctx, l, vv, ScriptStatus)
	}
	return ret
}

func Script(ctx context.Context, l log.Logger, v config.RequireScript, rules ...ScriptRule) fend.Fend {
	return func() []rule.Rule {
		ret := make([]rule.Rule, len(rules))
		for i, r := range rules {
			ret[i] = r(ctx, l, v)
		}
		return ret
	}
}

func ScriptStatus(ctx context.Context, l log.Logger, v config.RequireScript) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate script status:", v.String())
		if output, err := exec.CommandContext(ctx, "sh", "-c", v.Command).CombinedOutput(); err != nil {
			l.Debug(string(output))
			return nil, errors.Wrap(err, v.Help)
		}
		return nil, nil
	}
}
