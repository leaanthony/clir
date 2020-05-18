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

	// Now adding a StringFlag that also has a shortcut, FlagShortCut requires the longname, then the short name
	var newName string
	newNameFlag := nameCommand.StringFlag("newname", "New Name", &newName)
	newNameFlag.FlagShortCut("newname", "nn")

	// Adding an action for the 'namecommand'/'nn' command.
	nameCommand.Action(func() error {
		fmt.Println("The `namecommand` command was run!")

		// Handling the flag
		if newName != "" {
			fmt.Println("The flag `newname` was parsed! Your new name is: ", nameCommand.OtherArgs()[0])
		}
		return nil
	})

	// Adding another command with another shortcut
	locationCommand := nameCommand.NewSubCommand("location", "set your location!")
	// Adding the shortcut 'loc' for 'location'
	locationCommand.CommandShortCut("loc")
	// Adding the action for the 'location' command
	locationCommand.Action(func() error {
		fmt.Println("The `location` / `loc` command was run!")
		return nil
	})

	// Run!
	cli.Run()

}
