package require

import (
	"context"
	"errors"
	"os"

	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
)

func Envs(l log.Logger, v []config.RequireEnv) fend.Fends {
	ret := make([]fend.Fend, len(v))
	for i, vv := range v {
		ret[i] = fend.Var(vv, EnvExists(l))
	}

	return ret
}

func EnvExists(l log.Logger) rule.Rule[config.RequireEnv] {
	return func(ctx context.Context, v config.RequireEnv) error {
		l.Debug("validate env exists:", v.String())

		if value := os.Getenv(v.Name); value == "" {
			return errors.New(v.Help)
		}

		return nil
	}
}
