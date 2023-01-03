package config

import (
	"fmt"
)

type Plugin struct {
	Source   string `json:"source" yaml:"source"`
	Target   string `json:"target" yaml:"target"`
	Provider string `json:"provider" yaml:"provider"`
}

func (c Plugin) String() string {
	return fmt.Sprintf("Filename: %s, Provider: %s", c.Source, c.Provider)
}
