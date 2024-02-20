package clir

import (
	"fmt"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

// To understand where testscript comes from and what it does, see:
// - https://github.com/rogpeppe/go-internal/tree/master?tab=readme-ov-file#testscript
// - https://bitfieldconsulting.com/golang/test-scripts
func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"clir": Main,
	}))
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata",
	})
}

type Flags struct {
	Bar int `flag:"bar" description:"How many bars"`
}

func Main() int {
	cli := NewCli("clir", "a clir executable for testscript", "")
	flags := &Flags{}
	cli.NewSubCommand("foo", "Foo the bar").
		AddFlags(flags).
		Action(func() error { return foo(flags) })

	if err := cli.Run(); err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	return 0
}

func foo(flags *Flags) error {
	fmt.Println("this is Foo. Flags:", flags)
	return nil
}
