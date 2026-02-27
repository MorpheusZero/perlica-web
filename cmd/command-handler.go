package cmd

import (
	"errors"
	"os"
	"strings"
)

type ICommand interface {
	Run() error
}

type CommandHandler struct {
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{}
}

func (h *CommandHandler) HandleCommand() error {
	args := os.Args[1:]

	normalizedCommand := ""
	if len(args) > 0 {
		normalizedCommand = strings.ToLower(args[0])
	} else {
		normalizedCommand = "help"
	}

	switch normalizedCommand {
	case "help":
		cmd := NewHelpCommand()
		return cmd.Run()
	case "version":
		cmd := NewVersionCommand()
		return cmd.Run()
	case "migrate":
		cmd := NewMigrateCommand()
		return cmd.Run()
	case "server":
		cmd := NewServerCommand()
		return cmd.Run()
	default:
		return errors.New("unknown command: " + normalizedCommand + "\nUse 'help' command to see available commands")
	}
}
