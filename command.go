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
	sliceSeparator    map[string]string
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
		sliceSeparator:    make(map[string]string),
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
		sep := tag.Get("sep")
		c.positionalArgsMap[pos] = field
		if sep != "" {
			c.sliceSeparator[pos] = sep
		}
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
		case reflect.Int8:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int8 flag")
				}
				field.SetInt(int64(value))
			}
			c.Int8Flag(name, description, field.Addr().Interface().(*int8))
		case reflect.Int16:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int16 flag")
				}
				field.SetInt(int64(value))
			}
			c.Int16Flag(name, description, field.Addr().Interface().(*int16))
		case reflect.Int32:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int32 flag")
				}
				field.SetInt(int64(value))
			}
			c.Int32Flag(name, description, field.Addr().Interface().(*int32))
		case reflect.Int64:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for int64 flag")
				}
				field.SetInt(int64(value))
			}
			c.Int64Flag(name, description, field.Addr().Interface().(*int64))
		case reflect.Uint:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for uint flag")
				}
				field.SetUint(uint64(value))
			}
			c.UintFlag(name, description, field.Addr().Interface().(*uint))
		case reflect.Uint8:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for uint8 flag")
				}
				field.SetUint(uint64(value))
			}
			c.Uint8Flag(name, description, field.Addr().Interface().(*uint8))
		case reflect.Uint16:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for uint16 flag")
				}
				field.SetUint(uint64(value))
			}
			c.Uint16Flag(name, description, field.Addr().Interface().(*uint16))
		case reflect.Uint32:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for uint32 flag")
				}
				field.SetUint(uint64(value))
			}
			c.Uint32Flag(name, description, field.Addr().Interface().(*uint32))
		case reflect.Uint64:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.Atoi(defaultValue)
				if err != nil {
					panic("Invalid default value for uint64 flag")
				}
				field.SetUint(uint64(value))
			}
			c.UInt64Flag(name, description, field.Addr().Interface().(*uint64))
		case reflect.Float32:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.ParseFloat(defaultValue, 64)
				if err != nil {
					panic("Invalid default value for float32 flag")
				}
				field.SetFloat(value)
			}
			c.Float32Flag(name, description, field.Addr().Interface().(*float32))
		case reflect.Float64:
			if defaultValue != "" {
				// set value of field to default value
				value, err := strconv.ParseFloat(defaultValue, 64)
				if err != nil {
					panic("Invalid default value for float64 flag")
				}
				field.SetFloat(value)
			}
			c.Float64Flag(name, description, field.Addr().Interface().(*float64))
		case reflect.Slice:
			c.addSliceField(field, defaultValue, sep)
			c.addSliceFlags(name, description, field)
		default:
			if pos != "" {
				println("WARNING: Unsupported type for flag: ", fieldType.Type.Kind(), name)
			}
		}
	}

	return c
}

func (c *Command) addSliceFlags(name, description string, field reflect.Value) *Command {
	if field.Kind() != reflect.Slice {
		panic("addSliceFlags() requires a pointer to a slice")
	}
	t := reflect.TypeOf(field.Addr().Interface())
	if t.Kind() != reflect.Ptr {
		panic("addSliceFlags() requires a pointer to a slice")
	}
	if t.Elem().Kind() != reflect.Slice {
		panic("addSliceFlags() requires a pointer to a slice")
	}
	switch t.Elem().Elem().Kind() {
	case reflect.Bool:
		c.BoolsFlag(name, description, field.Addr().Interface().(*[]bool))
	case reflect.String:
		c.StringsFlag(name, description, field.Addr().Interface().(*[]string))
	case reflect.Int:
		c.IntsFlag(name, description, field.Addr().Interface().(*[]int))
	case reflect.Int8:
		c.Int8sFlag(name, description, field.Addr().Interface().(*[]int8))
	case reflect.Int16:
		c.Int16sFlag(name, description, field.Addr().Interface().(*[]int16))
	case reflect.Int32:
		c.Int32sFlag(name, description, field.Addr().Interface().(*[]int32))
	case reflect.Int64:
		c.Int64sFlag(name, description, field.Addr().Interface().(*[]int64))
	case reflect.Uint:
		c.UintsFlag(name, description, field.Addr().Interface().(*[]uint))
	case reflect.Uint8:
		c.Uint8sFlag(name, description, field.Addr().Interface().(*[]uint8))
	case reflect.Uint16:
		c.Uint16sFlag(name, description, field.Addr().Interface().(*[]uint16))
	case reflect.Uint32:
		c.Uint32sFlag(name, description, field.Addr().Interface().(*[]uint32))
	case reflect.Uint64:
		c.Uint64sFlag(name, description, field.Addr().Interface().(*[]uint64))
	case reflect.Float32:
		c.Float32sFlag(name, description, field.Addr().Interface().(*[]float32))
	case reflect.Float64:
		c.Float64sFlag(name, description, field.Addr().Interface().(*[]float64))
	default:
		panic(fmt.Sprintf("addSliceFlags() not supported slice type %s", t.Elem().Elem().Kind().String()))
	}
	return c
}

