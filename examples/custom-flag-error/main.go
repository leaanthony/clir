package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func customFlagError(cli *clir.Cli) string {
	return `Flag v0.0.1 - A custom error example

Flags:

  --help
	Get help on the 'flag' command.`
}

func main() {

	// Create new cli
	cli := clir.NewCli("Flag", "A custom error example", "v0.0.1")

	cli.SetFlagFunction(customFlagError)

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

}
