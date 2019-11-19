# clir

A Simple and Clear CLI library.

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