func (c *Command) addSliceField(field reflect.Value, defaultValue, separator string) *Command {
	if defaultValue == "" {
		return c
	}
	if field.Kind() != reflect.Slice {
		panic("addSliceField() requires a pointer to a slice")
	}
	t := reflect.TypeOf(field.Addr().Interface())
	if t.Kind() != reflect.Ptr {
		panic("addSliceField() requires a pointer to a slice")
	}
	if t.Elem().Kind() != reflect.Slice {
		panic("addSliceField() requires a pointer to a slice")
	}
	defaultSlice := []string{defaultValue}
	if separator != "" {
		defaultSlice = strings.Split(defaultValue, separator)
	}
	switch t.Elem().Elem().Kind() {
	case reflect.Bool:
		defaultValues := make([]bool, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.ParseBool(value)
			if err != nil {
				panic("Invalid default value for bool flag")
			}
			defaultValues = append(defaultValues, val)
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.String:
		field.Set(reflect.ValueOf(defaultSlice))
	case reflect.Int:
		defaultValues := make([]int, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for int flag")
			}
			defaultValues = append(defaultValues, val)
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Int8:
		defaultValues := make([]int8, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for int8 flag")
			}
			defaultValues = append(defaultValues, int8(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Int16:
		defaultValues := make([]int16, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for int16 flag")
			}
			defaultValues = append(defaultValues, int16(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Int32:
		defaultValues := make([]int32, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				panic("Invalid default value for int32 flag")
			}
			defaultValues = append(defaultValues, int32(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Int64:
		defaultValues := make([]int64, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				panic("Invalid default value for int64 flag")
			}
			defaultValues = append(defaultValues, val)
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Uint:
		defaultValues := make([]uint, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for uint flag")
			}
			defaultValues = append(defaultValues, uint(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Uint8:
		defaultValues := make([]uint8, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for uint8 flag")
			}
			defaultValues = append(defaultValues, uint8(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Uint16:
		defaultValues := make([]uint16, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for uint16 flag")
			}
			defaultValues = append(defaultValues, uint16(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Uint32:
		defaultValues := make([]uint32, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for uint32 flag")
			}
			defaultValues = append(defaultValues, uint32(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Uint64:
		defaultValues := make([]uint64, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for uint64 flag")
			}
			defaultValues = append(defaultValues, uint64(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Float32:
		defaultValues := make([]float32, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for float32 flag")
			}
			defaultValues = append(defaultValues, float32(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	case reflect.Float64:
		defaultValues := make([]float64, 0, len(defaultSlice))
		for _, value := range defaultSlice {
			val, err := strconv.Atoi(value)
			if err != nil {
				panic("Invalid default value for float64 flag")
			}
			defaultValues = append(defaultValues, float64(val))
		}
		field.Set(reflect.ValueOf(defaultValues))
	default:
		panic(fmt.Sprintf("addSliceField() not supported slice type %s", t.Elem().Elem().Kind().String()))
	}
	return c
}

// BoolFlag - Adds a boolean flag to the command
func (c *Command) BoolFlag(name, description string, variable *bool) *Command {
	c.flags.BoolVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// BoolsFlag - Adds a booleans flag to the command
func (c *Command) BoolsFlag(name, description string, variable *[]bool) *Command {
	c.flags.Var(newBoolsValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// StringFlag - Adds a string flag to the command
func (c *Command) StringFlag(name, description string, variable *string) *Command {
	c.flags.StringVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// StringsFlag - Adds a strings flag to the command
func (c *Command) StringsFlag(name, description string, variable *[]string) *Command {
	c.flags.Var(newStringsValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// IntFlag - Adds an int flag to the command
func (c *Command) IntFlag(name, description string, variable *int) *Command {
	c.flags.IntVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// IntsFlag - Adds an ints flag to the command
func (c *Command) IntsFlag(name, description string, variable *[]int) *Command {
	c.flags.Var(newIntsValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int8Flag - Adds an int8 flag to the command
func (c *Command) Int8Flag(name, description string, variable *int8) *Command {
	c.flags.Var(newInt8Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int8sFlag - Adds an int8 s flag to the command
func (c *Command) Int8sFlag(name, description string, variable *[]int8) *Command {
	c.flags.Var(newInt8sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int16Flag - Adds an int16 flag to the command
func (c *Command) Int16Flag(name, description string, variable *int16) *Command {
	c.flags.Var(newInt16Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int16sFlag - Adds an int16s flag to the command
func (c *Command) Int16sFlag(name, description string, variable *[]int16) *Command {
	c.flags.Var(newInt16sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int32Flag - Adds an int32 flag to the command
func (c *Command) Int32Flag(name, description string, variable *int32) *Command {
	c.flags.Var(newInt32Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int32sFlag - Adds an int32s flag to the command
func (c *Command) Int32sFlag(name, description string, variable *[]int32) *Command {
	c.flags.Var(newInt32sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Int64Flag - Adds an int64 flag to the command
func (c *Command) Int64Flag(name, description string, variable *int64) *Command {
	c.flags.Int64Var(variable, name, *variable, description)
	c.flagCount++
	return c
}

// Int64sFlag - Adds an int64s flag to the command
func (c *Command) Int64sFlag(name, description string, variable *[]int64) *Command {
	c.flags.Var(newInt64sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// UintFlag - Adds an uint flag to the command
func (c *Command) UintFlag(name, description string, variable *uint) *Command {
	c.flags.UintVar(variable, name, *variable, description)
	c.flagCount++
	return c
}

// UintsFlag - Adds an uints flag to the command
func (c *Command) UintsFlag(name, description string, variable *[]uint) *Command {
	c.flags.Var(newUintsValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Uint8Flag - Adds an uint8 flag to the command
func (c *Command) Uint8Flag(name, description string, variable *uint8) *Command {
	c.flags.Var(newUint8Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Uint8sFlag - Adds an uint8 s flag to the command
func (c *Command) Uint8sFlag(name, description string, variable *[]uint8) *Command {
	c.flags.Var(newUint8sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Uint16Flag - Adds an uint16 flag to the command
func (c *Command) Uint16Flag(name, description string, variable *uint16) *Command {
	c.flags.Var(newUint16Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Uint16sFlag - Adds an uint16s flag to the command
func (c *Command) Uint16sFlag(name, description string, variable *[]uint16) *Command {
	c.flags.Var(newUint16sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Uint32Flag - Adds an uint32 flag to the command
func (c *Command) Uint32Flag(name, description string, variable *uint32) *Command {
	c.flags.Var(newUint32Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Uint32sFlag - Adds an uint32s flag to the command
func (c *Command) Uint32sFlag(name, description string, variable *[]uint32) *Command {
	c.flags.Var(newUint32sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// UInt64Flag - Adds an uint64 flag to the command
func (c *Command) UInt64Flag(name, description string, variable *uint64) *Command {
	c.flags.Uint64Var(variable, name, *variable, description)
	c.flagCount++
	return c
}

// Uint64sFlag - Adds an uint64s flag to the command
func (c *Command) Uint64sFlag(name, description string, variable *[]uint64) *Command {
	c.flags.Var(newUint64sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Float64Flag - Adds a float64 flag to the command
func (c *Command) Float64Flag(name, description string, variable *float64) *Command {
	c.flags.Float64Var(variable, name, *variable, description)
	c.flagCount++
	return c
}

// Float32Flag - Adds a float32 flag to the command
func (c *Command) Float32Flag(name, description string, variable *float32) *Command {
	c.flags.Var(newFloat32Value(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Float32sFlag - Adds a float32s flag to the command
func (c *Command) Float32sFlag(name, description string, variable *[]float32) *Command {
	c.flags.Var(newFloat32sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

// Float64sFlag - Adds a float64s flag to the command
func (c *Command) Float64sFlag(name, description string, variable *[]float64) *Command {
	c.flags.Var(newFloat64sValue(*variable, variable), name, description)
	c.flagCount++
	return c
}

type boolsFlagVar []bool

func (f *boolsFlagVar) String() string { return fmt.Sprint([]bool(*f)) }

func (f *boolsFlagVar) Set(value string) error {
	if value == "" {
		*f = append(*f, false)
		return nil
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	*f = append(*f, b)
	return nil
}

func (f *boolsFlagVar) IsBoolFlag() bool {
	return true
}

func newBoolsValue(val []bool, p *[]bool) *boolsFlagVar {
	*p = val
	return (*boolsFlagVar)(p)
}

type stringsFlagVar []string

func (f *stringsFlagVar) String() string { return fmt.Sprint([]string(*f)) }

func (f *stringsFlagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func newStringsValue(val []string, p *[]string) *stringsFlagVar {
	*p = val
	return (*stringsFlagVar)(p)
}

type intsFlagVar []int

func (f *intsFlagVar) String() string { return fmt.Sprint([]int(*f)) }

func (f *intsFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, i)
	return nil
}

func newIntsValue(val []int, p *[]int) *intsFlagVar {
	*p = val
	return (*intsFlagVar)(p)
}

type int8Value int8

func newInt8Value(val int8, p *int8) *int8Value {
	*p = val
	return (*int8Value)(p)
}

func (f *int8Value) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = int8Value(i)
	return nil
}

func (f *int8Value) String() string { return fmt.Sprint(int8(*f)) }

type int8sFlagVar []int8

func (f *int8sFlagVar) String() string { return fmt.Sprint([]int8(*f)) }

func (f *int8sFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, int8(i))
	return nil
}

func newInt8sValue(val []int8, p *[]int8) *int8sFlagVar {
	*p = val
	return (*int8sFlagVar)(p)
}

type int16Value int16

func newInt16Value(val int16, p *int16) *int16Value {
	*p = val
	return (*int16Value)(p)
}

func (f *int16Value) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = int16Value(i)
	return nil
}

func (f *int16Value) String() string { return fmt.Sprint(int16(*f)) }

type int16sFlagVar []int16

func (f *int16sFlagVar) String() string { return fmt.Sprint([]int16(*f)) }

func (f *int16sFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, int16(i))
	return nil
}

func newInt16sValue(val []int16, p *[]int16) *int16sFlagVar {
	*p = val
	return (*int16sFlagVar)(p)
}

type int32Value int32

func newInt32Value(val int32, p *int32) *int32Value {
	*p = val
	return (*int32Value)(p)
}

func (f *int32Value) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = int32Value(i)
	return nil
}

func (f *int32Value) String() string { return fmt.Sprint(int32(*f)) }

type int32sFlagVar []int32

func (f *int32sFlagVar) String() string { return fmt.Sprint([]int32(*f)) }

func (f *int32sFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, int32(i))
	return nil
}

func newInt32sValue(val []int32, p *[]int32) *int32sFlagVar {
	*p = val
	return (*int32sFlagVar)(p)
}

type int64sFlagVar []int64

func (f *int64sFlagVar) String() string { return fmt.Sprint([]int64(*f)) }

func (f *int64sFlagVar) Set(value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*f = append(*f, i)
	return nil
}

func newInt64sValue(val []int64, p *[]int64) *int64sFlagVar {
	*p = val
	return (*int64sFlagVar)(p)
}

type uintsFlagVar []uint

func (f *uintsFlagVar) String() string {
	return fmt.Sprint([]uint(*f))
}

func (f *uintsFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, uint(i))
	return nil
}

func newUintsValue(val []uint, p *[]uint) *uintsFlagVar {
	*p = val
	return (*uintsFlagVar)(p)
}

type uint8FlagVar uint8

func newUint8Value(val uint8, p *uint8) *uint8FlagVar {
	*p = val
	return (*uint8FlagVar)(p)
}

func (f *uint8FlagVar) String() string {
	return fmt.Sprint(uint8(*f))
}

func (f *uint8FlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = uint8FlagVar(i)
	return nil
}

type uint8sFlagVar []uint8

func (f *uint8sFlagVar) String() string {
	return fmt.Sprint([]uint8(*f))
}

func (f *uint8sFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, uint8(i))
	return nil
}

func newUint8sValue(val []uint8, p *[]uint8) *uint8sFlagVar {
	*p = val
	return (*uint8sFlagVar)(p)
}

type uint16FlagVar uint16

func newUint16Value(val uint16, p *uint16) *uint16FlagVar {
	*p = val
	return (*uint16FlagVar)(p)
}

func (f *uint16FlagVar) String() string {
	return fmt.Sprint(uint16(*f))
}

func (f *uint16FlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = uint16FlagVar(i)
	return nil
}

type uint16sFlagVar []uint16

func (f *uint16sFlagVar) String() string {
	return fmt.Sprint([]uint16(*f))
}

func (f *uint16sFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, uint16(i))
	return nil
}

func newUint16sValue(val []uint16, p *[]uint16) *uint16sFlagVar {
	*p = val
	return (*uint16sFlagVar)(p)
}

type uint32FlagVar uint32

func newUint32Value(val uint32, p *uint32) *uint32FlagVar {
	*p = val
	return (*uint32FlagVar)(p)
}

func (f *uint32FlagVar) String() string {
	return fmt.Sprint(uint32(*f))
}

func (f *uint32FlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = uint32FlagVar(i)
	return nil
}

type uint32sFlagVar []uint32

func (f *uint32sFlagVar) String() string {
	return fmt.Sprint([]uint32(*f))
}

func (f *uint32sFlagVar) Set(value string) error {
	i, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*f = append(*f, uint32(i))
	return nil
}

func newUint32sValue(val []uint32, p *[]uint32) *uint32sFlagVar {
	*p = val
	return (*uint32sFlagVar)(p)
}

type uint64sFlagVar []uint64

func (f *uint64sFlagVar) String() string { return fmt.Sprint([]uint64(*f)) }

func (f *uint64sFlagVar) Set(value string) error {
	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	*f = append(*f, i)
	return nil
}

func newUint64sValue(val []uint64, p *[]uint64) *uint64sFlagVar {
	*p = val
	return (*uint64sFlagVar)(p)
}

type float32sFlagVar []float32

func (f *float32sFlagVar) String() string { return fmt.Sprint([]float32(*f)) }

func (f *float32sFlagVar) Set(value string) error {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	*f = append(*f, float32(i))
	return nil
}

func newFloat32sValue(val []float32, p *[]float32) *float32sFlagVar {
	*p = val
	return (*float32sFlagVar)(p)
}

type float32FlagVar float32

func (f *float32FlagVar) String() string { return fmt.Sprint(float32(*f)) }

func (f *float32FlagVar) Set(value string) error {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	*f = float32FlagVar(i)
	return nil
}

func newFloat32Value(val float32, p *float32) *float32FlagVar {
	*p = val
	return (*float32FlagVar)(p)
}

type float64sFlagVar []float64

func (f *float64sFlagVar) String() string { return fmt.Sprint([]float64(*f)) }

func (f *float64sFlagVar) Set(value string) error {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	*f = append(*f, i)
	return nil
}

func newFloat64sValue(val []float64, p *[]float64) *float64sFlagVar {
	*p = val
	return (*float64sFlagVar)(p)
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
		case reflect.Slice:
			c.addSliceField(field, posArg, c.sliceSeparator[key])
		default:
			return errors.New("Unsupported type for positional argument: " + fieldType.Name())
		}
	}
	return nil
}
