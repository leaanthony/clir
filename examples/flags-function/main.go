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
	Name string `default:"Bob" pos:"1"`
	Age  int    `default:"20" pos:"2"`
}

func main() {

	// Create new cli
	cli := clir.NewCli("Flags", "An example of using flags", "v0.0.1")

	cli.NewSubCommandFunction("create", "Create a new person", createPerson)
	cli.Run("create", "bob", "20", "--address", "123 Main Street")
}

func createPerson(flags *AppFlags) error {
	fmt.Println("Name:", flags.Name)
	fmt.Println("Age:", flags.Age)
	fmt.Println("Address:", flags.Address)
	return nil
}
