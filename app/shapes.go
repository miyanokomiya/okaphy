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
	Units []UnitInterface
}

// ToObject オブジェクト出力
func (s *Shape) ToObject() (v js.Value) {
	v = js.Global().Get("Object").New()
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

func newShape(body *box2d.B2Body) (s *Shape) {
	s = &Shape{
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

func getPointArray(points []box2d.B2Vec2) (array js.Value) {
	array = js.Global().Get("Array").New()
	for _, p := range points {
		array.Call("push", getPointObject(p))
	}
	return array
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
