package main

import (
	"fmt"
	"os"

	"github.com/leaanthony/clir"
)

// Example commands to run to test functionality
// go run main.go add all test // Shows error and exiting on excess input
// go run main.go add all  // Shows success for adding all files
// go run main.go add file.txt // Shows error when trying to run without a required flag
// go run main.go add -force file.txt // Show success when adding 'file.txt'.
// go run main.go add -f file.txt // Show success when adding 'file.txt' with shortcut for force
// go run main.go add -f -w *.txt // Show success when adding a wildcard (validation isn't currently done but could be added)
// go run main.go add -w *.txt // Show error again that the force flag is ALWAYS required
// go run main.go add file.txt -f // Show error that the flags MUST come before arguments
// go run main.go add -Force file.txt // Show error that case is sensitive

//Uncomment the 'secondforced' flag and required statement to test multiple forced flags
// go run main.go add -f -secondforced file.txt  // Show success when supplying both forced flags
// go run main.go add -f file.txt // Show error when not supplying the second required forced flag

func main() {

	// Create new cli
	cli := clir.NewCli("Other Args", "A basic example", "v0.0.1")

	// Set long description
	cli.LongDescription("This is a more advanced examples using a lot of features and handling response")

	// Setup our first subcommand "add"
	nameCommand := cli.NewSubCommand("add", "add a file/files")

	// Adding two flags, one is required and the other has shortcut
	var wildcard bool
	var forced bool

	// Adding a "wildcard" command, and giving it a shortcut
	wildcardFlag := nameCommand.BoolFlag("wildcard", "will treat as wildcard", &wildcard)
	wildcardFlag.FlagShortCut("wildcard", "w")

	// Adding a "force", giving it a shortcut and making it required
	forcedFlag := nameCommand.BoolFlag("force", "will force add", &forced)
	forcedFlag.FlagShortCut("force", "f")
	forcedFlag.FlagRequired("force")

	var secondforce bool
	secondforceFlag := nameCommand.BoolFlag("secondforced", "a second forced flag", &secondforce)
	secondforceFlag.FlagRequired("secondforced")

	// Adding an action that will evaluate any of the supplied arguments and flags to the "add" command
	nameCommand.Action(func() error {

		// Showing the user the complete list of arguments that were given
		fmt.Printf("The remaining arguments were: %+v\n", nameCommand.OtherArgs())

		// Checking to see if no other args were added
		if len(nameCommand.OtherArgs()) == 0 {
			nameCommand.PrintHelp()
			fmt.Println("ERROR: No file/folder name supplied!")
			os.Exit(0)
		}

		// If the "wildcard"/"w" flag was supplied, treat it as a wildcard, otherwise treat it as a file
		if wildcard {
			fmt.Println("Treating as wildcard: ", nameCommand.OtherArgs()[0])
		} else {
			fmt.Println("Adding File: ", nameCommand.OtherArgs()[0])
		}

		// Returning nil since everything is done
		return nil
	})

	// Adding a second subcommand, so this is a command that can be run AFTER "add".  Commands do not have "-" in front of them.
	secondCommand := nameCommand.NewSubCommand("all", "adds all files")
	secondCommand.Action(func() error {

		// We can choose to ignore anything after the subcommand if we want.. or we can handle it!
		if len(secondCommand.OtherArgs()) > 0 {
			//fmt.Println("No other flags or subcommands accepted, ignoring: ", secondCommand.OtherArgs())
			fmt.Println("No other flags or subcommands accepted after 'all', exiting...")
			os.Exit(1)
		}
		fmt.Println("Adding all files")
		return nil
	})

	// Run!
	cli.Run()

}
