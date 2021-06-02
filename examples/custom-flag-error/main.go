package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func customFlagError(cmdPath string, err error) error {
	return fmt.Errorf(`%s 
flag v0.0.1 - A custom error example

Flags:

  --help
	Get help on the '%s' command`, err, cmdPath)
}

func main() {

	// Create new cli
	cli := clir.NewCli("flag", "A custom error example", "v0.0.1")

	cli.SetErrorFunction(customFlagError)

	cli.NewSubCommand("test", "Testing whether subcommand path returns correctly via err callback")

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Println(err)
	}

}
