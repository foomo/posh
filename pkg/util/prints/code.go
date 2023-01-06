package prints

import (
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/foomo/posh/pkg/log"
)

func Code(l log.Logger, title, code, lexer string) {
	border := strings.Repeat("-", 80)
	l.Infof("\n%s\n%s\n%s", border, title, border)
	if err := quick.Highlight(os.Stdout, code, lexer, "terminal", "monokai"); err != nil {
		l.Debug(err.Error())
		l.Print(code)
	}
}
