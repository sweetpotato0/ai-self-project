package main

import (
	"os"

	"gin-web-framework/cmd/cli/commands"
)

func main() {
	rootCmd := commands.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
