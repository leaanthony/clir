---
title: "Getting Started"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 10
---

## Installation

`go get github.com/leaanthony/clir`

## Basic Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/leaanthony/clir"
)

func main() {

	// Create new cli
	cli := clir.NewCli("Flags", "A simple example", "v0.0.1")

	// Name
	name := "Anonymous"
	cli.StringFlag("name", "Your name", &name)
	
	// Define action for the command
	cli.Action(func() error {
		fmt.Printf("Hello %s!\n", name)
		return nil
	})

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```