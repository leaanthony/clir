package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

type EmbeddedFlags struct {
	Address string `flag:"address" description:"The address of the person" default:"123 Main Street"`
}

type AppFlags struct {
	EmbeddedFlags
	Name string `default:"Bob"`
	Age  int    `default:"20"`
}

func main() {

	// Create new cli
	cli := clir.NewCli("Flags", "An example of using flags", "v0.0.1")

	cli.NewSubCommandFunction("create", "Create a new person", createPerson)
	cli.Run()

}

func createPerson(flags *AppFlags) error {
	fmt.Println("Name:", flags.Name)
	fmt.Println("Age:", flags.Age)
	fmt.Println("Address:", flags.Address)
	return nil
}
