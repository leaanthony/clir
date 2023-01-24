---
title: "Custom Banner"
date: 2019-11-21T16:15:09+08:00
draft: false
weight: 90
chapter: false
---

It is possible to provide your own banner by setting the Banner function that the application will call: 

```go
package main

import (
	"log"
	
	"github.com/leaanthony/clir"
)

func customBanner(cli *clir.Cli) string {

	return `
 ______   __  __   ______   _________  ______   ___ __ __     
/_____/\ /_/\/_/\ /_____/\ /________/\/_____/\ /__//_//_/\    
\:::__\/ \:\ \:\ \\::::_\/_\__.::.__\/\:::_ \ \\::\| \| \ \   
 \:\ \  __\:\ \:\ \\:\/___/\  \::\ \   \:\ \ \ \\:.      \ \  
  \:\ \/_/\\:\ \:\ \\_::._\:\  \::\ \   \:\ \ \ \\:.\-/\  \ \ 
   \:\_\ \ \\:\_\:\ \ /____\:\  \::\ \   \:\_\ \ \\. \  \  \ \
    \_____\/ \_____\/ \_____\/   \__\/    \_____\/ \__\/ \__\/
     ` + cli.Version() + " - " + cli.ShortDescription()
}

func main() {

	// Create new cli
	cli := clir.NewCli("Banner", "A custom banner example", "v0.0.1")

	// Set the custom banner
	cli.SetBannerFunction(customBanner)

  // Run the application
  err := cli.Run()
  if err != nil {
    // We had an error
    log.Fatal(err)
  }

}
```

The `setBannerFunction` method expects a function with the signature `func (*clir.Cli) string`. 
