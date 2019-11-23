<p align="center" style="text-align: center">
   <img src="clir_logo.png" width="40%"><br/>
</p>
<p align="center">
   A Simple and Clear CLI library. Dependency free.<br/><br/>
   <a href="https://github.com/leaanthony/clir/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg"></a>
   <a href="https://goreportcard.com/report/github.com/leaanthony/clir"><img src="https://goreportcard.com/badge/github.com/leaanthony/clir"/></a>
   <a href="http://godoc.org/github.com/leaanthony/clir"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"/></a>
   <a href="https://www.codefactor.io/repository/github/leaanthony/clir"><img src="https://www.codefactor.io/repository/github/leaanthony/clir/badge" alt="CodeFactor" /></a>
   <a href="https://github.com/leaanthony/clir/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat" alt="CodeFactor" /></a>
   <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Fleaanthony%2Fclir?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Fleaanthony%2Fclir.svg?type=shield"/></a>
   <a href="https://houndci.com"><img src="https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg"/></a>

</p>


### Features

  * Nested Subcommands
  * Uses the standard library `flag` package
  * Auto-generated help
  * Custom banners
  * Hidden Subcommands
  * Default Subcommand
  * Dependency free

### Example

```go
package main

import (
	"fmt"

	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Flags", "A simple example", "v0.0.1")

	// Name
	var name string
	cli.StringFlag("name", "Your name", &name)
	
	// Define action for the command
	cli.Action(func() error {

		if name == "" {
			name = "Anonymous"
		}
		fmt.Printf("Hello %s!\n", name)

		return nil
	})

	cli.Run()

}
```

#### Generated Help

```shell
$ flags --help
Flags v0.0.1 - A simple example

Flags:

  -help
        Get help on the 'flags' command.
  -name string
        Your name
```

#### Documentation

The main documentation may be found [here](https://clir.leaanthony.com).
