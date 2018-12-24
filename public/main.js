window.GoWasm = {
  functions: {},
  onAddFunction(name) {
    console.log('add wasm function: ', name)
  },
  onLoad () {
    console.log('wasm functions loaded')

    this.functions.run({
      done: data => {
        console.log(data)
      },
      fail: console.error
    })
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
      .catch(console.error)
  }
}

GoWasm.init()

const canvas = document.getElementById('canvas')
const ctx = canvas.getContext('2d')
const startButton = document.getElementById('start')
startButton.onclick = start
document.getElementById('step').onclick = step

let moving = false
function start() {
  if (moving) {
    moving = false
    startButton.textContent = 'start'
    return
  }

  function loop() {
    if (!moving) {
      return
    }
    step()
    setTimeout(() => { loop() }, 1000 / 60)
  }

  moving = true
  startButton.textContent = 'stop'
  loop()
}

function step() {
  GoWasm.functions.step({
    done: shapes => {
      console.log(shapes)
      ctx.clearRect(0, 0, 1000, 1000)
      shapes.forEach(shape => {
        const x = shape.x
        const y = shape.y
        shape.units.forEach(unit => {
          ctx.beginPath()
          unit.points.forEach((p, i) => {
            if (i === 0) {
              ctx.moveTo(adjustX(x + p.x), adjustY(y + p.y))
            } else {
              ctx.lineTo(adjustX(x + p.x), adjustY(y + p.y))
            }
          })
          if (unit.type === 'polygon') {
            ctx.closePath()
            ctx.fill()
          } else {
            ctx.stroke()
          }
        })
      })
    },
    fail: console.error
  })
}

function adjustX(v) {
  return v * 10
}

function adjustY(v) {
  return -v * 10 + canvas.height
}

