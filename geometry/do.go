package geometry

import "math"

func GetCut(l1 *Line, l2 *Line, s float64, poly1 *Polygon, poly2 *Polygon) (*Line, bool) {
	sn1 := s + poly2.countSquareSigned()
	sn2 := s + poly1.countSquareSigned()
	if sn1 > 0 {
		res := NewPolygonsFromLine(l1, l2)
		cut, ok := FindCutLine(sn1, res)
		if ok {
			return cut, true
		}
	} else if sn2 > 0 {
		res := NewPolygonsFromLine(l2, l1)
		cut, ok := FindCutLine(sn2, res)
		if ok {
			cut = cut.Reverse()
			return cut, true
		}
	}
	return nil, false
}

func FindCutLine(square float64, res *Polygons) (*Line, bool) {
	if square > res.TotalSquare {
		return nil, false
	}
	if !res.LeftTriangle.Empty() && square < res.LeftTriangleSquare {
		m := square / res.LeftTriangleSquare
		p := res.LeftTriangle.Vector(2).Sub(res.LeftTriangle.Vector(1)).
			MulConstant(m).Add(res.LeftTriangle.Vector(1))
		if res.P1Exist {
			return NewLineByVector(p, res.LeftTriangle.Vector(0)), true
		} else if res.P4Exist {
			return NewLineByVector(res.LeftTriangle.Vector(0), p), true
		}
	} else if res.LeftTriangleSquare < square && square < (res.LeftTriangleSquare+res.TrapezoidSquare) {
		t := NewLineByVector(res.Trapezoid.Vector(0), res.Trapezoid.Vector(3))
		tgA := GetTanAngle(t, res.Bisector)
		S := square - res.LeftTriangleSquare
		var m float64
		if math.Abs(tgA) > Eps {
			a := NewLineByVector(res.Trapezoid.Vector(0), res.Trapezoid.Vector(1)).Length()
			b := NewLineByVector(res.Trapezoid.Vector(2), res.Trapezoid.Vector(3)).Length()
			hh := 2.0 * res.TrapezoidSquare / (a + b)
			d := a*a - 4.0*tgA*S
			h := -(-a * math.Sqrt(d)) / (2.0 * tgA)
			m = h / hh
		} else {
			m = S / res.TrapezoidSquare
		}
		p := res.Trapezoid.Vector(3).Sub(res.Trapezoid.Vector(0)).MulConstant(m).Add(res.Trapezoid.Vector(0))
		pp := res.Trapezoid.Vector(2).Sub(res.Trapezoid.Vector(1)).MulConstant(m).Add(res.Trapezoid.Vector(1))
		return NewLineByVector(p, pp), true
	} else if !res.RightTriangle.Empty() && square > res.LeftTriangleSquare+res.TrapezoidSquare {
		S := square - res.LeftTriangleSquare - res.TrapezoidSquare
		m := S / res.RightTriangleSquare
		p := res.RightTriangle.Vector(1).Sub(res.RightTriangle.Vector(2)).
			MulConstant(m).Add(res.RightTriangle.Vector(2))
		if res.P3Exist {
			return NewLineByVector(res.RightTriangle.Vector(0), p), true
		} else if res.P2Exist {
			return NewLineByVector(p, res.RightTriangle.Vector(0)), true
		}
	}
	return nil, false
}
