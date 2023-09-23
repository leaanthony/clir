---
title: "Pre/Post run"
date: 2023-09-23T16:15:09+08:00
draft: false
weight: 100
chapter: false
---

The PreRun and PostRun methods allow you to specify custom functions that should be executed before and after running a command, respectively. These functions can be used for tasks such as setup, cleanup, or any other actions you want to perform before or after executing a command.

PreRun

The PreRun method is used to specify a function that should run before executing a command. The function you pass to PreRun takes a *Cli parameter, which represents the CLI application itself. You can use this function to perform any necessary setup or validation before running the command.

PostRun
The PostRun method is used to specify a function that should run after executing a command. Similar to PreRun, the function you pass to PostRun takes a *Cli parameter. You can use this function to perform any cleanup or post-processing tasks after the command execution.

Example of postRun
```go
func main() {
    cli := clir.NewCli("MyApp", "My CLI application")

    // Define a command
    cmd := cli.NewSubCommand("mycommand", "Description of mycommand")

    // Add flags, actions, and other configurations for the command

    // Define a PostRun function for the command
    cmd.PostRun(func(c *clir.Cli) error {
        // Perform cleanup or post-processing here
        fmt.Println("Running PostRun for 'mycommand'")
        return nil // Return an error if something goes wrong
    })

    // Run the CLI application
    if err := cli.Run(os.Args[1:]...); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

```