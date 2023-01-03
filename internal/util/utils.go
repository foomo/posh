package util

import (
	"context"

	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/spf13/viper"
)

func LoadPlugin(ctx context.Context, m *plugin.Manager) (plugin.Plugin, error) {
	var cfg config.Plugin
	if err := viper.UnmarshalKey("plugin", &cfg); err != nil {
		return nil, err
	}
	return m.BuildAndLoadPlugin(ctx, cfg.Source, cfg.Provider)
}
