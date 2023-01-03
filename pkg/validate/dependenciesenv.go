package validate

import (
	"errors"
	"os"

	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
)

type DependenciesEnvRule func(l log.Logger, v config.DependenciesEnv) rule.Rule

func DependenciesEnvs(l log.Logger, v []config.DependenciesEnv) []fend.Fend {
	ret := make([]fend.Fend, len(v))
	for i, vv := range v {
		ret[i] = DependenciesEnv(l, vv, DependenciesEnvExists)
	}
	return ret
}

func DependenciesEnv(l log.Logger, v config.DependenciesEnv, rules ...DependenciesEnvRule) fend.Fend {
	return func() []rule.Rule {
		ret := make([]rule.Rule, len(rules))
		for i, r := range rules {
			ret[i] = r(l, v)
		}
		return ret
	}
}

func DependenciesEnvExists(l log.Logger, v config.DependenciesEnv) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate env exists:", v.String())
		if value := os.Getenv(v.Name); value == "" {
			return nil, errors.New(v.Help)
		}
		return nil, nil
	}
}
