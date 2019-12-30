// Package clir provides a simple API for creating command line apps
package clir

import (
	"fmt"
)

// defaultBannerFunction prints a banner for the application.
// If version is a blank string, it is ignored.
func defaultBannerFunction(c *Cli) string {
	version := ""
	if len(c.Version()) > 0 {
		version = " " + c.Version()
	}
	return fmt.Sprintf("%s%s - %s", c.Name(), version, c.ShortDescription())
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
