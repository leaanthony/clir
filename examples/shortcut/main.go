package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Shortcuts", "A basic example", "v0.0.1")

	// Set long description
	cli.LongDescription("This app shows creating shortcuts for commands and flags")

	// Using a subcommand instead of a flag
	nameCommand := cli.NewSubCommand("namecommand", "Shows your name!")
	// CommandShortCut is used to set a shortcut for a command. Create a shortcut for "namecommand" that will be "nc"
	nameCommand.CommandShortCut("nc")
	var newName string
	newNameFlag := nameCommand.StringFlag("newname", "New Name", &newName)
	newNameFlag.FlagShortCut("newname", "nn")
	newNameFlag.FlagRequired("newname")
	nameCommand.Action(func() error {
		fmt.Println("The `namecommand` command was run!, other args: ", nameCommand.OtherArgs())
		if len(nameCommand.OtherArgs()) == 0 {
			fmt.Println("A name must be supplied to the command!")
			return nil
		}

		if newName != "" {
			fmt.Println("The flag `newname` was parsed! Your new name is: ", nameCommand.OtherArgs()[0])
		}
		return nil
	})

	// Run!
	cli.Run()

}
