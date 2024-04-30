package main

import (
	"fmt"
	"github.com/leaanthony/clir"
)

type Flags struct {
	String         string   `name:"string" description:"The string"`
	Strings        []string `name:"strings" description:"The strings"`
	StringsDefault []string `name:"strings_default" description:"The strings default" default:"one,two,three"`

	Int         int   `name:"int" description:"The int"`
	Ints        []int `name:"ints" description:"The ints"`
	IntsDefault []int `name:"ints_default" description:"The ints default" default:"3,4,5"`

	Int64         int64   `name:"int64" description:"The int64"`
	Int64s        []int64 `name:"int64s" description:"The int64s"`
	Int64sDefault []int64 `name:"int64s_default" description:"The int64s default" default:"3,4,5"`

	Uint         uint   `name:"uint" description:"The uint"`
	Uints        []uint `name:"uints" description:"The uints"`
	UintsDefault []uint `name:"uints_default" description:"The uints default" default:"3,4,5"`

	Uint64         uint64   `name:"uint64" description:"The uint64"`
	Uint64s        []uint64 `name:"uint64s" description:"The uint64s"`
	Uint64sDefault []uint64 `name:"uint64s_default" description:"The uint64s default" default:"3,4,5"`

	Float64         float64   `name:"float64" description:"The float64"`
	Float64s        []float64 `name:"float64s" description:"The float64s"`
	Float64sDefault []float64 `name:"float64s_default" description:"The float64s default" default:"3,4,5"`

	Bool         bool   `name:"bool" description:"The bool"`
	Bools        []bool `name:"bools" description:"The bools"`
	BoolsDefault []bool `name:"bools_default" description:"The bools default" default:"false,true,false,true"`
}

func main() {

	// Create new cli
	cli := clir.NewCli("flagstruct", "An example of subcommands with flag inherence", "v0.0.1")

	// Create an init subcommand with flag inheritance
	init := cli.NewSubCommand("flag", "print default")

	cli.DefaultCommand(init)

	flags := &Flags{
		String:   "zkep",
		Strings:  []string{"one", "two", "three"},
		Int:      1,
		Ints:     []int{1, 2, 3},
		Int64:    1,
		Int64s:   []int64{1, 2, 3},
		Uint:     uint(1),
		Uints:    []uint{1, 2, 3},
		Uint64:   uint64(1),
		Uint64s:  []uint64{1, 2, 3},
		Float64:  3.14,
		Float64s: []float64{1.1, 2.2, 3.3},
		Bool:     false,
		Bools:    []bool{true, false, false},
	}
	init.AddFlags(flags)
	init.Action(func() error {

		println("string:", fmt.Sprintf("%#v", flags.String))
		println("strings:", fmt.Sprintf("%#v", flags.Strings))
		println("strings_default:", fmt.Sprintf("%#v", flags.StringsDefault))
		println("\n")

		println("int:", fmt.Sprintf("%#v", flags.Int))
		println("ints:", fmt.Sprintf("%#v", flags.Ints))
		println("ints_default:", fmt.Sprintf("%#v", flags.IntsDefault))
		println("\n")

		println("int64:", fmt.Sprintf("%#v", flags.Int64))
		println("int64s:", fmt.Sprintf("%#v", flags.Int64s))
		println("int64s_default:", fmt.Sprintf("%#v", flags.Int64sDefault))
		println("\n")

		println("uint:", fmt.Sprintf("%#v", flags.Uint))
		println("uints:", fmt.Sprintf("%#v", flags.Uints))
		println("uints_default:", fmt.Sprintf("%#v", flags.UintsDefault))
		println("\n")

		println("uint64:", fmt.Sprintf("%#v", flags.Uint))
		println("uint64s:", fmt.Sprintf("%#v", flags.Uint64s))
		println("uint64s_default:", fmt.Sprintf("%#v", flags.Uint64sDefault))
		println("\n")

		println("float64:", fmt.Sprintf("%#v", flags.Float64))
		println("float64s:", fmt.Sprintf("%#v", flags.Float64s))
		println("float64s_default:", fmt.Sprintf("%#v", flags.Float64sDefault))
		println("\n")

		println("bool:", fmt.Sprintf("%#v", flags.Bool))
		println("bools:", fmt.Sprintf("%#v", flags.Bools))
		println("bools_default:", fmt.Sprintf("%#v", flags.BoolsDefault))
		return nil
	})

	// Run!
	if err := cli.Run(); err != nil {
		panic(err)
	}

}
