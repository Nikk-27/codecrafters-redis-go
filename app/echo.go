package main

// Commands map that maps command names to handler functions
var Commands = map[string]func([]Value) Value{
    "ECHO": echo,
}

// echo function to handle the ECHO command
func echo(args []Value) Value {
    if len(args) != 1 {
        return Value{typ: "error", bulk: "ERR wrong number of arguments for 'echo' command"}
    }
    return Value{typ: "string", bulk: args[0].bulk}
}
