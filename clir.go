// Package clir provides a simple API for creating command line apps
package clir

// NewCli - Creates a new Cli application object
func NewCli(name, description, version string) *Cli {
	result := &Cli{
		version: version,
	}
	result.rootCommand = NewCommand(name, description)
	result.rootCommand.setApp(result)
	result.rootCommand.setParentCommandPath("")
	return result
}
