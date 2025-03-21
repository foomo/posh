package browser

import (
	"context"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/pkg/errors"
)

func OpenRawURL(ctx context.Context, u string) error {
	if u, err := url.Parse(u); err != nil {
		return err
	} else {
		return OpenURL(ctx, u)
	}
}

func OpenURL(ctx context.Context, u *url.URL) error {
	if u == nil {
		return errors.New("empty url")
	}
	switch runtime.GOOS {
	case "linux":
		return exec.CommandContext(ctx, "xdg-open", u.String()).Start()
	case "windows":
		return exec.CommandContext(ctx, "rundll32", "url.dll,FileProtocolHandler", u.String()).Start()
	case "darwin":
		return exec.CommandContext(ctx, "open", u.String()).Start()
	default:
		return errors.New("unsupported platform")
	}
}
