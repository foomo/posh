package git

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/shell"
	"github.com/pkg/errors"
)

func Ref(ctx context.Context, l log.Logger) (string, error) {
	value, err := shell.New(ctx, l, "git rev-parse --abbrev-ref HEAD").Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve git ref")
	}

	if strings.TrimSpace(string(value)) == "HEAD" {
		value, err = shell.New(ctx, l, "git describe --tags").Output()
		if err != nil {
			return "", errors.Wrap(err, "failed to retrieve git tag")
		}
	}

	return strings.TrimSpace(string(value)), nil
}
