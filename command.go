package clir

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Command represents a command that may be run by the user
type Command struct {
	name              string
	commandPath       string
	shortdescription  string
	shortCut          string
	longdescription   string
	subCommands       []*Command
	subCommandsMap    map[string]*Command
	longestSubcommand int
	actionCallback    Action
	app               *Cli
	flags             *flag.FlagSet
	flagList          []flagDetails
	flagCount         int
	helpFlag          bool
	hidden            bool
}

type flagDetails struct {
	flagName string
	shortCut string
}

// NewCommand creates a new Command
// func NewCommand(name string, description string, app *Cli, parentCommandPath string) *Command {
func NewCommand(name string, description string) *Command {
	result := &Command{
		name:             name,
		shortdescription: description,
		subCommandsMap:   make(map[string]*Command),
		hidden:           false,
	}

	return result
}

func (c *Command) setParentCommandPath(parentCommandPath string) {
	// Set up command path
	if parentCommandPath != "" {
		c.commandPath += parentCommandPath + " "
	}
	c.commandPath += c.name

	// Set up flag set
	c.flags = flag.NewFlagSet(c.commandPath, flag.ContinueOnError)
	c.BoolFlag("help", "Get help on the '"+strings.ToLower(c.commandPath)+"' command.", &c.helpFlag)

	// result.Flags.Usage = result.PrintHelp

}

func (c *Command) setApp(app *Cli) {
	c.app = app
}

// parseFlags parses the given flags
func (c *Command) parseFlags(args []string) error {
	// Parse flags
	tmp := os.Stderr
	os.Stderr = nil
	err := c.flags.Parse(args)
	os.Stderr = tmp
	return err
}

// Run - Runs the Command with the given arguments
func (c *Command) run(args []string) error {
	// If we have arguments, process them
	if len(args) > 0 {
		// Convert command shortCut to full command name
		for _, command := range c.subCommands {
			if args[0] == command.shortCut {
				args[0] = command.name
			}
		}
		// Check for subcommand
		subcommand := c.subCommandsMap[args[0]]
		if subcommand != nil {
			return subcommand.run(args[1:])
		}
		// Convert flag shortCut to full flag
		for _, arg := range args {
			if arg[0] == '-' {
				for _, flagDetails := range c.flagList {
					if flagDetails.shortCut == arg[1:] {
						args[0] = "-" + flagDetails.flagName
					}
				}
			}
		}

		// Parse flags
		err := c.parseFlags(args)
		if err != nil {
			fmt.Printf("Error: %s\n\n", err.Error())
			c.PrintHelp()
			return err
		}

		// Help takes precedence
		if c.helpFlag {
			c.PrintHelp()
			return nil
		}
	}

	// Do we have an action?
	if c.actionCallback != nil {
		return c.actionCallback()
	}

	// If we haven't specified a subcommand
	// check for an app level default command
	if c.app.defaultCommand != nil {
		// Prevent recursion!
		if c.app.defaultCommand != c {
			// only run default command if no args passed
			if len(args) == 0 {
				return c.app.defaultCommand.run(args)
			}
		}
	}

	// Nothing left we can do
	c.PrintHelp()

	return nil
}

// Action - Define an action from this command
func (c *Command) Action(callback Action) *Command {
	c.actionCallback = callback
	return c
}

// PrintHelp - Output the help text for this command
func (c *Command) PrintHelp() {
	c.app.PrintBanner()

	commandTitle := c.commandPath
	if c.shortdescription != "" {
		commandTitle += " - " + c.shortdescription
	}
	// Ignore root command
	if c.commandPath != c.name {
		fmt.Println(commandTitle)
	}
	if c.longdescription != "" {
		fmt.Println(c.longdescription + "\n")
	}
	if c.shortCut != "" {
		fmt.Println(c.shortCut + "\n")
	}
	if len(c.subCommands) > 0 {
		fmt.Println("Available commands:")
		fmt.Println("")
		for _, subcommand := range c.subCommands {
			if subcommand.isHidden() {
				continue
			}
			spacer := strings.Repeat(" ", 3+c.longestSubcommand-len(subcommand.name))
			isDefault := ""
			if subcommand.isDefaultCommand() {
				isDefault = "[default]"
			}
			fmt.Printf("   %s%s%s%s%s %s\n", subcommand.name, spacer, subcommand.shortCut, spacer, subcommand.shortdescription, isDefault)
		}
		fmt.Println("")
	}
	if c.flagCount > 0 {
		fmt.Println("Flags:")
		fmt.Println()
		c.flags.SetOutput(os.Stdout)
		c.flags.PrintDefaults()
		c.flags.SetOutput(os.Stderr)

	}
	fmt.Println()
}

// isDefaultCommand returns true if called on the default command
func (c *Command) isDefaultCommand() bool {
	return c.app.defaultCommand == c
}

// isHidden returns true if the command is a hidden command
func (c *Command) isHidden() bool {
	return c.hidden
}

// Hidden hides the command from the Help system
func (c *Command) Hidden() {
	c.hidden = true
}

// NewSubCommand - Creates a new subcommand
func (c *Command) NewSubCommand(name, description string) *Command {
	result := NewCommand(name, description)
	c.AddCommand(result)
	return result
}

// AddCommand - Adds a subcommand
func (c *Command) AddCommand(command *Command) {
	command.setApp(c.app)
	command.setParentCommandPath(c.commandPath)
	name := command.name
	c.subCommands = append(c.subCommands, command)
	c.subCommandsMap[name] = command
	if len(name) > c.longestSubcommand {
		c.longestSubcommand = len(name)
	}
}

// BoolFlag - Adds a boolean flag to the command
func (c *Command) BoolFlag(name, description string, variable *bool) *Command {
	c.flags.BoolVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// StringFlag - Adds a string flag to the command
func (c *Command) StringFlag(name, description string, variable *string) *Command {
	c.flags.StringVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// IntFlag - Adds an int flag to the command
func (c *Command) IntFlag(name, description string, variable *int) *Command {
	c.flags.IntVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// FlagShortCut  creates a shortcut or shorter call to a flags (i.e. "-readfiles" or "-rf")
func (c *Command) FlagShortCut(flagLongName string, flagShortCut string) *Command {
	for _, singleFlag := range c.flagList {
		if singleFlag.flagName == flagLongName {
			singleFlag.shortCut = flagShortCut
		}
	}
	return c
}

// CommandShortCut creates a shortcut or shorter call to a command (i.e. "readfiles" or "rf")
func (c *Command) CommandShortCut(cmdShortCut string) *Command {
	c.shortCut = cmdShortCut
	return c
}

// LongDescription - Sets the long description for the command
func (c *Command) LongDescription(longdescription string) *Command {
	c.longdescription = longdescription
	return c
}

// OtherArgs - Returns the non-flag arguments passed to the subcommand. NOTE: This should only be called within the context of an action.
func (c *Command) OtherArgs() []string {
	return c.flags.Args()
}
