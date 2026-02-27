package main

import (
	"fmt"

	"github.com/moprheuszero/perlica-web/cmd"
)

func main() {
	fmt.Println(GetTerminalArt())
	handler := cmd.NewCommandHandler()
	err := handler.HandleCommand()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
