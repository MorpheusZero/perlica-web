package cmd

import (
	"fmt"

	"github.com/moprheuszero/perlica-web/constants"
)

type VersionCommand struct {
}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

func (c *VersionCommand) Run() error {
	versionText := fmt.Sprintf(`Perlica - v%s
`, constants.AppReleaseVersion+" (Build: "+constants.APIBuildVersion+")")
	fmt.Println(versionText)
	return nil
}
