package app

import (
	"math"
	"sort"
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
	characters map[string]*box2d.B2Body
}

// NewApp アプリ作成
func NewApp() App {
	gravity := box2d.MakeB2Vec2(0.0, -10.0)
	return &app{
		gravity:    gravity,
		world:      box2d.MakeB2World(gravity),
		characters: make(map[string]*box2d.B2Body),
	}
}

func (a *app) AddShapes(value js.Value) {
	newBodyFromObjects(a, value)
}

func (a *app) GetShapes() []Shape {
	characterNames := make([]string, 0)
	for k := range a.characters {
		characterNames = append(characterNames, k)
	}
	sort.Strings(characterNames)

	shapes := []Shape{}
	var character *box2d.B2Body
	for _, name := range characterNames {
		character = a.characters[name]
		shapes = append(shapes, *newShape(character, name))
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
		shape.Set(box2d.MakeB2Vec2(0.0, 0.0), box2d.MakeB2Vec2(20.0, 0.0))
		ground.CreateFixture(&shape, 0.0)
		a.characters["00_ground"] = ground
	}

	// Square tiles. This shows that adjacency shapes may
	// have non-smooth  There is no solution
	// to this problem.
	{
		bd := box2d.MakeB2BodyDef()
		ground := a.world.CreateBody(&bd)

		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBoxFromCenterAndAngle(1.0, 1.0, box2d.MakeB2Vec2(4.0, 3.0), 0.0)
		ground.CreateFixture(&shape, 0.0)
		shape.SetAsBoxFromCenterAndAngle(1.0, 1.0, box2d.MakeB2Vec2(6.0, 3.0), 0.0)
		ground.CreateFixture(&shape, 0.0)
		shape.SetAsBoxFromCenterAndAngle(1.0, 1.0, box2d.MakeB2Vec2(8.0, 3.0), 0.0)
		ground.CreateFixture(&shape, 0.0)
		a.characters["03_squaretiles"] = ground
	}

	// Square character 1
	{
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(3.0, 8.0)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.FixedRotation = true
		bd.AllowSleep = false

		body := a.world.CreateBody(&bd)

		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBox(0.5, 0.5)

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 20.0
		body.CreateFixtureFromDef(&fd)
		a.characters["06_squarecharacter1"] = body
	}

	// Hexagon character
	{
		bd := box2d.MakeB2BodyDef()
		bd.Position.Set(5.0, 8.0)
		bd.Type = box2d.B2BodyType.B2_dynamicBody
		bd.FixedRotation = true
		bd.AllowSleep = false

		body := a.world.CreateBody(&bd)

		angle := 0.0
		delta := box2d.B2_pi / 3.0
		vertices := make([]box2d.B2Vec2, 6)
		for i := 0; i < 6; i++ {
			vertices[i].Set(0.5*math.Cos(angle), 0.5*math.Sin(angle))
			angle += delta
		}

		shape := box2d.MakeB2PolygonShape()
		shape.Set(vertices, 6)

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &shape
		fd.Density = 20.0
		body.CreateFixtureFromDef(&fd)
		a.characters["08_hexagoncharacter"] = body
	}
}
