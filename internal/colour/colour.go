package colour

import "github.com/fatih/color"

// YellowString returns a function that returns a Yellow String
var YellowString = color.New(color.FgYellow).SprintFunc()

// RedString returns a function that returns a Red String
var RedString = color.New(color.FgRed).SprintFunc()

// WhiteString returns a function that returns a White String
var WhiteString = color.New(color.FgWhite).SprintFunc()

// GreenString returns a function that returns a Green String
var GreenString = color.New(color.FgGreen).SprintFunc()

// Yellow returns a function that prints the given string in Yellow
var Yellow = color.New(color.FgYellow).PrintFunc()

// Red returns a function that prints the given string in Red
var Red = color.New(color.FgRed).PrintFunc()

// White returns a function that prints the given string in White
var White = color.New(color.FgWhite).PrintFunc()

// Green returns a function that prints the given string in Green
var Green = color.New(color.FgGreen).PrintFunc()

// YellowStringLn returns a function that prints the given string in YellowStringLn
// It appends a line ending
var YellowStringLn = color.New(color.FgYellow).PrintlnFunc()

// RedStringLn returns a function that prints the given string in RedStringLn
// It appends a line ending
var RedStringLn = color.New(color.FgRed).PrintlnFunc()

// WhiteStringLn returns a function that prints the given string in WhiteStringLn
// It appends a line ending
var WhiteStringLn = color.New(color.FgWhite).PrintlnFunc()

// GreenStringLn returns a function that prints the given string in GreenStringLn
// It appends a line ending
var GreenStringLn = color.New(color.FgGreen).PrintlnFunc()
