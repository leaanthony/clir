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

	Int8         int8   `name:"int8" description:"The int8"  pos:"7"`
	Int8s        []int8 `name:"int8s" description:"The int8s" pos:"8"`
	Int8sDefault []int8 `name:"int8s_default" description:"The int8s default" default:"3,4,5" pos:"9"`

	Int16         int16   `name:"int16" description:"The int16"  pos:"10"`
	Int16s        []int16 `name:"int16s" description:"The int16s" pos:"11"`
	Int16sDefault []int16 `name:"int16s_default" description:"The int16s default" default:"3,4,5" pos:"12"`

	Int32         int32   `name:"int32" description:"The int32"  pos:"13"`
	Int32s        []int32 `name:"int32s" description:"The int32s" pos:"14"`
	Int32sDefault []int32 `name:"int32s_default" description:"The int32 default" default:"3,4,5" pos:"15"`

	Int64         int64   `name:"int64" description:"The int64" pos:"16"`
	Int64s        []int64 `name:"int64s" description:"The int64s" pos:"17"`
	Int64sDefault []int64 `name:"int64s_default" description:"The int64s default" default:"3,4,5" pos:"18"`

	Uint         uint   `name:"uint" description:"The uint" pos:"19"`
	Uints        []uint `name:"uints" description:"The uints" pos:"20"`
	UintsDefault []uint `name:"uints_default" description:"The uints default" default:"3,4,5" pos:"21"`

	Uint8         uint8   `name:"uint8" description:"The uint8" pos:"22"`
	Uint8s        []uint8 `name:"uint8s" description:"The uint8s" pos:"23"`
	Uint8sDefault []uint8 `name:"uint8s_default" description:"The uint8s default" default:"3,4,5" pos:"24"`

	Uint16         uint16   `name:"uint16" description:"The uint16" pos:"25"`
	Uint16s        []uint16 `name:"uint16s" description:"The uint16s" pos:"26"`
	Uint16sDefault []uint16 `name:"uint16s_default" description:"The uint16 default" default:"3,4,5" pos:"27"`

	Uint32         uint32   `name:"uint32" description:"The uint32" pos:"28"`
	Uint32s        []uint32 `name:"uint32s" description:"The uint32s" pos:"29"`
	Uint32sDefault []uint32 `name:"uint32s_default" description:"The uint32s default" default:"3,4,5" pos:"30"`

	Uint64         uint64   `name:"uint64" description:"The uint64" pos:"31"`
	Uint64s        []uint64 `name:"uint64s" description:"The uint64s" pos:"32"`
	Uint64sDefault []uint64 `name:"uint64s_default" description:"The uint64s default" default:"3,4,5" pos:"33"`

	Float32         float32   `name:"float32" description:"The float32" pos:"34"`
	Float32s        []float32 `name:"float32s" description:"The float32s" pos:"35"`
	Float32sDefault []float32 `name:"float32s_default" description:"The float32s default" default:"3,4,5" pos:"36"`

	Float64         float64   `name:"float64" description:"The float64" pos:"37"`
	Float64s        []float64 `name:"float64s" description:"The float64s" pos:"38"`
	Float64sDefault []float64 `name:"float64s_default" description:"The float64s default" default:"3,4,5" pos:"39"`

	Bool         bool   `name:"bool" description:"The bool" pos:"40"`
	Bools        []bool `name:"bools" description:"The bools" pos:"41"`
	BoolsDefault []bool `name:"bools_default" description:"The bools default" default:"false,true,false,true" pos:"42"`
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
		Int8:     int8(1),
		Int8s:    []int8{1, 2, 3},
		Int16:    int16(1),
		Int16s:   []int16{1, 2, 3},
		Int32:    int32(1),
		Int32s:   []int32{1, 2, 3},
		Int64:    int64(1),
		Int64s:   []int64{1, 2, 3},
		Uint:     uint(1),
		Uints:    []uint{1, 2, 3},
		Uint8:    uint8(1),
		Uint8s:   []uint8{1, 2, 3},
		Uint16:   uint16(1),
		Uint16s:  []uint16{1, 2, 3},
		Uint32:   uint32(1),
		Uint32s:  []uint32{1, 2, 3},
		Uint64:   uint64(1),
		Uint64s:  []uint64{1, 2, 3},
		Float32:  float32(3.14),
		Float32s: []float32{1.1, 2.2, 3.3},
		Float64:  float64(3.14),
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

		println("int8:", fmt.Sprintf("%#v", flags.Int8))
		println("int8s:", fmt.Sprintf("%#v", flags.Int8s))
		println("int8s_default:", fmt.Sprintf("%#v", flags.Int8sDefault))
		println("\n")

		println("int16:", fmt.Sprintf("%#v", flags.Int16))
		println("int16s:", fmt.Sprintf("%#v", flags.Int16s))
		println("int16s_default:", fmt.Sprintf("%#v", flags.Int16sDefault))
		println("\n")

		println("int32:", fmt.Sprintf("%#v", flags.Int32))
		println("int32s:", fmt.Sprintf("%#v", flags.Int32s))
		println("int32s_default:", fmt.Sprintf("%#v", flags.Int32sDefault))
		println("\n")

		println("int64:", fmt.Sprintf("%#v", flags.Int64))
		println("int64s:", fmt.Sprintf("%#v", flags.Int64s))
		println("int64s_default:", fmt.Sprintf("%#v", flags.Int64sDefault))
		println("\n")

		println("uint:", fmt.Sprintf("%#v", flags.Uint))
		println("uints:", fmt.Sprintf("%#v", flags.Uints))
		println("uints_default:", fmt.Sprintf("%#v", flags.UintsDefault))
		println("\n")

		println("uint8:", fmt.Sprintf("%#v", flags.Uint8))
		println("uint8s:", fmt.Sprintf("%#v", flags.Uint8s))
		println("uint8s_default:", fmt.Sprintf("%#v", flags.Uint8sDefault))
		println("\n")

		println("uint16:", fmt.Sprintf("%#v", flags.Uint16))
		println("uint16s:", fmt.Sprintf("%#v", flags.Uint16s))
		println("uint16s_default:", fmt.Sprintf("%#v", flags.Uint16sDefault))
		println("\n")

		println("uint32:", fmt.Sprintf("%#v", flags.Uint32))
		println("uint32s:", fmt.Sprintf("%#v", flags.Uint32s))
		println("uint32s_default:", fmt.Sprintf("%#v", flags.Uint32sDefault))
		println("\n")

		println("uint64:", fmt.Sprintf("%#v", flags.Uint))
		println("uint64s:", fmt.Sprintf("%#v", flags.Uint64s))
		println("uint64s_default:", fmt.Sprintf("%#v", flags.Uint64sDefault))
		println("\n")

		println("float32:", fmt.Sprintf("%#v", flags.Float32))
		println("float32s:", fmt.Sprintf("%#v", flags.Float32s))
		println("float32s_default:", fmt.Sprintf("%#v", flags.Float32sDefault))
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
