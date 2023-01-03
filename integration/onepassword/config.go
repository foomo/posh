package onepassword

type Config struct {
	Account       string `json:"account" yaml:"account"`
	TokenFilename string `json:"tokenFilename" yaml:"tokenFilename"`
}
