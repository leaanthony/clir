// Package clir provides a simple API for creating command line apps
package clir

import (
	"fmt"
	"log"
	"os"

	"github.com/leaanthony/clir/internal/colour"
)

// NewCli - Creates a new Cli application object
func NewCli(name, description, version string) *Cli {
	result := &Cli{
		version: version,
	}
	result.rootCommand = NewCommand(name, description)
	result.rootCommand.setApp(result)
	result.rootCommand.setParentCommandPath("")
	return result
}

// Cli - The main application object
type Cli struct {
	version        string
	rootCommand    *Command
	defaultCommand *Command
	preRunCommand  func(*Cli) error
}

// Version - Get the Application version string
func (c *Cli) Version() string {
	return c.version
}

// Abort prints the given error and terminates the application
func (c *Cli) Abort(err error) {
	log.Fatal(err)
	os.Exit(1)
}

// AddCommand - Adds a command to the application
func (c *Cli) AddCommand(command *Command) {
	c.rootCommand.addCommand(command)
}

// PrintBanner prints the application banner!
func (c *Cli) PrintBanner() {
	fmt.Printf("%s %s - %s\n\n", colour.YellowString(c.rootCommand.name), colour.RedString(c.version), c.rootCommand.shortdescription)
}

// PrintHelp - Prints the application's help
func (c *Cli) PrintHelp() {
	c.rootCommand.PrintHelp()
}

// Run - Runs the application with the given arguments
func (c *Cli) Run(args ...string) error {
	if c.preRunCommand != nil {
		err := c.preRunCommand(c)
		if err != nil {
			return err
		}
	}
	if len(args) == 0 {
		args = os.Args[1:]
	}
	return c.rootCommand.run(args)
}

// DefaultCommand - Sets the given command as the command to run when
// no other commands given
func (c *Cli) DefaultCommand(defaultCommand *Command) *Cli {
	c.defaultCommand = defaultCommand
	return c
}

// Command - Adds a command to the application
func (c *Cli) Command(name, description string) *Command {
	return c.rootCommand.NewSubCommand(name, description)
}

// PreRun - Calls the given function before running the specific command
func (c *Cli) PreRun(callback func(*Cli) error) {
	c.preRunCommand = callback
}

// BoolFlag - Adds a boolean flag to the root command
func (c *Cli) BoolFlag(name, description string, variable *bool) *Command {
	c.rootCommand.BoolFlag(name, description, variable)
	return c.rootCommand
}

// StringFlag - Adds a string flag to the root command
func (c *Cli) StringFlag(name, description string, variable *string) *Command {
	c.rootCommand.StringFlag(name, description, variable)
	return c.rootCommand
}
