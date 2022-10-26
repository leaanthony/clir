---
title: "CLI"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 30
chapter: false
---

The main entry point for a Clîr application is the Cli object. This is created using the `NewCli` command which takes an application name, description and optional version:

```go
package main

import "github.com/leaanthony/clir"
import "log"

func main() {
  
  // Create the application
  cli := clir.NewCli("Basic", "A basic example", "v0.0.1")

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }
}
```

When you run this app, you will get the default help text:

```shell
> basic
Basic v0.0.1 - A basic example

Flags:

  -help
        Get help on the 'basic' command.

```

### API

**NewCli(name string, description string, version string) *Cli**

The [NewCli](https://godoc.org/github.com/leaanthony/clir#NewCli) function creates a new Clîr application. 

**Cli.Run(args ...string) error**

The [Run](https://godoc.org/github.com/leaanthony/clir#Cli.Run) method starts the application. By default it will use `os.Args`, though you are free to pass in arguments for testing purposes. Run returns an error, which may be handled appropriately.
