package config

type Env []struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}
