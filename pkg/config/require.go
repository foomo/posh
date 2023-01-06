package config

import (
	"fmt"
)

type (
	Require struct {
		Envs     []RequireEnv     `json:"envs" yaml:"envs"`
		Scripts  []RequireScript  `json:"scripts" yaml:"scripts"`
		Packages []RequirePackage `json:"packages" yaml:"packages"`
	}
	RequireEnv struct {
		Name string `json:"name" yaml:"name"`
		Help string `json:"help" yaml:"help"`
	}
	RequireScript struct {
		Name    string `json:"name" yaml:"name"`
		Command string `json:"command" yaml:"command"`
		Help    string `json:"help" yaml:"help"`
	}
	RequirePackage struct {
		Name    string `json:"name" yaml:"name"`
		Version string `json:"version" yaml:"version"`
		Command string `json:"command" yaml:"command"`
		Help    string `json:"help" yaml:"help"`
	}
)

func (c RequireEnv) String() string {
	return fmt.Sprintf("Name: %s", c.Name)
}

func (c RequireScript) String() string {
	return fmt.Sprintf("Name: %s, Command: %s", c.Name, c.Command)
}

func (c RequirePackage) String() string {
	return fmt.Sprintf("Name: %s, Version: %s Command: %s", c.Name, c.Version, c.Command)
}
