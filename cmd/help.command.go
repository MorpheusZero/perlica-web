package cmd

import (
	"fmt"

	"github.com/moprheuszero/perlica-web/constants"
)

type HelpCommand struct {
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

func (c *HelpCommand) Run() error {
	helpText := fmt.Sprintf(`Perlica CLI - v%s

Available Commands:
- help: Display this help message
- version: Display the current version of the API
- server: Start the API Server
- migrate: Run database migrations

Usage:
To run a command, use the following format:
./perlica <command>

Example:
./perlica version
`, constants.AppReleaseVersion+" (Build: "+constants.APIBuildVersion+")")
	fmt.Println(helpText)
	return nil
}
