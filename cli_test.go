package clir

import (
	"errors"
	"testing"
)

func TestCli(t *testing.T) {
	c := NewCli("test", "description", "0")

	t.Run("Run SetBannerFunction()", func(t *testing.T) {
		c.SetBannerFunction(func(*Cli) string { return "" })
	})

	// t.Run("Run Abort()", func(t *testing.T) {
	// 	cl := NewCli("test", "description", "0")
	// 	cl.Abort(errors.New("test error"))
	// })

	t.Run("Run AddCommand()", func(t *testing.T) {
		c.AddCommand(&Command{name: "test"})
	})

	t.Run("Run PrintBanner()", func(t *testing.T) {
		c.PrintBanner()
	})

	t.Run("Run Run()", func(t *testing.T) {
		c.Run("test")
		c.Run()

		c.preRunCommand = func(*Cli) error { return errors.New("testing coverage") }
		c.Run("test")
	})

	t.Run("Run DefaultCommand()", func(t *testing.T) {
		c.DefaultCommand(&Command{})
	})

	t.Run("Run NewSubCommand()", func(t *testing.T) {
		c.NewSubCommand("name", "description")
	})

	t.Run("Run PreRun()", func(t *testing.T) {
		c.PreRun(func(*Cli) error { return nil })
	})

	t.Run("Run BoolFlag()", func(t *testing.T) {
		var variable bool
		c.BoolFlag("bool", "description", &variable)
	})

	t.Run("Run StringFlag()", func(t *testing.T) {
		var variable string
		c.StringFlag("string", "description", &variable)
	})

	t.Run("Run IntFlag()", func(t *testing.T) {
		var variable int
		c.IntFlag("int", "description", &variable)
	})

	t.Run("Run Action()", func(t *testing.T) {
		c.Action(func() error { return nil })
	})

	t.Run("Run LongDescription()", func(t *testing.T) {
		c.LongDescription("long description")
	})
}
