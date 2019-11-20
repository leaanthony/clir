# clir

A Simple and Clear CLI library. Dependency free.

### Example

```
package main

import (
	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Basic", "A basic example", "v0.0.1")

	// Define action
	cli.Action(func() error {
		println("Hello World!")
		return nil
	})

	// Run!
	cli.Run()

}
```

### Features

  * Nested Subcommands
  * Auto-generated help
  * Custom banners
  * Uses the standard library `flag` package
  * 