import * as geo from 'okageo/src/geo'
import * as svg from 'okageo/src/svg'
import { IVec2 } from 'okageo/types'

interface IShape {
  x: number
  y: number
  angle: number
  units: IUnit[]
}

interface IUnit {
  points: IVec2[]
}

interface IWebAssembly {
  instantiateStreaming: (a1: any, a2: any) => any
  instantiate: (a1: any, a2: any) => any
}
const MyWebAssembly = (window as any).WebAssembly as IWebAssembly

interface IGoArgs {
  data?: any
  done?: (data: any) => void
  fail?: () => void
}
const goFunc = (args: IGoArgs) => { console.log(args) }

const goWasm: any = {
  functions: {
    add: goFunc,
    run: goFunc,
    step: goFunc
  },
  onAddFunction (name: string) { console.log('add wasm function: ', name) },
  onLoad () {
    console.log('wasm functions loaded')

    this.functions.run({
      done: (data: any) => { console.log(data) },
      fail: console.error
    })
  },
  init () {
    if (!MyWebAssembly.instantiateStreaming) {
      // polyfill
      MyWebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer
        return MyWebAssembly.instantiate(source, importObject)
      }
    }

    const go = new ((window as any).Go)()
    MyWebAssembly.instantiateStreaming(fetch('main.wasm'), go.importObject)
      .then((result: any) => { go.run(result.instance) })
      .catch(console.error)
  }
}

goWasm.init()

interface IWindow { GoWasm: any }
declare var window: IWindow
window.GoWasm = goWasm

let shapeList: IShape[] = []

const fileInput = document.getElementById('input') as HTMLInputElement
fileInput.onchange = (e) => {
  const file = (e.target as HTMLInputElement).files
  if (!file || file.length === 0) return

  const reader = new FileReader()
  reader.readAsText(file[0])
  reader.onload = () => {
    const pathInfoList = svg.parseSvgGraphicsStr(reader.result as string)
    const inRectList = svg.fitRect(pathInfoList, 0, 0, canvas.width, canvas.height)
    shapeList = inRectList.map((info) => ({
      angle: 0,
      units: geo.triangleSplit(info.d).map((points) => ({ points })),
      x: 0,
      y: 0
    }))
    goWasm.functions.add({
      data: shapeList,
      done: (data: any) => { console.log(data) },
      fail: console.error
    })
  }
}

const canvas = document.getElementById('canvas') as HTMLCanvasElement
const ctx = canvas.getContext('2d')
const startButton = document.getElementById('start') as HTMLInputElement
startButton.onclick = start
const stepButton = document.getElementById('step') as HTMLInputElement
stepButton.onclick = step

let moving = false
function start () {
  if (moving) {
    moving = false
    startButton.textContent = 'start'
    return
  }

  function loop () {
    if (!moving) {
      return
    }
    step()
    setTimeout(() => {
      loop()
    }, 1000 / 30)
  }

  moving = true
  startButton.textContent = 'stop'
  loop()
}

function step () {
  goWasm.functions.step({
    done: (shapes: IShape[]) => {
      if (!ctx) return
      ctx.clearRect(0, 0, canvas.width, canvas.height)

      shapeList = shapeList.map((old, i) => {
        const next = shapes[i]
        return { ...old, x: next.x, y: next.y, angle: next.angle }
      })

      shapeList.forEach((shape: IShape) => {
        ctx.save()
        ctx.translate(shape.x, shape.y)
        ctx.rotate(shape.angle)

        shape.units.forEach((unit) => {
          ctx.beginPath()
          unit.points.forEach((p: IVec2, j: number) => {
            if (j === 0) {
              ctx.moveTo(p.x, p.y)
            } else {
              ctx.lineTo(p.x, p.y)
            }
          })
          ctx.closePath()
          ctx.fill()
          ctx.stroke()
        })

        ctx.restore()
      })
    },
    fail: console.error
  })
}
