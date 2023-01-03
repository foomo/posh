package config

import (
	"fmt"
)

type (
	Dependencies struct {
		Envs     []DependenciesEnv     `json:"envs" yaml:"envs"`
		Scripts  []DependenciesScript  `json:"script" yaml:"script"`
		Packages []DependenciesPackage `json:"packages" yaml:"packages"`
	}
	DependenciesEnv struct {
		Name string `json:"name" yaml:"name"`
		Help string `json:"help" yaml:"help"`
	}
	DependenciesScript struct {
		Name    string `json:"name" yaml:"name"`
		Command string `json:"command" yaml:"command"`
		Help    string `json:"help" yaml:"help"`
	}
	DependenciesPackage struct {
		Name    string `json:"name" yaml:"name"`
		Version string `json:"version" yaml:"version"`
		Command string `json:"command" yaml:"command"`
		Help    string `json:"help" yaml:"help"`
	}
)

func (c DependenciesEnv) String() string {
	return fmt.Sprintf("Name: %s", c.Name)
}

func (c DependenciesScript) String() string {
	return fmt.Sprintf("Name: %s, Command: %s", c.Name, c.Command)
}

func (c DependenciesPackage) String() string {
	return fmt.Sprintf("Name: %s, Version: %s Command: %s", c.Name, c.Version, c.Command)
}
