package config

import (
	"fmt"
)

type Package struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Command string `json:"command" yaml:"command"`
}

func (c Package) String() string {
	return fmt.Sprintf("Name: %s, Version: %s Command: %s", c.Name, c.Version, c.Command)
}
