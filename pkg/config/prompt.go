package config

type (
	Prompt struct {
		Title   string            `json:"title" yaml:"title"`
		Prefix  string            `json:"prefix" yaml:"prefix"`
		History PromptHistory     `json:"history" yaml:"history"`
		Aliases map[string]string `json:"aliases" yaml:"aliases"`
	}
	PromptHistory struct {
		Limit        int    `json:"limit" yaml:"limit"`
		Filename     string `json:"filename" yaml:"filename"`
		LockFilename string `json:"lockFilename" yaml:"lockFilename"`
	}
)
