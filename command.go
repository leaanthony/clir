package clir

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
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
	positionalArgsMap map[string]reflect.Value
}

// NewCommand creates a new Command
// func NewCommand(name string, description string, app *Cli, parentCommandPath string) *Command {
func NewCommand(name string, description string) *Command {
	result := &Command{
		name:              name,
		shortdescription:  description,
		subCommandsMap:    make(map[string]*Command),
		hidden:            false,
		positionalArgsMap: make(map[string]reflect.Value),
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
	defer func() {
		os.Stderr = tmp
	}()

	// Credit: https://stackoverflow.com/a/74146375
	var positionalArgs []string
	for {
		if err := c.flags.Parse(args); err != nil {
			return err
		}
		// Consume all the flags that were parsed as flags.
		args = args[len(args)-c.flags.NArg():]
		if len(args) == 0 {
			break
		}
		// There's at least one flag remaining and it must be a positional arg since
		// we consumed all args that were parsed as flags. Consume just the first
		// one, and retry parsing, since subsequent args may be flags.
		positionalArgs = append(positionalArgs, args[0])
		args = args[1:]
	}

	// Parse just the positional args so that flagset.Args()/flagset.NArgs()
	// return the expected value.
	// Note: This should never return an error.
	err := c.flags.Parse(positionalArgs)
	if err != nil {
		return err
	}

	if len(positionalArgs) > 0 {
		return c.parsePositionalArgs(positionalArgs)
	}
	return nil
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
		defaultValue := tag.Get("default")
		pos := tag.Get("pos")
		c.positionalArgsMap[pos] = field
		if name == "" {
			name = strings.ToLower(t.Elem().Field(i).Name)
		}
		switch field.Kind() {
		case reflect.Bool:
			var defaultValueBool bool
			if defaultValue != "" {
				var err error
				defaultValueBool, err = strconv.ParseBool(defaultValue)
				if err != nil {
					panic("Invalid default value for bool flag")
				}
			}
			field.SetBool(defaultValueBool)
			c.BoolFlag(name, description, field.Addr().Interface().(*bool))
		case reflect.String:
			if defaultValue != "" {
				// set value of field to default value
				field.SetString(defaultValue)
			}
			c.StringFlag(name, description, field.Addr().Interface().(*string))
		case reflect.Int:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int flag")
				}
				field.SetInt(int64(value))
			}
			c.IntFlag(name, description, field.Addr().Interface().(*int))
		case reflect.Int64:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int flag")
				}
				field.SetInt(int64(value))
			}
			c.Int64Flag(name, description, field.Addr().Interface().(*int64))
		case reflect.Uint:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int flag")
				}
				field.SetUint(uint64(value))
			}
			c.UintFlag(name, description, field.Addr().Interface().(*uint))
		case reflect.Uint64:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int flag")
				}
				field.SetUint(uint64(value))
			}
			c.UInt64Flag(name, description, field.Addr().Interface().(*uint64))
		case reflect.Float64:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.ParseFloat(defaultValue, 64)
				if err != nil {
					panic("Invalid default value for float flag")
				}
				field.SetFloat(value)
			}
			c.Float64Flag(name, description, field.Addr().Interface().(*float64))
		default:
			if pos != "" {
				println("WARNING: Unsupported type for flag: ", fieldType.Type.Kind())
			}
		}
	}

	return c
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

// Int64Flag - Adds an int flag to the command
func (c *Command) Int64Flag(name, description string, variable *int64) *Command {
	c.flags.Int64Var(variable, name, *variable, description)
	c.flagCount++
	return c
}

// UintFlag - Adds an int flag to the command
func (c *Command) UintFlag(name, description string, variable *uint) *Command {
	c.flags.UintVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// UInt64Flag - Adds an int flag to the command
func (c *Command) UInt64Flag(name, description string, variable *uint64) *Command {
	c.flags.Uint64Var(variable, name, *variable, description)
	c.flagCount++
	return c
}

// Float64Flag - Adds a float64 flag to the command
func (c *Command) Float64Flag(name, description string, variable *float64) *Command {
	c.flags.Float64Var(variable, name, *variable, description)
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

func (c *Command) parsePositionalArgs(args []string) error {
	for index, posArg := range args {
		// Check the map for a field for this arg
		key := strconv.Itoa(index + 1)
		field, ok := c.positionalArgsMap[key]
		if !ok {
			continue
		}
		fieldType := field.Type()
		switch fieldType.Kind() {
		case reflect.Bool:
			// set value of field to true
			field.SetBool(true)
		case reflect.String:
			field.SetString(posArg)
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			value, err := strconv.ParseInt(posArg, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(value)
		case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
			value, err := strconv.ParseUint(posArg, 10, 64)
			if err != nil {
				return err
			}
			field.SetUint(value)
		case reflect.Float64, reflect.Float32:
			value, err := strconv.ParseFloat(posArg, 64)
			if err != nil {
				return err
			}
			field.SetFloat(value)
		default:
			return errors.New("Unsupported type for positional argument: " + fieldType.Name())
		}
	}
	return nil
}
