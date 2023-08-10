package require

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
			return fmt.Errorf(v.Help, v.Version)
		} else if output == "" {
			l.Debugf("missing executable %s", v.Name)
			return fmt.Errorf(v.Help, v.Version)
		}
		return nil
	}
}

func PackageVersion(l log.Logger) rule.Rule[config.RequirePackage] {
	return func(ctx context.Context, v config.RequirePackage) error {
		l.Debug("validate package version:", v.String())
		if output, err := exec.CommandContext(ctx, "sh", "-c", v.Command).CombinedOutput(); err != nil {
			return err
		} else if actual := strings.TrimPrefix(strings.TrimSpace(string(output)), "v"); actual == "" {
			l.Debugf("failed to retrieve version: %s", string(output))
			return fmt.Errorf(v.Help, v.Version)
		} else if c, err := semver.NewConstraint(v.Version); err != nil {
			return errors.Wrapf(err, "failed to create version constraint: %s", v.Version)
		} else if version, err := semver.NewVersion(actual); err != nil {
			return errors.Wrapf(err, "failed to create version")
		} else if !c.Check(version) {
			l.Debug("wrong package version:", actual)
			return fmt.Errorf(v.Help, v.Version)
		}
		return nil
	}
}
