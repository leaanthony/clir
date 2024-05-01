package main

import (
	"fmt"
	"github.com/leaanthony/clir"
	"reflect"
)

type Flags struct {
	String         string   `name:"string" description:"The string" pos:"1"`
	Strings        []string `name:"strings" description:"The strings" pos:"2"`
	StringsDefault []string `name:"strings_default" description:"The strings default" default:"one,two,three" pos:"3"`

	Int         int   `name:"int" description:"The int" pos:"4"`
	Ints        []int `name:"ints" description:"The ints" pos:"5"`
	IntsDefault []int `name:"ints_default" description:"The ints default" default:"3,4,5" pos:"6"`

	Int64         int64   `name:"int64" description:"The int64" pos:"7"`
	Int64s        []int64 `name:"int64s" description:"The int64s" pos:"8"`
	Int64sDefault []int64 `name:"int64s_default" description:"The int64s default" default:"3,4,5" pos:"9"`

	Uint         uint   `name:"uint" description:"The uint" pos:"10"`
	Uints        []uint `name:"uints" description:"The uints" pos:"11"`
	UintsDefault []uint `name:"uints_default" description:"The uints default" default:"3,4,5" pos:"12"`

	Uint64         uint64   `name:"uint64" description:"The uint64" pos:"13"`
	Uint64s        []uint64 `name:"uint64s" description:"The uint64s" pos:"14"`
	Uint64sDefault []uint64 `name:"uint64s_default" description:"The uint64s default" default:"3,4,5" pos:"15"`

	Float64         float64   `name:"float64" description:"The float64" pos:"17"`
	Float64s        []float64 `name:"float64s" description:"The float64s" pos:"18"`
	Float64sDefault []float64 `name:"float64s_default" description:"The float64s default" default:"3,4,5" pos:"19"`

	Bool         bool   `name:"bool" description:"The bool" pos:"20"`
	Bools        []bool `name:"bools" description:"The bools" pos:"21"`
	BoolsDefault []bool `name:"bools_default" description:"The bools default" default:"false,true,false,true" pos:"22"`
}

func main() {

	// Create new cli
	cli := clir.NewCli("flagstruct", "An example of subcommands with flag inherence", "v0.0.1")

	// Create an init subcommand with flag inheritance
	init := cli.NewSubCommand("flag", "print default")

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
		println("\n")

		return nil
	})

	// Run!
	if err := cli.Run("flag"); err != nil {
		panic(err)
	}

	cli.NewSubCommandFunction("positional", "test positional args", func(f *Flags) error {

		if f.String != "hello" {
			panic(fmt.Sprintf("expected 'hello', got '%v'", f.String))
		}

		if !reflect.DeepEqual(f.Strings, []string{"zkep", "hello", "clir"}) {
			panic(fmt.Sprintf("expected '[zkep hello clir]', got '%v'", f.Strings))
		}

		if !reflect.DeepEqual(f.StringsDefault, []string{"zkep", "clir", "hello"}) {
			panic(fmt.Sprintf("expected '[zkep clir hello]', got '%v'", f.StringsDefault))
		}

		println("string:", fmt.Sprintf("%#v", f.String))
		println("strings:", fmt.Sprintf("%#v", f.Strings))
		println("strings_default:", fmt.Sprintf("%#v", f.StringsDefault))
		println("\n")

		return nil
	})

	// Run!
	if err := cli.Run("positional", "hello", "zkep,hello,clir", "zkep,clir,hello"); err != nil {
		panic(err)
	}

}
