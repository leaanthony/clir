package main

import (
	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Subcommand", "A basic example", "v0.0.1")

	// Create a subcommand
	init := cli.NewSubCommand("init", "Initialise the app")
	init.Action(func() error {
		println("Initialising!")
		return nil
	})

	// Run!
	cli.Run()

}
