package main

import (
	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("subcommandinheritflags", "An example of subcommands with flag inherence", "v0.0.1")

	inheritFlag := false
	cli.BoolFlag("inherit", "flag to inherit", &inheritFlag)

	// Create an init subcommand with flag inheritance
	init := cli.NewSubCommandInheritFlags("init", "Initialise the app")
	init.Action(func() error {
		println("I am initializing!", "inherit flag:", inheritFlag)
		return nil
	})

	// Run!
	if err := cli.Run(); err != nil {
		panic(err)
	}

}
