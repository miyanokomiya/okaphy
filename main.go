package main

import (
	"fmt"
	"syscall/js"

	"github.com/miyanokomiya/okaphy/app"
)

var goWasm = js.Global().Get("GoWasm")
var a app.App

func main() {
	c := make(chan struct{}, 0)

	wasmWrapper("echo", echo)
	wasmWrapper("run", run)
	wasmWrapper("step", step)
	wasmWrapper("add", add)

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

		arg := values[0]

		data := arg.Get("data")
		res, err := fn(data)

		if err != nil {
			fail := arg.Get("fail")
			if fail != js.Undefined() {
				fail.Invoke(err.Error())
			}
			return
		}

		done := arg.Get("done")
		if done != js.Undefined() {
			done.Invoke(res)
		}
	})

	goWasm.Get("functions").Set(name, handler)
	goWasm.Call("onAddFunction", name)
}

func echo(value js.Value) (interface{}, error) {
	return value, nil
}

func run(value js.Value) (interface{}, error) {
	a = app.NewApp()
	a.Run()
	return nil, nil
}

func step(value js.Value) (interface{}, error) {
	a.Step()
	array := js.Global().Get("Array").New()
	for _, shape := range a.GetShapes() {
		array.Call("push", shape.ToObject())
	}
	return array, nil
}

func add(value js.Value) (interface{}, error) {
	a.AddShapes(value)
	array := js.Global().Get("Array").New()
	for _, shape := range a.GetShapes() {
		array.Call("push", shape.ToObject())
	}
	return array, nil
}
