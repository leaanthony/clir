package main

import (
	"context"
	"fmt"

	"github.com/leaanthony/clir"
)

func customFlagError(ctx context.Context, err error) error {
	return fmt.Errorf(`%s 
flag v0.0.1 - A custom error example

Flags:

  --help
	Get help on the '%s' command`, err, ctx.Value("path"))
}

func main() {

	// Create new cli
	cli := clir.NewCli("flag", "A custom error example", "v0.0.1")

	cli.SetErrorFunction(customFlagError)

	cli.NewSubCommand("test", "Testing whether subcommands return via context to err")

	// Run!
	if err := cli.Run(); err != nil {
		fmt.Println(err)
	}

}
