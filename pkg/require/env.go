package require

import (
	"errors"
	"os"

	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
)

type EnvRule func(l log.Logger, v config.RequireEnv) rule.Rule

func Envs(l log.Logger, v []config.RequireEnv) []fend.Fend {
	ret := make([]fend.Fend, len(v))
	for i, vv := range v {
		ret[i] = Env(l, vv, EnvExists)
	}
	return ret
}

func Env(l log.Logger, v config.RequireEnv, rules ...EnvRule) fend.Fend {
	return func() []rule.Rule {
		ret := make([]rule.Rule, len(rules))
		for i, r := range rules {
			ret[i] = r(l, v)
		}
		return ret
	}
}

func EnvExists(l log.Logger, v config.RequireEnv) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate env exists:", v.String())
		if value := os.Getenv(v.Name); value == "" {
			return nil, errors.New(v.Help)
		}
		return nil, nil
	}
}
