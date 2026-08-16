package cmd

import (
	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/log"
	"github.com/spf13/viper"
)

func NewLogger() log.Logger {
	level := log.GetLevel(viper.GetString("level"))

	if agent.IsAgentMode() {
		return log.NewAgentJSON(log.AgentJSONWithLevel(level))
	}

	return log.NewPTerm(
		log.PTermWithDisableColor(viper.GetBool("no-color")),
		log.PTermWithLevel(level),
	)
}
