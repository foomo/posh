package validate

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

type DependenciesPackageRule func(ctx context.Context, l log.Logger, v config.DependenciesPackage) rule.Rule

func DependenciesPackages(ctx context.Context, l log.Logger, v []config.DependenciesPackage) []fend.Fend {
	ret := make([]fend.Fend, len(v))
	for i, vv := range v {
		ret[i] = DependenciesPackage(ctx, l, vv, DependenciesPackageExists, DependenciesPackageVersion)
	}
	return ret
}

func DependenciesPackage(ctx context.Context, l log.Logger, v config.DependenciesPackage, rules ...DependenciesPackageRule) fend.Fend {
	return func() []rule.Rule {
		ret := make([]rule.Rule, len(rules))
		for i, r := range rules {
			ret[i] = r(ctx, l, v)
		}
		return ret
	}
}

func DependenciesPackageExists(ctx context.Context, l log.Logger, v config.DependenciesPackage) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate package exists:", v.String())
		if output, err := exec.LookPath(v.Name); err != nil {
			l.Debug(err.Error(), output)
			return nil, fmt.Errorf(v.Help, v.Version)
		} else if output == "" {
			l.Debugf("missing executable %s", v.Name)
			return nil, fmt.Errorf(v.Help, v.Version)
		} else {
			return nil, nil
		}
	}
}

func DependenciesPackageVersion(ctx context.Context, l log.Logger, v config.DependenciesPackage) rule.Rule {
	return func() (*rule.Error, error) {
		l.Debug("validate package version:", v.String())
		if output, err := exec.CommandContext(ctx, "sh", "-c", v.Command).CombinedOutput(); err != nil {
			return nil, err
		} else if actual := strings.TrimPrefix(strings.TrimSpace(string(output)), "v"); actual == "" {
			return nil, fmt.Errorf("failed to retrieve version: %s", string(output))
		} else if c, err := semver.NewConstraint(v.Version); err != nil {
			return nil, errors.Wrapf(err, "failed to create version constraint: %s", v.Version)
		} else if version, err := semver.NewVersion(actual); err != nil {
			return nil, errors.Wrapf(err, "failed to create version")
		} else if !c.Check(version) {
			l.Debug("wrong package version:", actual)
			return nil, fmt.Errorf(v.Help, v.Version)
		} else {
			return nil, nil
		}
	}
}
