window.GoWasm = {
  functions: {},
  onAddFunction(name) {
    console.log('add wasm function: ', name)
    const createButton = (name) => {
      const btn = document.createElement('button')
      btn.type = 'button'
      btn.onclick = () => {
        this.functions[name]({
          data: { a: 1, b: 2 },
          done: (arg) => {
            console.log('done: ', arg)
          },
          fail: (err) => {
            console.error(err)
          }
        })
      }
      btn.textContent = name
      document.body.appendChild(btn)
    }
    createButton(name)
  },
  onLoad () {
    document.getElementById("wasmReady").textContent = 'wasm ready'
    console.log('wasm functions loaded')
  },
  init () {
    if (!WebAssembly.instantiateStreaming) { // polyfill
      WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer
        return await WebAssembly.instantiate(source, importObject)
      }
    }

    const go = new Go()
    let mod, inst
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
      .then(result => {
        mod = result.module
        inst = result.instance
        go.run(inst)
      })
      .catch(err => {
        console.error(err)
      })
  }
}

GoWasm.init()

