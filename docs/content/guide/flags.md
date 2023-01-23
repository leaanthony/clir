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

### Defining flags using a struct

It's also possible to add flags using a struct. This is useful as you can easily
add flags to your application without having to define them in multiple places:

```go
package main

import (
	"github.com/leaanthony/clir"
)

type Flags struct {
	Name string `name:"name" description:"The name of the person"`
	Age  int    `name:"age" description:"The age of the person"`
}

func main() {

	// Create new cli
	cli := clir.NewCli("flagstruct", "An example of subcommands with flag inherence", "v0.0.1")

	// Create an init subcommand with flag inheritance
	init := cli.NewSubCommand("create", "Create a person")
	person := &Flags{
		Age: 20,
	}
	init.AddFlags(person)
	init.Action(func() error {
		println("Name:", person.Name, "Age:", person.Age)
		return nil
	})

	// Run!
	if err := cli.Run(); err != nil {
		panic(err)
	}

}
```

### Defining flags using a struct with default values

It's also possible to add flags using a struct with default values. This is useful as you can easily add flags to your application without having to define them in multiple places:

```go
package main

import (
    "github.com/leaanthony/clir"
)

type Flags struct {
    Name string `name:"name" description:"The name of the person" default:"John"`
    Age  int    `name:"age" description:"The age of the person" default:"20"`
}

func main() {

    // Create new cli
    cli := clir.NewCli("default", "An example of subcommands with flag inheritance", "v0.0.1")

    // Create an init subcommand with flag inheritance
    init := cli.NewSubCommand("create", "Create a person")
    person := &Flags{}
    init.AddFlags(person)
    init.Action(func() error {
        println("Name:", person.Name, "Age:", person.Age)
        return nil
    })

    // Run!
    if err := cli.Run(); err != nil {
        panic(err)
    }

}
```

### Defining positional arguments

It is possible to define positional arguments. These are arguments that are not
flags and are defined in the order they are passed to the application. For
example:


```go
package main

import (
    "github.com/leaanthony/clir"
)

type Flags struct {
    Name string `pos:"1" description:"The name of the person" default:"John"`
    Age  int    `pos:"2" description:"The age of the person" default:"20"`
}

func main() {

    // Create new cli
    cli := clir.NewCli("default", "An example of subcommands with positional args", "v0.0.1")

    // Create an init subcommand with flag inheritance
    init := cli.NewSubCommand("create", "Create a person")
    person := &Flags{}
    init.AddFlags(person)
    init.Action(func() error {
        println("Name:", person.Name, "Age:", person.Age)
        return nil
    })

    // Run!
    if err := cli.Run("create", "bob", "30"); err != nil {
        panic(err)
    }

}
```

### API

#### Cli.StringFlag(name string, description string, variable *string)

The [StringFlag](https://godoc.org/github.com/leaanthony/clir#StringFlag) method defines a string flag for your Clîr application. 

For the example above, you would pass in a name as follows:

```shell
> flags -name John
```


#### Cli.IntFlag(name string, description string, variable *int)

The [IntFlag](https://godoc.org/github.com/leaanthony/clir#IntFlag) method defines an integer flag for your Clîr application. 

For the example above, you would pass in a value for age as follows:

```shell
> flags -age 32
```


#### Cli.BoolFlag(name string, description string, variable *bool)

The [BoolFlag](https://godoc.org/github.com/leaanthony/clir#BoolFlag) method defines a boolean flag for your Clîr application. 

For the example above, you would you were awesome by simply passing in the flag:

```shell
> flags -awesome
```

#### Cli.AddFlags(config interface{})

The [AddFlags](https://godoc.org/github.com/leaanthony/clir#AddFlags) method defines flags for your Clîr application 
using a struct. It uses the `name` and `description` tags to define the flag name and description.
If no `name` tag is given, the field name is used. If no `description` tag is given, the description will be blank.
A struct pointer must be passed in otherwise the method will panic.