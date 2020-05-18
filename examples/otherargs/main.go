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
	cli.StringFlag("name", "Your name", &name)

	// Define action
	cli.Action(func() error {
		println("Your name is", name)
		fmt.Printf("The remaining arguments were: %+v\n", cli.OtherArgs())
		return nil
	})
	// Using a subcommand instead of a flag
	nameCommand := cli.NewSubCommand("namecommand", "Shows your name!")
	nameCommand.Action(func() error {
		fmt.Printf("The remaining arguments were: %+v\n", nameCommand.OtherArgs())
		return nil
	})
	var newName string
	newNameFlag := nameCommand.StringFlag("newname", "New Name", &newName)
	newNameFlag.Action(func() error {
		fmt.Println("The flag `newname` was parsed!")
		return nil
	})

	// Run!
	cli.Run()

}
