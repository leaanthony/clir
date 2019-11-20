package main

import (
	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Nested Subcommands", "An example of nested subcommands", "v0.0.1")

	// Create a top level subcommand
	// Run with: `nested-subcommand top`
	top := cli.NewSubCommand("top", "top level command")
	top.Action(func() error {
		println("I am the top-level command!")
		return nil
	})

	// Create a middle subcommand on the top command
	// Run with: `nested-subcommand top middle`
	middle := top.NewSubCommand("middle", "middle level command")
	middle.Action(func() error {
		println("I am the middle-level command!")
		return nil
	})

	// Create a bottom subcommand on the middle command
	// Run with: `nested-subcommand top middle bottom`
	bottom := middle.NewSubCommand("bottom", "bottom level command")
	bottom.Action(func() error {
		println("I am the bottom-level command!")
		return nil
	})

	// Run!
	cli.Run()

}
