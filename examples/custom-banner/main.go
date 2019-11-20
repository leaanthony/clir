package main

import (
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

	cli.SetBannerFunction(customBanner)

	// Define action
	cli.Action(func() error {
		println("Hello World!")
		return nil
	})

	// Run!
	cli.Run()

}
