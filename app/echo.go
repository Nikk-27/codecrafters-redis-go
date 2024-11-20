package main

type Value struct {
    typ  string
    str  string
    bulk string
}

var Commands = map[string]func([]Value) Value{"ECHO":echo,
}

func echo(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'echo' command"}
	}
	return Value{typ: "string", str: args[0].bulk}
}