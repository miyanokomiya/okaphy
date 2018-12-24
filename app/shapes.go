package app

import (
	"syscall/js"

	"github.com/ByteArena/box2d"
)

// Shape 図形情報
type Shape struct {
	ID    string
	X     float64
	Y     float64
	Angle float64
	Units []UnitInterface
}

// ToObject オブジェクト出力
func (s *Shape) ToObject() (v js.Value) {
	v = js.Global().Get("Object").New()
	v.Set("id", s.ID)
	v.Set("x", s.X)
	v.Set("y", s.Y)
	v.Set("angle", s.Angle)
	units := js.Global().Get("Array").New()
	for _, u := range s.Units {
		units.Call("push", u.ToObject())
	}
	v.Set("units", units)
	return v
}

func newShape(body *box2d.B2Body, id string) (s *Shape) {
	s = &Shape{
		ID:    id,
		X:     body.GetPosition().X,
		Y:     body.GetPosition().Y,
		Angle: body.GetAngle(),
		Units: getUnits(body),
	}
	return s
}

func getPointObject(p box2d.B2Vec2) (v js.Value) {
	v = js.Global().Get("Object").New()
	v.Set("x", p.X)
	v.Set("y", p.Y)
	return v
}

func getPointFromObject(v js.Value) (p box2d.B2Vec2) {
	p = box2d.MakeB2Vec2(v.Get("x").Float(), v.Get("y").Float())
	return p
}

func getPointArray(points []box2d.B2Vec2) (array js.Value) {
	array = js.Global().Get("Array").New()
	for _, p := range points {
		array.Call("push", getPointObject(p))
	}
	return array
}

func getPointFromArray(array js.Value) (points []box2d.B2Vec2) {
	points = []box2d.B2Vec2{}
	for i := 0; i < array.Get("length").Int(); i++ {
		points = append(points, getPointFromObject(array.Index(i)))
	}
	return points
}

// UnitInterface 図形単位インタフェース
type UnitInterface interface {
	ToObject() js.Value
}

// EdgeUnit 単位エッジ
type EdgeUnit struct {
	points Points
}

// ToObject オブジェクト出力
func (u *EdgeUnit) ToObject() (v js.Value) {
	v = js.Global().Get("Object").New()
	v.Set("type", "edge")
	points := getPointArray(u.points)
	v.Set("points", points)
	return v
}

// PolygonUnit 単位ポリゴン
type PolygonUnit struct {
	points Points
}

// ToObject オブジェクト出力
func (u *PolygonUnit) ToObject() (v js.Value) {
	v = js.Global().Get("Object").New()
	v.Set("type", "polygon")
	points := getPointArray(u.points)
	v.Set("points", points)
	return v
}

// Points 座標リスト
type Points []box2d.B2Vec2

func getUnits(body *box2d.B2Body) (units []UnitInterface) {
	units = []UnitInterface{}
	for fixture := body.GetFixtureList(); fixture != nil; fixture = fixture.GetNext() {
		if fixture.GetType() == box2d.B2Shape_Type.E_polygon {
			if shape, ok := fixture.GetShape().(*box2d.B2PolygonShape); ok {
				units = append(units, getUnitFromPolygon(shape))
			}
		} else if fixture.GetType() == box2d.B2Shape_Type.E_edge {
			if shape, ok := fixture.GetShape().(*box2d.B2EdgeShape); ok {
				units = append(units, getUnitFromEdge(shape))
			}
		}
	}
	return units
}

func getUnitFromPolygon(shape *box2d.B2PolygonShape) (unit UnitInterface) {
	points := []box2d.B2Vec2{}
	// M_countに従ってループ
	for i := 0; i < shape.M_count; i++ {
		points = append(points, shape.M_vertices[i])
	}
	unit = &PolygonUnit{
		points: points,
	}
	return unit
}

func getUnitFromEdge(shape *box2d.B2EdgeShape) (unit UnitInterface) {
	points := []box2d.B2Vec2{
		shape.M_vertex1,
		shape.M_vertex2,
	}
	unit = &EdgeUnit{
		points: points,
	}
	return unit
}

func newBodyFromObjects(a *app, value js.Value) {
	for i := 0; i < value.Get("length").Int(); i++ {
		obj := value.Index(i)
		objType := obj.Get("type").String()
		if objType == "polygon" {
			newPolygonBody(a, obj)
		}
	}
}

func newPolygonBody(a *app, value js.Value) {
	bd := box2d.MakeB2BodyDef()
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.Position.Set(5.0, 8.0)
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

	a.characters[value.Get("id").String()] = body
}
