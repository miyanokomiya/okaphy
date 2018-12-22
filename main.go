package main

import (
	"errors"
	"fmt"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	wasmWrapper("getHello", getHello)
	wasmWrapper("get100", get100)
	wasmWrapper("echo", echo)

	js.Global().Get("GoWasm").Call("onLoad")
	<-c
}

func wasmWrapper(name string, fn func(value js.Value) (interface{}, error)) {
	handler := js.NewCallback(func(i []js.Value) {
		fmt.Println("exec " + name)

		if len(i) == 0 {
			var val js.Value
			fn(val)
			return
		}

		data := i[0].Get("data")
		res, err := fn(data)

		if err != nil {
			fail := i[0].Get("fail")
			if fail != js.Undefined() {
				fail.Invoke(err.Error())
			}
			return
		}

		done := i[0].Get("done")
		if done != js.Undefined() {
			done.Invoke(res)
		}
	})

	js.Global().Get("GoWasm").Get("functions").Set(name, handler)
	js.Global().Get("GoWasm").Call("onAddFunction", name)
}

func getHello(value js.Value) (interface{}, error) {
	return fmt.Sprintf("Hello " + value.String()), nil
}

func get100(value js.Value) (interface{}, error) {
	return 100, nil
}

func echo(value js.Value) (interface{}, error) {
	return value, errors.New("error")
}
