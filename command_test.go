package clir

import (
	"testing"
)

func TestCommand(t *testing.T) {
	c := &Command{}

	t.Run("Run NewCommand()", func(t *testing.T) {
		c = NewCommand("test", "test description")
	})

	t.Run("Run setParentCommandPath()", func(t *testing.T) {
		c.setParentCommandPath("path")
	})

	t.Run("Run setApp()", func(t *testing.T) {
		c.setApp(&Cli{})
	})

	t.Run("Run parseFlags()", func(t *testing.T) {
		err := c.parseFlags([]string{"test", "flags"})
		t.Log(err)
	})

	t.Run("Run run()", func(t *testing.T) {
		cl := NewCli("test", "description", "0")
		cl.rootCommand.run([]string{"test"})

		cl.rootCommand.subCommandsMap["test"] = &Command{
			name:             "subcom",
			shortdescription: "short description",
			hidden:           false,
			app:              cl,
		}
		cl.rootCommand.run([]string{"test"})

		cl.rootCommand.run([]string{"---"})

		cl.rootCommand.run([]string{"-help"})

		// cl.rootCommand.actionCallback = func() error {
		// 	println("Hello World!")
		// 	return nil
		// }
		// cl.rootCommand.run([]string{"test"})

		cl.rootCommand.app.defaultCommand = &Command{
			name:             "subcom",
			shortdescription: "short description",
			hidden:           false,
			app:              cl,
		}
		cl.rootCommand.run([]string{"test"})
	})

	t.Run("Run Action()", func(t *testing.T) {
		c.Action(func() error { return nil })

	})

	t.Run("Run PrintHelp()", func(t *testing.T) {
		cl := NewCli("test", "description", "0")

		// co.shortdescription = ""
		cl.PrintHelp()

		cl.rootCommand.shortdescription = "test"
		cl.PrintHelp()

		cl.rootCommand.commandPath = "notTest"
		cl.PrintHelp()

		cl.rootCommand.longdescription = ""
		cl.PrintHelp()

		cl.rootCommand.longdescription = "test"
		cl.PrintHelp()

		mockCommand := &Command{
			name:              "test",
			shortdescription:  "short description",
			hidden:            true,
			longestSubcommand: len("test"),
		}
		cl.rootCommand.subCommands = append(cl.rootCommand.subCommands, mockCommand)
		cl.PrintHelp()

		mockCommand = &Command{
			name:             "subcom",
			shortdescription: "short description",
			hidden:           false,
			app:              cl,
		}
		cl.rootCommand.longestSubcommand = 10
		cl.rootCommand.subCommands = append(cl.rootCommand.subCommands, mockCommand)
		cl.PrintHelp()

		mockCommand = &Command{
			name:             "subcom",
			shortdescription: "short description",
			hidden:           false,
			app:              cl,
		}
		cl.rootCommand.longestSubcommand = 10
		cl.defaultCommand = mockCommand
		cl.rootCommand.subCommands = append(cl.rootCommand.subCommands, mockCommand)
		cl.PrintHelp()

		cl.rootCommand.flagCount = 3
		cl.PrintHelp()
	})

	t.Run("Run isDefaultCommand()", func(t *testing.T) {
		c.isDefaultCommand()

	})

	t.Run("Run isHidden()", func(t *testing.T) {
		c.isHidden()

	})

	t.Run("Run Hidden()", func(t *testing.T) {
		c.Hidden()
	})

	t.Run("Run NewSubCommand()", func(t *testing.T) {
		c.NewSubCommand("name", "description")

	})

	t.Run("Run AddCommand()", func(t *testing.T) {
		c.AddCommand(c)
	})

	t.Run("Run StringFlag()", func(t *testing.T) {
		var variable = "variable"
		c.StringFlag("name", "description", &variable)

	})

	t.Run("Run IntFlag()", func(t *testing.T) {
		var variable int
		c.IntFlag("test", "description", &variable)

	})

	t.Run("Run LongDescription()", func(t *testing.T) {
		c.LongDescription("name")

	})
}
