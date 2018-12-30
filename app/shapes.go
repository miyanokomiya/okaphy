package app

import (
	"syscall/js"

	"github.com/ByteArena/box2d"
)

// Shape 図形情報
type Shape struct {
	X     float64
	Y     float64
	Angle float64
}

// ToObject オブジェクト出力
func (s *Shape) ToObject() (v js.Value) {
	v = js.Global().Get("Object").New()
	v.Set("x", s.X)
	v.Set("y", s.Y)
	v.Set("angle", s.Angle)
	return v
}

func newShape(body *box2d.B2Body) (s *Shape) {
	s = &Shape{
		X:     body.GetPosition().X,
		Y:     body.GetPosition().Y,
		Angle: body.GetAngle(),
	}
	return s
}

func getPointFromObject(v js.Value) (p box2d.B2Vec2) {
	p = box2d.MakeB2Vec2(v.Get("x").Float(), v.Get("y").Float())
	return p
}

func getPointFromArray(array js.Value) (points []box2d.B2Vec2) {
	points = []box2d.B2Vec2{}
	for i := 0; i < array.Get("length").Int(); i++ {
		points = append(points, getPointFromObject(array.Index(i)))
	}
	return points
}

func newBodyFromObjects(a *app, value js.Value) {
	for i := 0; i < value.Get("length").Int(); i++ {
		obj := value.Index(i)
		newPolygonBody(a, obj)
	}
}

func newPolygonBody(a *app, value js.Value) {
	bd := box2d.MakeB2BodyDef()
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	body := a.world.CreateBody(&bd)

	units := value.Get("units")
	for i := 0; i < units.Get("length").Int(); i++ {
		unit := units.Index(i)
		points := getPointFromArray(unit.Get("points"))
		vertices := make([]box2d.B2Vec2, len(points))
		for j, p := range points {
			vertices[j].Set(p.X, p.Y)
		}
		shape := box2d.MakeB2PolygonShape()
		shape.Set(vertices, len(points))

		body.CreateFixture(&shape, 20.0)
	}

	a.characters = append(a.characters, body)
}
