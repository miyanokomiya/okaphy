package app

import (
	"syscall/js"

	"github.com/ByteArena/box2d"
)

var timeStep = 1.0 / 60.0
var velocityIterations = 8
var positionIterations = 3

// App アプリ本体
type App interface {
	Run()
	Step()
	GetShapes() []Shape
	AddShapes(value js.Value)
}

type app struct {
	gravity    box2d.B2Vec2
	world      box2d.B2World
	characters []*box2d.B2Body
}

// NewApp アプリ作成
func NewApp() App {
	gravity := box2d.MakeB2Vec2(0.0, -50.0)
	return &app{
		gravity:    gravity,
		world:      box2d.MakeB2World(gravity),
		characters: []*box2d.B2Body{},
	}
}

func (a *app) AddShapes(value js.Value) {
	newBodyFromObjects(a, value)
}

func (a *app) GetShapes() []Shape {
	shapes := []Shape{}
	for _, character := range a.characters {
		shapes = append(shapes, *newShape(character))
	}
	return shapes
}

func (a *app) Step() {
	a.world.Step(timeStep, velocityIterations, positionIterations)
}

func (a *app) Run() {
	// Ground body
	{
		bd := box2d.MakeB2BodyDef()
		ground := a.world.CreateBody(&bd)

		shape := box2d.MakeB2EdgeShape()
		shape.Set(box2d.MakeB2Vec2(-1000.0, 0.0), box2d.MakeB2Vec2(2000.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
	}
}
