package log

import (
	"github.com/foomo/posh/pkg/log"
)

func Init(level string, noColor bool) log.Logger {
	return log.NewPTerm(
		log.PTermWithDisableColor(noColor),
		log.PTermWithLevel(log.GetLevel(level)),
	)
}
