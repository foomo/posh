package plugin

import (
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/plugin"
)

var m *plugin.Manager

func manager(l log.Logger) (*plugin.Manager, error) {
	if m == nil {
		if value, err := plugin.NewManager(l); err != nil {
			return nil, err
		} else {
			m = value
		}
	}
	return m, nil
}
