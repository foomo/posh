package plugin

import (
	"github.com/foomo/posh/pkg/log"
)

type Provider func(l log.Logger) (Plugin, error)
