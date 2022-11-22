package clir

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Command represents a command that may be run by the user
type Command struct {
	name              string
	commandPath       string
	shortdescription  string
	longdescription   string
	subCommands       []*Command
	subCommandsMap    map[string]*Command
	longestSubcommand int
	actionCallback    Action
	app               *Cli
	flags             *flag.FlagSet
	flagCount         int
	helpFlag          bool
	hidden            bool
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

func (c *Command) inheritFlags(inheritFlags *flag.FlagSet) {
	// inherit flags
	inheritFlags.VisitAll(func(f *flag.Flag) {
		if f.Name != "help" {
			c.flags.Var(f.Value, f.Name, f.Usage)
		}
	})
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
		// Check for subcommand
		subcommand := c.subCommandsMap[args[0]]
		if subcommand != nil {
			return subcommand.run(args[1:])
		}

		// Parse flags
		err := c.parseFlags(args)
		if err != nil {
			if c.app.errorHandler != nil {
				return c.app.errorHandler(c.commandPath, err)
			}
			return fmt.Errorf("Error: %s\nSee '%s --help' for usage", err, c.commandPath)
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
			fmt.Printf("   %s%s%s %s\n", subcommand.name, spacer, subcommand.shortdescription, isDefault)
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

// NewSubCommandInheritFlags - Creates a new subcommand, inherits flags from command
func (c *Command) NewSubCommandInheritFlags(name, description string) *Command {
	result := c.NewSubCommand(name, description)
	result.inheritFlags(c.flags)
	return result
}

// BoolFlag - Adds a boolean flag to the command
func (c *Command) BoolFlag(name, description string, variable *bool) *Command {
	c.flags.BoolVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

func (c *Command) AddFlags(optionStruct interface{}) *Command {
	// use reflection to determine if this is a pointer to a struct
	// if not, panic

	t := reflect.TypeOf(optionStruct)

	// Check for a pointer to a struct
	if t.Kind() != reflect.Ptr {
		panic("AddFlags() requires a pointer to a struct")
	}
	if t.Elem().Kind() != reflect.Struct {
		panic("AddFlags() requires a pointer to a struct")
	}

	// Iterate through the fields of the struct reading the struct tags
	// and adding the flags
	v := reflect.ValueOf(optionStruct).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Elem().Field(i)
		if !fieldType.IsExported() {
			continue
		}
		// If this is an embedded struct, recurse
		if fieldType.Type.Kind() == reflect.Struct {
			c.AddFlags(field.Addr().Interface())
			continue
		}

		tag := t.Elem().Field(i).Tag
		name := tag.Get("name")
		description := tag.Get("description")
		if name == "" {
			name = strings.ToLower(t.Elem().Field(i).Name)
		}
		switch field.Kind() {
		case reflect.Bool:
			c.BoolFlag(name, description, field.Addr().Interface().(*bool))
		case reflect.String:
			c.StringFlag(name, description, field.Addr().Interface().(*string))
		case reflect.Int:
			c.IntFlag(name, description, field.Addr().Interface().(*int))
		}
	}

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

// LongDescription - Sets the long description for the command
func (c *Command) LongDescription(longdescription string) *Command {
	c.longdescription = longdescription
	return c
}

// OtherArgs - Returns the non-flag arguments passed to the subcommand. NOTE: This should only be called within the context of an action.
func (c *Command) OtherArgs() []string {
	return c.flags.Args()
}

func (c *Command) NewSubCommandFunction(name string, description string, fn interface{}) *Command {
	result := c.NewSubCommand(name, description)
	// use reflection to determine if this is a function
	// if not, panic
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}

	// Check the function has 1 input ant it's a struct pointer
	fnValue := reflect.ValueOf(fn)
	if t.NumIn() != 1 {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	// Check the input is a struct pointer
	if t.In(0).Kind() != reflect.Ptr {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	if t.In(0).Elem().Kind() != reflect.Struct {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	// Check only 1 output and it's an error
	if t.NumOut() != 1 {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	if t.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("NewSubFunction '" + name + "' requires a function with the signature 'func(*struct) error'")
	}
	flags := reflect.New(t.In(0).Elem())
	defaultMethod, ok := t.In(0).MethodByName("Default")

	if ok {
		// Check the default method has no inputs
		if defaultMethod.Type.NumIn() != 1 {
			panic("'Default' method on struct '" + t.In(0).Elem().Name() + "' must have the signature 'Default() *" + t.In(0).Elem().Name() + "'")
		}

		// Check the default method has a single struct output
		if defaultMethod.Type.NumOut() != 1 {
			panic("'Default' method on struct '" + t.In(0).Elem().Name() + "' must have the signature 'Default() *" + t.In(0).Elem().Name() + "'")
		}

		// Check the default method has a single struct output
		if defaultMethod.Type.Out(0) != t.In(0) {
			panic("'Default' method on struct '" + t.In(0).Elem().Name() + "' must have the signature 'Default() *" + t.In(0).Elem().Name() + "'")
		}

		// Call defaultMethod to get default flags
		results := defaultMethod.Func.Call([]reflect.Value{flags})
		flags = results[0]
	}
	result.Action(func() error {
		result := fnValue.Call([]reflect.Value{flags})[0].Interface()
		if result != nil {
			return result.(error)
		}
		return nil
	})
	result.AddFlags(flags.Interface())
	return result
}
