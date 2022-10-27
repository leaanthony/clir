package main

import (
	"github.com/leaanthony/clir"
)

type Flags struct {
	Name string `name:"name" description:"The name of the person"`
	Age  int    `name:"age" description:"The age of the person"`
}

func main() {

	// Create new cli
	cli := clir.NewCli("flagstruct", "An example of subcommands with flag inherence", "v0.0.1")

	// Create an init subcommand with flag inheritance
	init := cli.NewSubCommand("create", "Create a person")
	person := &Flags{
		Age: 20,
	}
	init.AddFlags(person)
	init.Action(func() error {
		println("Name:", person.Name, "Age:", person.Age)
		return nil
	})

	// Run!
	if err := cli.Run(); err != nil {
		panic(err)
	}

}
