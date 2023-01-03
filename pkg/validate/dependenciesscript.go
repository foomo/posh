package validate

import (
	"context"
	"os/exec"

	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type DependenciesScriptRule func(ctx context.Context, l log.Logger, v config.DependenciesScript) rule.Rule

func DependenciesScripts(ctx context.Context, l log.Logger, v []config.DependenciesScript) []fend.Fend {
	ret := make([]fend.Fend, len(v))
	for i, vv := range v {
		ret[i] = DependenciesScript(ctx, l, vv, DependenciesScriptStatus)
	}
	return ret
}

func DependenciesScript(ctx context.Context, l log.Logger, v config.DependenciesScript, rules ...DependenciesScriptRule) fend.Fend {
	return func() []rule.Rule {
		ret := make([]rule.Rule, len(rules))
		for i, r := range rules {
			ret[i] = r(ctx, l, v)
		}
		return ret
	}
}

func DependenciesScriptStatus(ctx context.Context, l log.Logger, v config.DependenciesScript) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate script status:", v.String())
		if output, err := exec.CommandContext(ctx, "sh", "-c", v.Command).CombinedOutput(); err != nil {
			l.Debug(string(output))
			return nil, errors.Wrap(err, v.Help)
		}
		return nil, nil
	}
}
