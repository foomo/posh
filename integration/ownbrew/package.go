package ownbrew

import (
	"fmt"
	"strings"
)

type Package struct {
	Tap     string   `json:"tap" yaml:"tap"`
	Name    string   `json:"name" yaml:"name"`
	Names   []string `json:"names" yaml:"names"`
	Args    []string `json:"args" yaml:"args"`
	Version string   `json:"version" yaml:"version"`
}

func (c Package) AllNames() []string {
	names := c.Names
	if len(names) == 0 {
		names = append(names, c.Name)
	}
	return names
}

func (c Package) String() string {
	return fmt.Sprintf("Names: %s, Version: %s Tap: %s", c.AllNames(), c.Version, c.Tap)
}

func (c Package) URL() (string, error) {
	// foomo/tap/aws/kubectl
	parts := strings.Split(c.Tap, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid tap format: %s", c.Tap)
	}
	return fmt.Sprintf(
		"https://raw.githubusercontent.com/%s/ownbrew-%s/main/%s/%s.sh",
		parts[0],
		parts[1],
		parts[2],
		parts[3],
	), nil
}
