---
title: "SubCommands"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 60
chapter: false
---

SubCommands allow you to add extra functionality to your application. This is done by adding child commands to your main cli app. SubCommands may be nested so there's no limit to how you can structure your app.

### Creating Subcommands for your app

```go
package main

import (
  "fmt"
  "log"

  "github.com/leaanthony/clir"
)

func main() {

  // Create the application
  cli := clir.NewCli("SubCommand", "A simple example", "v0.0.1")

  // Create subcommand
  initCmd := cli.NewSubCommand("init", "Initialise a new Project")

  // Define action for the command
  initCmd.Action(func() error {
    fmt.Println("Initialising Project!")
    return nil
  })

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

Running this command will display the default help:

```shell
> subcommand
SubCommand v0.0.1 - A simple example

Available commands:

   init   Initialise a new Project

Flags:

  -help
        Get help on the 'subcommand' command.
```

If we run `subcommand init`, we will get our output:

```shell
> subcommand init
Initialising Project!
```

Running `subcommand init -help` will give us the help message for that subcommand. 

```shell
> subcommand init --help
SubCommand v0.0.1 - A simple example

SubCommand init - Initialise a new Project
Flags:

  -help
        Get help on the 'subcommand init' command.
```

If you wish to add more information to help messages, use the `Command.LongDescription()` method:

```go
package main

import (
  "fmt"
  "log"

  "github.com/leaanthony/clir"
)

func main() {

  // Create the application
  cli := clir.NewCli("SubCommand", "A simple example", "v0.0.1")

  // Create subcommand
  initCmd := cli.NewSubCommand("init", "Initialise a new Project")

  // More help
  initCmd.LongDescription("The init command initialises a new project in the current working directory.")
  
  // Define action for the command
  initCmd.Action(func() error {
    fmt.Println("Initialising Project!")
    return nil
  })

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```
Then when we pass `-help`, we get more information:

```shell
> subcommand init --help
SubCommand v0.0.1 - A simple example

SubCommand init - Initialise a new Project
The init command initialises a new project in the current working directory.

Flags:

  -help
        Get help on the 'subcommand init' command.
```

You can add as many subcommands as you like. You can even nest them!

### Nested SubCommands

As Commands have (basically) the same API as the Cli object, we can do everything we did in the previous pages: Add Flags, Actions and SubCommands to any Command:

```go
package main

import (
  "fmt"
  "log"

  "github.com/leaanthony/clir"
)

func main() {

  // Create the application
  cli := clir.NewCli("SubCommand", "A simple example", "v0.0.1")

  // Create subcommand
  initCmd := cli.NewSubCommand("init", "Initialise a component")

  // Create a new "project" command below the "init" command
  projectCmd := initCmd.NewSubCommand("project", "Creates a new project")
  projectCmd.Action(func() error {
    fmt.Println("Initialising Project!")
    return nil
  })

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

Running `subcommand init` shows the following help:

```shell
> subcommand init
SubCommand v0.0.1 - A simple example

SubCommand init - Initialise a component
Available commands:

   project   Creates a new project

Flags:

  -help
        Get help on the 'subcommand init' command.

```

Whilst running `subcommand init project` will call the project command action:

```shell
> subcommand init project
Initialising Project!
```

### Adding SubCommands

It's possible to define SubCommands in isolation, then add them to a command later. To do this, we use the `Command.AddCommand` method:

```go
package main

import (
  "fmt"
  "log"

  "github.com/leaanthony/clir"
)

func newProjectCommand() *clir.Command {
  // Create a new Command
  result := clir.NewCommand("project", "Creates a new project")

  // Define the Action
  result.Action(func() error {
    fmt.Println("Initialising Project!")
    return nil
  })
  
  return result
}

func main() {

  // Create the application
  cli := clir.NewCli("SubCommand", "A simple example", "v0.0.1")

  // Create subcommand
  initCmd := cli.NewSubCommand("init", "Initialise a component")

  // Create a new "project" command and add it to the "init" command
  initCmd.AddCommand(newProjectCommand())

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

This really helps with better code organisation.

### Inheriting Flags

The `NewSubCommandInheritFlags` method will create a subcommand in the usual way but will inherit all previously defined flags in the parent.

### Hidden SubCommands

It's possible to hide subcommands by calling the `Hidden` method. This will omit the subcommand from any help text.

```go
package main

import (
  "fmt"
  "log"

	"github.com/leaanthony/clir"
)

func main() {

	// Create the application
	cli := clir.NewCli("SubCommand", "A simple example", "v0.0.1")

	// Create subcommand
	initCmd := cli.NewSubCommand("init", "Initialise a component")
	initCmd.Action(func() error {
		fmt.Println("Initialising")
		return nil
	})

	// Create a hidden developer command
	devtoolsCommand := cli.NewSubCommand("dev", "Developer tools")
	devtoolsCommand.Action(func() error {
		fmt.Println("I'm a secret command")
		return nil
	})
	devtoolsCommand.Hidden()

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

The main help hides the command, but it is still possible to run it:

```
> subcommand
SubCommand v0.0.1 - A simple example

Available commands:

   init   Initialise a component

Flags:

  -help
        Get help on the 'subcommand' command.


> subcommand dev
I'm a secret command
```

### Default Command

If you would like the default action of your application to be a particular subcommand, you can set it using the `DefaultCommand` method:

```go
package main

import (
  "fmt"
  "log"

	"github.com/leaanthony/clir"
)

func main() {

	// Create the application
	cli := clir.NewCli("SubCommand", "A simple example", "v0.0.1")

	// Create subcommand
	initCmd := cli.NewSubCommand("init", "Initialise a component")
	initCmd.Action(func() error {
		fmt.Println("Initialising")
		return nil
	})

	// Make init the default command
	cli.DefaultCommand(initCmd)

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

Now when we run the application with no parameters, it will run the init subcommand:

```shell
> subcommand
Initialising
```

The help text also indicates that there is a default option:

```
> subcommand -help
SubCommand v0.0.1 - A simple example

Available commands:

   init   Initialise a component [default]

Flags:

  -help
        Get help on the 'subcommand' command.
```

### API

**NewCommand(name string, description string) *Command**

Creates a new Command in isolation. It may be attached to the application through `AddCommand(cmd *Command)`. The returned Command may be further configured.

**AddCommand(cmd *Command)**

Attaches the given Command as a SubCommand. This API is valid for both Cli and Command.

**Hidden()**

Hides the command from help message.
