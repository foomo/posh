package git

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/shell"
	"github.com/pkg/errors"
)

func ConfigUserName(ctx context.Context, l log.Logger) (string, error) {
	value, err := shell.New(ctx, l, "git config user.name").Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve git user name")
	}

	return strings.TrimSpace(string(value)), nil
}

func ConfigUserEmail(ctx context.Context, l log.Logger) (string, error) {
	value, err := shell.New(ctx, l, "git config user.email").Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve git user name")
	}

	return strings.TrimSpace(string(value)), nil
}
