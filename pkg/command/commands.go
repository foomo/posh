package command

import (
	"plugin"
	"sort"

	"github.com/pkg/errors"
)

type Commands map[string]Command

func (c Commands) Get(name string) Command {
	if value, ok := c[name]; ok {
		return value
	} else {
		return nil
	}
}

func (c Commands) List() []Command {
	ret := make([]Command, 0, len(c))
	for _, command := range c {
		ret = append(ret, command)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name() < ret[j].Name()
	})
	return ret
}

func (c Commands) Has(name string) bool {
	return c.Get(name) != nil
}

func (c Commands) Add(commands ...Command) {
	for _, command := range commands {
		c[command.Name()] = command
	}
}

func (c Commands) Load(paths ...string) error {
	for _, path := range paths {
		if plg, err := plugin.Open(path); err != nil {
			return errors.Wrapf(err, "failed to load plugin (%s)", path)
		} else if sym, err := plg.Lookup("Commands"); err != nil {
			return errors.Wrapf(err, "failed to lookup Commands from plugin (%s)", path)
		} else if cmds, ok := sym.([]Command); !ok {
			return errors.Wrapf(err, "failed to lookup Commands type from plugin (%s)", path)
		} else {
			c.Add(cmds...)
		}
	}
	return nil
}
