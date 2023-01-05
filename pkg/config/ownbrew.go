package config

import (
	"github.com/foomo/posh/integration/ownbrew"
)

type Ownbrew struct {
	Dry       bool              `json:"dry" yaml:"dry"`
	BinDir    string            `json:"binDir" yaml:"binDir"`
	TapDir    string            `json:"tapDir" yaml:"tapDir"`
	TempDir   string            `json:"tempDir" yaml:"tempDir"`
	CellarDir string            `json:"cellarDir" yaml:"cellarDir"`
	Packages  []ownbrew.Package `json:"packages" yaml:"packages"`
}
