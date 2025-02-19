package require

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/foomo/fender/fend"
	"github.com/foomo/fender/rule"
	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/pkg/errors"
)

func Packages(l log.Logger, v []config.RequirePackage) fend.Fends {
	ret := make(fend.Fends, len(v))
	for i, vv := range v {
		ret[i] = fend.Var(vv, PackageExists(l), PackageVersion(l))
	}
	return ret
}

func PackageExists(l log.Logger) rule.Rule[config.RequirePackage] {
	return func(ctx context.Context, v config.RequirePackage) error {
		l.Debug("validate package exists:", v.String())
		if output, err := exec.LookPath(v.Name); err != nil {
			l.Debug(err.Error(), output)
			return errors.Errorf(v.Help, v.Version)
		} else if output == "" {
			l.Debugf("missing executable %s", v.Name)
			return errors.Errorf(v.Help, v.Version)
		}
		return nil
	}
}

func PackageVersion(l log.Logger) rule.Rule[config.RequirePackage] {
	return func(ctx context.Context, v config.RequirePackage) error {
		l.Debug("validate package version:", v.String())
		output, err := exec.CommandContext(ctx, "sh", "-c", v.Command).CombinedOutput()
		if err != nil {
			l.Error("failed to validate package version:", v.String(), ", output:", string(output), ", err:", err.Error())
			return errors.Wrap(err, string(output))
		}

		actual := strings.TrimPrefix(strings.TrimSpace(string(output)), "v")
		if actual == "" {
			l.Error("failed to retrieve version: ", string(output))
			return errors.Errorf(v.Help, v.Version)
		}

		if c, err := semver.NewConstraint(v.Version); err != nil {
			return errors.Wrapf(err, "failed to create version constraint: %s", v.Version)
		} else if version, err := semver.NewVersion(actual); err != nil {
			return errors.Wrapf(err, "failed to create version")
		} else if !c.Check(version) {
			l.Debug("wrong package version:", actual)
			return errors.Errorf(v.Help, v.Version)
		}
		return nil
	}
}
