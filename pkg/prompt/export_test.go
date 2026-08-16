package prompt

import (
	"github.com/c-bata/go-prompt"
)

// ThemeOptions exposes the theme's color options for the black-box
// prompt_test package, so a test can assert that an unconfigured prompt adds
// none of them without having to start an interactive session.
func ThemeOptions(p *Prompt) []prompt.Option {
	return p.themeOptions()
}
