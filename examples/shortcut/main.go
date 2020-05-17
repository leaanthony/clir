package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Other Args", "A basic example", "v0.0.1")

	// Set long description
	cli.LongDescription("This app shows positional arguments")

	// Name
	var name string
	nameFlag := cli.StringFlag("name", "Your name", &name)
	// FlagShort is used to set a shortcut for a flag. Create a shortcut for "-name", that will be "-n"
	nameFlag.FlagShortCut("name", "n")
	// Define action
	cli.Action(func() error {
		println("Your name is", name)
		fmt.Printf("The remaining arguments were: %+v\n", cli.OtherArgs())
		return nil
	})
	// Using a subcommand instead of a flag
	nameCommand := cli.NewSubCommand("namecommand", "Shows your name!")
	// CommandShortCut is used to set a shortcut for a command. Create a shortcut for "namecommand" that will be "nc"
	nameCommand.CommandShortCut("nc")
	nameCommand.Action(func() error {
		fmt.Printf("The remaining arguments were: %+v\n", nameCommand.OtherArgs())
		return nil
	})

	// Run!
	cli.Run()

}
