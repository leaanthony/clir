---
title: "Flags"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 40
chapter: false
---

Flags may be added to your application as follows:

```go
package main

import "github.com/leaanthony/clir"
import "log"

func main() {
  
  // Create the application
  cli := clir.NewCli("Flags", "A simple example", "v0.0.1")

  // Add a string flag
  var name string
  cli.StringFlag("name", "Your name", &name)

  // Add an int flag
  var age int
  cli.IntFlag("age", "Your age", &age)

  // Add a bool flag
  var awesome bool
  cli.BoolFlag("awesome", "Are you awesome?", &awesome)

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }
}
```

Running this command prints the default help text as expected:

```shell
Flags v0.0.1 - A simple example

Flags:

  -age int
        Your age
  -awesome
        Are you awesome?
  -help
        Get help on the 'flags' command.
  -name string
        Your name
```

## API

**Cli.StringFlag(name string, description string, variable *string)**

The [StringFlag](https://godoc.org/github.com/leaanthony/clir#StringFlag) method defines a string flag for your Clîr application. 

For the example above, you would pass in a name as follows:

```shell
> flags -name John
```


**Cli.IntFlag(name string, description string, variable *int)**

The [IntFlag](https://godoc.org/github.com/leaanthony/clir#IntFlag) method defines an integer flag for your Clîr application. 

For the example above, you would pass in a value for age as follows:

```shell
> flags -age 32
```


**Cli.BoolFlag(name string, description string, variable *bool)**

The [BoolFlag](https://godoc.org/github.com/leaanthony/clir#BoolFlag) method defines a boolean flag for your Clîr application. 

For the example above, you would you were awesome by simply passing in the flag:

```shell
> flags -awesome
```
