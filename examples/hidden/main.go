package main

import (
	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("subcommands", "An example of subcommands", "v0.0.1")

	// Create an init subcommand
	init := cli.NewSubCommand("init", "Initialise the app")
	init.Action(func() error {
		println("I am initialising!")
		return nil
	})

	// Create a test subcommand
	test := cli.NewSubCommand("test", "Test the app")
	test.Action(func() error {
		println("I am testing!")
		return nil
	})

	// Make test hidden
	test.Hidden()

	// Run!
	cli.Run()

}
