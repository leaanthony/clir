package clir

import (
	"errors"
	"fmt"
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
func TestCli_InheritFlags(t *testing.T) {
	c := NewCli("test", "description", "0")
	var name string
	c.StringFlag("name", "name of person", &name)
	sub := c.NewSubCommandInheritFlags("sub", "sub description")

	sub.Action(func() error {
		if name != "Janet" {
			t.Errorf("expected name to be Janet")
		}
		return nil
	})
	e := c.Run("sub", "-name", "Janet")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFlagsWrongType(t *testing.T) {
	c := NewCli("test", "description", "0")
	sub := c.NewSubCommand("sub", "sub description")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "AddFlags() requires a pointer to a struct" {
			t.Errorf(`Expected: "AddFlags() requires a pointer to a struct". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	ts := testStruct{}
	sub.AddFlags(ts)

	e := c.Run("sub", "-mode", "123")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFlagsWrongPointerType(t *testing.T) {
	c := NewCli("test", "description", "0")
	sub := c.NewSubCommand("sub", "sub description")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "AddFlags() requires a pointer to a struct" {
			t.Errorf(`Expected: "AddFlags() requires a pointer to a struct". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	var i *int
	sub.AddFlags(i)

	e := c.Run("sub", "-mode", "123")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandOtherSubArgs(t *testing.T) {
	c := NewCli("test", "description", "0")
	sub := c.NewSubCommand("sub", "sub description")

	ts := &testStruct{}
	sub.AddFlags(ts)

	sub.Action(func() error {
		if ts.Mode != "123" {
			t.Errorf("expected flag value to be set")
		}
		other := sub.OtherArgs()
		if len(other) != 2 {
			t.Errorf("expected 2 other args, got %v", other)
		}
		if other[0] != "other" {
			t.Errorf("expected other arg to be 'other', got %v", other[0])
		}
		if other[1] != "args" {
			t.Errorf("expected other arg to be 'args', got %v", other[1])
		}
		return nil
	})
	e := c.Run("sub", "-mode", "123", "other", "args")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}
func TestCli_CommandOtherCLIArgs(t *testing.T) {
	c := NewCli("test", "description", "0")

	c.Action(func() error {
		other := c.OtherArgs()
		if len(other) != 2 {
			t.Errorf("expected 2 other args, got %v", other)
		}
		if other[0] != "other" {
			t.Errorf("expected other arg to be 'other', got %v", other[0])
		}
		if other[1] != "args" {
			t.Errorf("expected other arg to be 'args', got %v", other[1])
		}
		return nil
	})
	e := c.Run("other", "args")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

type Person struct {
	Name        string  `description:"The name of the person"`
	Age         int     `description:"The age of the person"`
	SSN         uint    `description:"The SSN of the person"`
	Age64       int64   `description:"The age of the person"`
	SSN64       uint64  `description:"The SSN of the person"`
	BankBalance float64 `description:"The bank balance of the person"`
	Married     bool    `description:"Whether the person is married"`
}

func TestCli_CommandAddFunction(t *testing.T) {
	c := NewCli("test", "description", "0")

	createPerson := func(person *Person) error {
		if person.Name != "Bob" {
			t.Errorf("expected name flag to be 'Bob', got %v", person.Name)
		}
		if person.Age != 30 {
			t.Errorf("expected age flag to be 30, got %v", person.Age)
		}
		if person.SSN != 123456789 {
			t.Errorf("expected ssn flag to be 123456789, got %v", person.SSN)
		}
		if person.Age64 != 30 {
			t.Errorf("expected age64 flag to be 30, got %v", person.Age64)
		}
		if person.SSN64 != 123456789 {
			t.Errorf("expected ssn64 flag to be 123456789, got %v", person.SSN64)
		}
		if person.BankBalance != 123.45 {
			t.Errorf("expected bankbalance flag to be 123.45, got %v", person.BankBalance)
		}
		if person.Married != true {
			t.Errorf("expected married flag to be true, got %v", person.Married)
		}
		return nil
	}

	c.NewSubCommandFunction("create", "create a person", createPerson)

	e := c.Run("create", "-name", "Bob", "-age", "30", "-ssn", "123456789", "-age64", "30", "-ssn64", "123456789", "-bankbalance", "123.45", "-married")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

type Person2 struct {
	Name string `description:"The name of the person"`
}

func (p *Person2) Default() *Person2 {
	return &Person2{Name: "Janet"}
}

func TestCli_CommandAddFunctionDefaults(t *testing.T) {
	c := NewCli("test", "description", "0")

	createPerson := func(person *Person2) error {
		if person.Name != "Janet" {
			t.Errorf("expected name flag to be 'Janet', got %v", person.Name)
		}
		return nil
	}

	c.NewSubCommandFunction("create", "create a person", createPerson)

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFunctionNoFunctionError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", 0)

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}
func TestCli_CommandAddFunctionNoPointerError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person Person2) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFunctionMultipleInputError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person2, count int) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFunctionNoInputError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func() error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}

}

func TestCli_CommandAddFunctionNoReturnError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person2) {
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

func TestCli_CommandAddFunctionWrongReturnError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person2) int {
		return 0
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

type Person3 struct {
	Name string `description:"The name of the person"`
}

func (p *Person3) Default(_ int) *Person3 {
	return &Person3{Name: "Janet"}
}

func TestCli_CommandAddFunctionDefaultWrongInputError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "'Default' method on struct 'Person3' must have the signature 'Default() *Person3'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person3) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

type Person4 struct {
	Name string `description:"The name of the person"`
}

func (p *Person4) Default() int {
	return 0
}

func TestCli_CommandAddFunctionDefaultWrongOutputError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "'Default' method on struct 'Person4' must have the signature 'Default() *Person4'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person4) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

type Person5 struct {
	Name string `description:"The name of the person"`
}

func (p *Person5) Default(_ int, _ int) int {
	return 0
}

func TestCli_CommandAddFunctionDefaultWrongNumberOfInputsError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "'Default' method on struct 'Person5' must have the signature 'Default() *Person5'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person5) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

type Person6 struct {
	Name  string `description:"The name of the person"`
	Admin bool   `description:"Is the person an admin"`
}

func (p *Person6) Default() (int, int) {
	return 0, 0
}

func TestCli_CommandAddFunctionDefaultWrongNumberOfOutputsError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "'Default' method on struct 'Person6' must have the signature 'Default() *Person6'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *Person6) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

func TestCli_CommandAddFunctionDefaultWrongTypeOfInputsError(t *testing.T) {
	c := NewCli("test", "description", "0")

	recoverTest := func() {
		var r interface{}
		if r = recover(); r == nil {
			t.Errorf("expected panic")
		}
		if r.(string) != "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'" {
			t.Errorf(`expected error message to be: "NewSubFunction 'create' requires a function with the signature 'func(*struct) error'". Got: "` + r.(string) + `"`)
		}
	}

	defer recoverTest()

	c.NewSubCommandFunction("create", "create a person", func(person *int) error {
		return nil
	})

	e := c.Run("create")
	if e != nil {
		t.Errorf("expected no error, got %v", e)
	}
}

type Person7 struct {
	Name  string `description:"The name of the person"`
	Admin bool   `description:"Is the person an admin"`
}

func TestCli_CommandAddFunctionReturnsAnError(t *testing.T) {
	c := NewCli("test", "description", "0")

	c.NewSubCommandFunction("create", "create a person", func(person *Person7) error {
		return fmt.Errorf("error")
	})

	e := c.Run("create")
	if e == nil {
		t.Errorf("expected error, got nil")
	}
}
