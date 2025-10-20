package git

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/shell"
	"github.com/pkg/errors"
)

func Ref(ctx context.Context, l log.Logger) (string, error) {
	value, err := shell.New(ctx, l,
		"git", "describe", "--tags", "--exact-match", "2>", "/dev/null",
		"||", "git", "symbolic-ref -q", "--short HEAD",
		"||", "git rev-parse", "--short", "HEAD",
	).Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve git ref")
	}

	return strings.TrimSpace(string(value)), nil
}
