package log

import (
	"github.com/foomo/posh/pkg/log"
)

func Init(level string, noColor bool) (log.Logger, error) {
	if value, err := log.NewPTerm(
		log.PTermWithDisableColor(noColor),
		log.PTermWithLevel(log.GetLevel(level)),
	); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
