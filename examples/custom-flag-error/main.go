package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func customFlagError(err error) error {
	return fmt.Errorf(`%s 
Flag v0.0.1 - A custom error example

Flags:

  --help
	Get help on the 'flag' command`, err)
}

func main() {

	// Create new cli
	cli := clir.NewCli("Flag", "A custom error example", "v0.0.1")

	cli.SetErrorFunction(customFlagError)

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Println(err)
	}

}
