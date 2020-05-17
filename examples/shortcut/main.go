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

	// Name
	var name string
	nameFlag := cli.StringFlag("name", "Your name", &name)
	// FlagShort is used to set a shortcut for a flag. Create a shortcut for "-name", that will be "-n"
	nameFlag.FlagShortCut("name", "n")
	// Define action
	cli.Action(func() error {
		println("Your name is", name)
		return nil
	})
	// Using a subcommand instead of a flag
	nameCommand := cli.NewSubCommand("namecommand", "Shows your name!")
	// CommandShortCut is used to set a shortcut for a command. Create a shortcut for "namecommand" that will be "nc"
	nameCommand.CommandShortCut("nc")
	nameCommand.Action(func() error {
		fmt.Println("The `namecommand` command was run!")
		return nil
	})
	var newName string
	newNameFlag := nameCommand.StringFlag("newname", "New Name", &newName)
	//newNameFlag.FlagShortCut("newname", "nn")
	//newNameFlag.FlagRequired("newname")
	newNameFlag.Action(func() error {
		fmt.Println("The flag `newname` was parsed!")
		return nil
	})

	// Run!
	cli.Run()

}
