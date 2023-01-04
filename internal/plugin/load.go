package plugin

import (
	"context"

	"github.com/foomo/posh/pkg/config"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/plugin"
	"github.com/spf13/viper"
)

func Load(ctx context.Context, l log.Logger) (plugin.Plugin, error) {
	m, err := manager(l)
	if err != nil {
		return nil, err
	}

	var cfg config.Plugin
	if err := viper.UnmarshalKey("plugin", &cfg); err != nil {
		return nil, err
	}
	return m.BuildAndLoadPlugin(ctx, cfg.Source, cfg.Provider)
}
