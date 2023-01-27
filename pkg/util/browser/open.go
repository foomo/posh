package browser

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

func OpenRawURL(u string) error {
	if u, err := url.Parse(u); err != nil {
		return err
	} else {
		return OpenURL(u)
	}
}

func OpenURL(u *url.URL) error {
	if u == nil {
		return fmt.Errorf("empty url")
	}
	switch runtime.GOOS {
	case "linux":
		return exec.CommandContext(context.TODO(), "xdg-open", u.String()).Start()
	case "windows":
		return exec.CommandContext(context.TODO(), "rundll32", "url.dll,FileProtocolHandler", u.String()).Start()
	case "darwin":
		return exec.CommandContext(context.TODO(), "open", u.String()).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
