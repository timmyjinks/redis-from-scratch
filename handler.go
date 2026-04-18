package main

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

func set(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "error", str: "PANIC YOU ARE GOING TO DIE"}
	}

	key := args[0].bulk
	val := args[0].bulk

	DataMutex.Lock()
	Data[key] = val
	DataMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "error", str: "PANIC YOU ARE GOING TO DIE"}
	}

	key := args[0].bulk

	DataMutex.Lock()
	val, ok := Data[key]
	DataMutex.Unlock()
	if !ok {
		return Value{typ: "error", str: "value DNE"}
	}

	return Value{typ: "string", str: val}
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	val := args[0].bulk

	return Value{typ: "string", str: val}
}
