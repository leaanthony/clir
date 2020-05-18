package main

import (
	"fmt"

	"github.com/deranjer/clir"
)

// Required flag errors are parsed before and other errors

// Example commands to run to test functionality
// go run main.go add file.txt // Shows error since required flag '-force' is required
// go run main.go add all //  Shows error since the add command requires '-force' before the subcommand is evaluated
// go run main.go add all -force // Shows error since the flag is on the 'all' subcommand, which does not have a '-force' flag
// go run main.go add file.txt -force // Shows success since the required flag can be in any position in the args
// go run main.go add -wildcard *.txt // Shows error since the required flag '-force' is not supplied
// go run main.go add -force -wildcard *.txt // Shows success

// Uncomment the 'add' command required line to run these tests
// go run main.go add // Shows error since the 'all' command is required
// go run main.go add all // Shows error since the '-force' flag is still required on the 'add' command
// go run main.go add -force all // Shows success since the -force flag is ONLY required if the 'add' command is the final command, but 'all' is the final command

func main() {

	// Create new cli
	cli := clir.NewCli("Other Args", "A basic example", "v0.0.1")

	// Set long description
	cli.LongDescription("This app shows required flags/commands")

	// Adding a subcommand
	nameCommand := cli.NewSubCommand("add", "add a file/files")

	// Adding two flags, one required
	var wildcard bool
	var forced bool

	// wildcard is a bool flag that is not required
	nameCommand.BoolFlag("wildcard", "will treat as wildcard", &wildcard)

	// force is a bool flag that is required
	forcedFlag := nameCommand.BoolFlag("force", "will force add", &forced)
	forcedFlag.FlagRequired("force")

	// Adding the Action to handle all the flags
	// Adding an action that will evaluate any of the supplied arguments and flags to the "add" command
	nameCommand.Action(func() error {

		// If the "wildcard"/"w" flag was supplied, treat it as a wildcard, otherwise treat it as a file
		// Just grabbing the first argument after (if one supplied) and printing it, see "advanced" for a more complete solution
		if wildcard {
			fmt.Println("Treating as wildcard: ", nameCommand.OtherArgs()[0])
		} else {
			fmt.Println("Adding File: ", nameCommand.OtherArgs()[0])
		}

		// Returning nil since everything is done
		return nil
	})

	// Showing a required subcommand instead of a flag.
	allCommand := nameCommand.NewSubCommand("all", "Required Subcommand")

	// NOTE: Uncomment this to make this subcommand required
	//allCommand.CommandRequired()

	allCommand.Action(func() error {

		// Printing out response when command run
		fmt.Println("Adding all files!!")
		return nil
	})

	// Run!
	cli.Run()

}
