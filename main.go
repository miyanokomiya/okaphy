package main

import (
	"fmt"
	"syscall/js"
)

var goWasm = js.Global().Get("GoWasm")

func main() {
	c := make(chan struct{}, 0)

	wasmWrapper("getHello", getHello)
	wasmWrapper("get100", get100)
	wasmWrapper("echo", echo)

	goWasm.Call("onLoad")
	<-c
}

func wasmWrapper(name string, fn func(value js.Value) (interface{}, error)) {
	handler := js.NewCallback(func(values []js.Value) {
		fmt.Println("exec " + name)

		if len(values) == 0 {
			var val js.Value
			fn(val)
			return
		}

		data := values[0].Get("data")
		res, err := fn(data)

		if err != nil {
			fail := values[0].Get("fail")
			if fail != js.Undefined() {
				fail.Invoke(err.Error())
			}
			return
		}

		done := values[0].Get("done")
		if done != js.Undefined() {
			done.Invoke(res)
		}
	})

	goWasm.Get("functions").Set(name, handler)
	goWasm.Call("onAddFunction", name)
}

func getHello(value js.Value) (interface{}, error) {
	return fmt.Sprintf("Hello " + value.String()), nil
}

func get100(value js.Value) (interface{}, error) {
	return 100, nil
}

func echo(value js.Value) (interface{}, error) {
	return value, nil
}
