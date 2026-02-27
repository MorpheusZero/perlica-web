package cmd

import (
	"fmt"

	"github.com/moprheuszero/perlica-web/config"
	"github.com/moprheuszero/perlica-web/constants"
	"github.com/morpheuszero/go-heimdall/v2"
)

type MigrateCommand struct {
	envProvider *config.EnvProvider
}

func NewMigrateCommand() *MigrateCommand {
	envProvider := config.NewEnvProvider()
	return &MigrateCommand{
		envProvider: envProvider,
	}
}

func (c *MigrateCommand) Run() error {
	migrateText := fmt.Sprintf(`Running database migrations for Perlica - v%s
`, constants.AppReleaseVersion+" (Build: "+constants.APIBuildVersion+")")
	fmt.Println(migrateText)

	config := heimdall.HeimdallConfig{
		ConnectionString:            c.envProvider.GetDBConnectionString(),
		MigrationTableName:          "migration_history",
		MigrationFilesDirectoryPath: "./migrations",
		Verbose:                     true,
	}

	h, err := heimdall.NewHeimdall(config)
	if err != nil {
		return err
	}
	defer h.Close()

	err = h.RunMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Database migrations completed successfully.")
	return nil
}
