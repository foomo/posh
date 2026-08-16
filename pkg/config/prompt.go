package config

type (
	Prompt struct {
		Title string `json:"title" yaml:"title"`
		// Theme selects a Catppuccin palette for the prompt: mocha, macchiato, frappe or latte
		Theme         string            `json:"theme" yaml:"theme"`
		Prefix        string            `json:"prefix" yaml:"prefix"`
		PrefixGit     bool              `json:"prefixGit" yaml:"prefixGit"`
		History       PromptHistory     `json:"history" yaml:"history"`
		HistorySearch bool              `json:"historySearch" yaml:"historySearch"`
		Aliases       map[string]string `json:"aliases" yaml:"aliases"`
	}
	PromptHistory struct {
		Limit        int    `json:"limit" yaml:"limit"`
		Filename     string `json:"filename" yaml:"filename"`
		LockFilename string `json:"lockFilename" yaml:"lockFilename"`
	}
)
