package clir

// Action represents a function that gets calls when the command is called by
// the user
type Action func() error
