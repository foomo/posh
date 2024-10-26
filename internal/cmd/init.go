package cmd

import (
	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/viper"
)

func NewLogger() log.Logger {
	return log.NewPTerm(
		log.PTermWithDisableColor(viper.GetBool("no-color")),
		log.PTermWithLevel(log.GetLevel(viper.GetString("level"))),
	)
}
