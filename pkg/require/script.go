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

func Scripts(l log.Logger, v []config.RequireScript) fend.Fends {
	ret := make(fend.Fends, len(v))
	for i, vv := range v {
		ret[i] = fend.Var(vv, ScriptStatus(l))
	}
	return ret
}

func ScriptStatus(l log.Logger) rule.Rule[config.RequireScript] {
	return func(ctx context.Context, v config.RequireScript) error {
		l.Debug("validate script status:", v.String())
		if output, err := exec.CommandContext(ctx, "sh", "-c", v.Command).CombinedOutput(); err != nil {
			l.Debug(string(output))
			return errors.Wrap(err, v.Help)
		}
		return nil
	}
}
