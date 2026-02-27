package cmd

import "github.com/moprheuszero/perlica-web/server"

type ServerCommand struct {
}

func NewServerCommand() *ServerCommand {
	return &ServerCommand{}
}

func (c *ServerCommand) Run() error {
	server := server.NewAppServer()
	err := server.Start()
	if err != nil {
		return err
	}
	return nil
}
