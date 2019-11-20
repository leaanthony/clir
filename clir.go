// Package clir provides a simple API for creating command line apps
package clir

import (
	"fmt"
)

func defaultBannerFunction(c *Cli) string {
	return fmt.Sprintf("%s %s - %s", c.Name(), c.Version(), c.ShortDescription())
}

// NewCli - Creates a new Cli application object
func NewCli(name, description, version string) *Cli {
	result := &Cli{
		version:        version,
		bannerFunction: defaultBannerFunction,
	}
	result.rootCommand = NewCommand(name, description)
	result.rootCommand.setApp(result)
	result.rootCommand.setParentCommandPath("")
	return result
}
