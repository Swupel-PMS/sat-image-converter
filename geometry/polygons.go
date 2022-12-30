package geometry

type Polygons struct {
	Bisector *Line

	LeftTriangle  *Polygon
	Trapezoid     *Polygon
	RightTriangle *Polygon

	P1Exist bool
	P2Exist bool
	P3Exist bool
	P4Exist bool

	LeftTriangleSquare  float64
	TrapezoidSquare     float64
	RightTriangleSquare float64
	TotalSquare         float64
}

func NewEmptyPolygons() *Polygons {
	return &Polygons{
		Bisector:      NewEmptyLine(),
		LeftTriangle:  NewEmptyPolygon(),
		Trapezoid:     NewEmptyPolygon(),
		RightTriangle: NewEmptyPolygon(),
	}
}

func NewPolygonsFromLine(l1 *Line, l2 *Line) *Polygons {
	res := NewEmptyPolygons()
	res.Bisector = GetBisector(l1, l2)
	v1 := l1.Start
	v2 := l1.End
	v3 := l2.Start
	v4 := l2.End

	res.P1Exist = false
	res.P4Exist = false

	if !v1.Equal(v4) {
		l1s := NewLineByVector(v1, res.Bisector.GetLineNearestPoint(v1))
		p1, ok := l1s.CrossLineSegment(l2)
		res.P1Exist = ok && !p1.Equal(v4)
		if res.P1Exist {
			res.LeftTriangle.Add(v1)
			res.LeftTriangle.Add(v4)
			res.LeftTriangle.Add(p1)
			res.Trapezoid.Add(p1)
		} else {
			res.Trapezoid.Add(v4)
		}
		l2e := NewLineByVector(v4, res.Bisector.GetLineNearestPoint(v4))
		p4, ok := l2e.CrossLineSegment(l1)
		res.P4Exist = ok && !p4.Equal(v1)
		if res.P4Exist {
			res.LeftTriangle.Add(v4)
			res.LeftTriangle.Add(v1)
			res.LeftTriangle.Add(p4)
			res.Trapezoid.Add(p4)
		} else {
			res.Trapezoid.Add(v1)
		}
	} else {
		res.Trapezoid.Add(v4)
		res.Trapezoid.Add(v1)
	}
	res.P2Exist = false
	res.P3Exist = false

	if !v2.Equal(v3) {
		l2s := NewLineByVector(v3, res.Bisector.GetLineNearestPoint(v3))
		p3, ok := l2s.CrossLineSegment(l1)
		res.P3Exist = ok && !p3.Equal(v2)
		if res.P3Exist {
			res.RightTriangle.Add(v3)
			res.RightTriangle.Add(v2)
			res.RightTriangle.Add(p3)

			res.Trapezoid.Add(p3)
		} else {
			res.Trapezoid.Add(v2)
		}
		l1e := NewLineByVector(v2, res.Bisector.GetLineNearestPoint(v2))
		p2, ok := l1e.CrossLineSegment(l2)
		res.P2Exist = ok && !p2.Equal(v3)
		if res.P2Exist {
			res.RightTriangle.Add(v2)
			res.RightTriangle.Add(v3)
			res.RightTriangle.Add(p2)

			res.Trapezoid.Add(p2)
		} else {
			res.Trapezoid.Add(v3)
		}
	} else {
		res.Trapezoid.Add(v2)
		res.Trapezoid.Add(v3)
	}
	res.LeftTriangleSquare = res.LeftTriangle.CountSquare()
	res.TrapezoidSquare = res.Trapezoid.CountSquare()
	res.RightTriangleSquare = res.RightTriangle.CountSquare()

	res.TotalSquare = res.LeftTriangleSquare + res.TrapezoidSquare + res.RightTriangleSquare
	return res
}
