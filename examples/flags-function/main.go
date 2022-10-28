package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

type AppFlags struct {
	Name string
	Age  int
}

// Default is an optional method that provides the default values for the flags
func (t AppFlags) Default() *AppFlags {
	return &AppFlags{
		Name: "Bob",
		Age:  20,
	}
}

func main() {

	// Create new cli
	cli := clir.NewCli("Flags", "An example of using flags", "v0.0.1")

	cli.NewSubCommandFunction("create", "Create a new person", createPerson)
	cli.Run()

}

func createPerson(flags *AppFlags) error {
	fmt.Printf("%+v\n", flags)
	return nil
}
