package clir

type flagFunc func(c *Command)

func BoolFlag(name, description string, variable *bool) flagFunc {
	return func(c *Command) {
		c.flags.BoolVar(variable, name, *variable, description)
	}
}

func BoolsFlag(name, description string, variable *[]bool) flagFunc {
	return func(c *Command) {
		c.flags.Var(newBoolsValue(*variable, variable), name, description)
	}
}

func StringFlag(name, description string, variable *string) flagFunc {
	return func(c *Command) {
		c.flags.StringVar(variable, name, *variable, description)
	}
}

func StringsFlag(name, description string, variable *[]string) flagFunc {
	return func(c *Command) {
		c.flags.Var(newStringsValue(*variable, variable), name, description)
	}
}

func IntFlag(name, description string, variable *int) flagFunc {
	return func(c *Command) {
		c.flags.IntVar(variable, name, *variable, description)
	}
}

func IntsFlag(name, description string, variable *[]int) flagFunc {
	return func(c *Command) {
		c.flags.Var(newIntsValue(*variable, variable), name, description)
	}
}

func Int8Flag(name, description string, variable *int8) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt8Value(*variable, variable), name, description)
	}
}

func Int8sFlag(name, description string, variable *[]int8) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt8sValue(*variable, variable), name, description)
	}
}

func Int16Flag(name, description string, variable *int16) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt16Value(*variable, variable), name, description)
	}
}

func Int16sFlag(name, description string, variable *[]int16) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt16sValue(*variable, variable), name, description)
	}
}

func Int32Flag(name, description string, variable *int32) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt32Value(*variable, variable), name, description)
	}
}

func Int32sFlag(name, description string, variable *[]int32) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt32sValue(*variable, variable), name, description)
	}
}

func Int64Flag(name, description string, variable *int64) flagFunc {
	return func(c *Command) {
		c.flags.Int64Var(variable, name, *variable, description)
	}
}

func Int64sFlag(name, description string, variable *[]int64) flagFunc {
	return func(c *Command) {
		c.flags.Var(newInt64sValue(*variable, variable), name, description)
	}
}

func UintFlag(name, description string, variable *uint) flagFunc {
	return func(c *Command) {
		c.flags.UintVar(variable, name, *variable, description)
	}
}

func UintsFlag(name, description string, variable *[]uint) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUintsValue(*variable, variable), name, description)
	}
}

func Uint8Flag(name, description string, variable *uint8) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint8Value(*variable, variable), name, description)
	}
}

func Uint8sFlag(name, description string, variable *[]uint8) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint8sValue(*variable, variable), name, description)
	}
}

func Uint16Flag(name, description string, variable *uint16) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint16Value(*variable, variable), name, description)
	}
}

func Uint16sFlag(name, description string, variable *[]uint16) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint16sValue(*variable, variable), name, description)
	}
}

func Uint32Flag(name, description string, variable *uint32) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint32Value(*variable, variable), name, description)
	}
}

func Uint32sFlag(name, description string, variable *[]uint32) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint32sValue(*variable, variable), name, description)
	}
}

func Uint64Flag(name, description string, variable *uint64) flagFunc {
	return func(c *Command) {
		c.flags.Uint64Var(variable, name, *variable, description)
	}
}

func Uint64sFlag(name, description string, variable *[]uint64) flagFunc {
	return func(c *Command) {
		c.flags.Var(newUint64sValue(*variable, variable), name, description)
	}
}

func Float64Flag(name, description string, variable *float64) flagFunc {
	return func(c *Command) {
		c.flags.Float64Var(variable, name, *variable, description)
	}
}

func Float64sFlag(name, description string, variable *[]float64) flagFunc {
	return func(c *Command) {
		c.flags.Var(newFloat64sValue(*variable, variable), name, description)
	}
}

func Float32Flag(name, description string, variable *float32) flagFunc {
	return func(c *Command) {
		c.flags.Var(newFloat32Value(*variable, variable), name, description)
	}
}

func Float32sFlag(name, description string, variable *[]float32) flagFunc {
	return func(c *Command) {
		c.flags.Var(newFloat32sValue(*variable, variable), name, description)
	}
}
