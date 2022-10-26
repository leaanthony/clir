---
title: "Other Args"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 45
chapter: false
---

Other arguments passed to your application are accessible within actions using `cli.OtherArgs()`:

```go
package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Other Args", "Access other arguments", "v0.0.1")

	// Set long description
	cli.LongDescription("This app shows how to access non-flag arguments")

	// Name
	var name string
	cli.StringFlag("name", "Your name", &name)

	// Define action
	cli.Action(func() error {
		println("Your name is", name)
		fmt.Printf("The remaining arguments were: %+v\n", cli.OtherArgs())
		return nil
	})

	// Run!
	cli.Run()

}
```

Running this command prints the following:

```shell
$ ./otherargs -name test other args
Your name is test
The remaining arguments were: [other args]
```

**Cli.OtherArgs() []string**

The [OtherArgs](https://godoc.org/github.com/leaanthony/clir#Cli.OtherArgs) method returns all arguments to the application that are not handled by the defined flags. *NOTE*: This will only return correct values if accessed in an Action.
