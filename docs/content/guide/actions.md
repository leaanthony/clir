---
title: "Actions"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 50
chapter: false
---

Actions define what code will be executed when your application is run. A basic example of this is:

```go
package main

import "github.com/leaanthony/clir"
import "log"

func main() {
  
  // Create the application
  cli := clir.NewCli("Actions", "A simple example", "v0.0.1")

  // Define main action
  cli.Action(func() error {
      println("Hello World!")
      return nil
  })

  // Run application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }
}
```

Running this command will simply run our action:

```shell
> actions
Hello World!
```

To interact with the flags passed in, we simply use scoping:

```go
package main

import (
  "fmt"
  "log"
  
  "github.com/leaanthony/clir"
)

func main() {

  // Create the application
  cli := clir.NewCli("Actions", "A simple example", "v0.0.1")

  // Set our default name to "Anonymous"
  name := "Anonymous"
  cli.StringFlag("name", "Your name", &name)
  
  // Define action for the command
  cli.Action(func() error {
    fmt.Printf("Hello %s!\n", name)
    return nil
  })

  // Run application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

When we run this with no flags we get:
```shell
> actions
Hello Anonymous!
```

Passing in a name produces the expected output:
```shell
> actions -name Debbie
Hello Debbie!
```

## API

**Cli.Action(fn func() error)**

Action binds the given function to the application. It is called when the application is executed. Any errors returned by actions are passed back to via the main `Cli.Run` method.

Example:

```go
package main

import (
  "fmt"
  "log"

  "github.com/leaanthony/clir"
)

func main() {

  // Create the application
  cli := clir.NewCli("Actions", "A simple example", "v0.0.1")

  // Define action for the command
  cli.Action(func() error {
    return fmt.Errorf("I am an error")
  })

  // We will receive our error here
  err := cli.Run()
  if err != nil {
    log.Fatal(err)
  }
}
```

Running this will produce:

```shell
> actions
2019/11/23 08:03:56 I am an error
```