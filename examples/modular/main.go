/*
 * An example of passing an argument to a function
 * not defined within the call to cli.Action
 */

package main

import (
	"github.com/leaanthony/clir"
	"fmt"
)


func printMessage(name *string) error {
	fmt.Printf("Hello %s!\n", *name)
	return nil
}


func main() {

	// Create new cli
	cli := clir.NewCli("Basic", "A basic example", "v0.0.1")

	// Set long description
	cli.LongDescription("This app prints hello world")

	// Name
	var name string
	cli.StringFlag("name", "Your name", &name)

	// Define action
	cli.CustomAction(printMessage, &name)

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Println(err)
	}
}