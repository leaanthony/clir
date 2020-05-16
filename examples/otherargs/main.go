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

	// Run!
	cli.Run()

}
