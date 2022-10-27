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

	t.Run("Run SetFlagFunction()", func(t *testing.T) {
		c.SetErrorFunction(func(cmdPath string, err error) error { return err })
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

type testStruct struct {
	Mode  string `name:"mode" description:"The mode of build"`
	Count int
}

func TestCli_CLIAddFlags(t *testing.T) {
	c := NewCli("test", "description", "0")

	ts := &testStruct{}
	c.AddFlags(ts)

	modeFlag := c.rootCommand.flags.Lookup("mode")
	if modeFlag == nil {
		t.Errorf("expected flag mode to be added")
	}
	if modeFlag.Name != "mode" {
		t.Errorf("expected flag name to be added")
	}
	if modeFlag.Usage != "The mode of build" {
		t.Errorf("expected flag description to be added")
	}

	c.Action(func() error {
		if ts.Mode != "123" {
			t.Errorf("expected flag value to be set")
		}
		return nil
	})
	e := c.Run("-mode", "123")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFlags(t *testing.T) {
	c := NewCli("test", "description", "0")
	sub := c.NewSubCommand("sub", "sub description")

	ts := &testStruct{}
	sub.AddFlags(ts)

	sub.Action(func() error {
		if ts.Mode != "123" {
			t.Errorf("expected flag value to be set")
		}
		return nil
	})
	e := c.Run("sub", "-mode", "123")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}
