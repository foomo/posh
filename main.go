package main

import (
	"github.com/foomo/posh/cmd"
)

func init() {
	cmd.Init(nil)
}

func main() {
	cmd.Execute()
}
