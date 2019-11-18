package main

import (
	"github.com/leaanthony/clir"
)

func main() {

	app := clir.NewCli("Basic", "A Basic Example", "v0.0.1")
	err := app.Run()
	if err != nil {
		app.Abort(err)
	}
}
