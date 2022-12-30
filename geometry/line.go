package geometry

// checked
import (
	"math"
)

type Line struct {
	Start, End *Vector
	A, B, C    float64
}

func NewEmptyLine() *Line {
	return &Line{
		Start: &Vector{},
		End:   &Vector{},
		A:     0,
		B:     0,
		C:     0,
	}
}

func NewLineByValues(a, b, c float64) *Line {
	start := &Vector{}
	end := &Vector{}
	if math.Abs(a) <= Eps && math.Abs(b) >= Eps {
		start.X = -1000
		start.Y = -(c / b)

		end.X = 1000
		end.Y = start.Y
	} else if math.Abs(b) <= Eps && math.Abs(a) >= Eps {
		start.X = -(c / a)
		start.Y = -1000
		end.X = start.X
		end.Y = 1000
	} else {
		start.X = -1000
		start.Y = -((a*start.X + c) / b)
		end.X = 1000
		end.Y = -((a*end.X + c) / b)
	}
	return &Line{
		Start: start,
		End:   end,
		A:     a,
		B:     b,
		C:     c,
	}
}

func NewLineByVector(start, end *Vector) *Line {
	return &Line{
		Start: start,
		End:   end,
		A:     start.Y - end.Y,
		B:     end.X - start.X,
		C:     start.X*end.Y - end.X*start.Y,
	}
}

func NewDirectedLine(p *Vector, d *Vector) *Line {
	return NewLineByVector(p, p.Add(d))
}

func (l *Line) GetDistance(point *Vector) float64 {
	n := l.A*point.X + l.B*point.Y + l.C
	m := math.Sqrt(l.A*l.A + l.B*l.B)
	return n / m
}

func (l *Line) GetLineNearestPoint(point *Vector) *Vector {
	direction := NewVector(l.B, -l.A, 0)
	u := point.Sub(l.Start).Dot(direction) / direction.SquareLength()
	return direction.MulConstant(u).Add(l.Start)
}

func (l *Line) GetSegmentNearestPoint(point *Vector) *Vector {
	direction := NewVector(l.B, -l.A, 0)
	u := point.Sub(l.Start).Dot(direction) / direction.SquareLength()
	if u < 0 {
		return l.Start
	}
	if u > 1 {
		return l.End
	}
	return direction.MulConstant(u).Add(l.Start)
}

func (l *Line) PointSide(point *Vector) int {
	s := l.A*(point.X-l.Start.X) + l.B*(point.Y-l.Start.Y)
	if s > 0 {
		return 1
	}
	if s < 0 {
		return -1
	}
	return 0
}

func (l *Line) CrossLineSegment(line *Line) (*Vector, bool) {
	d := Det(l.A, l.B, line.A, line.B)
	result := NewVector(-(Det(l.C, l.B, line.C, line.B) / d), -(Det(l.A, l.C, line.A, line.C) / d), 0)
	return result, Inside(result.X, math.Min(line.Start.X, line.End.X), math.Max(line.Start.X, line.End.X)) &&
		Inside(result.Y, math.Min(line.Start.Y, line.End.Y), math.Max(line.Start.Y, line.End.Y))
}

func (l *Line) CrossSegmentSegment(line *Line) (*Vector, bool) {
	d := Det(l.A, l.B, line.A, line.B)
	if d == 0 {
		return nil, false
	}
	result := NewVector(-(Det(l.C, l.B, line.C, line.B) / d), -(Det(l.A, l.C, line.A, line.C) / d), 0)
	return result, Inside(result.X, math.Min(l.Start.X, l.End.X), math.Max(l.Start.X, l.End.X)) &&
		Inside(result.Y, math.Min(l.Start.Y, l.End.Y), math.Max(l.Start.Y, l.End.Y)) &&
		Inside(result.X, math.Min(line.Start.X, line.End.X), math.Max(line.Start.X, line.End.X)) &&
		Inside(result.Y, math.Min(line.Start.Y, line.End.Y), math.Max(line.Start.Y, line.End.Y))
}

func (l *Line) CrossLineLine(line *Line) (*Vector, bool) {
	d := Det(l.A, l.B, line.A, line.B)
	if d == 0 {
		return nil, false
	}
	result := NewVector(-(Det(l.C, l.B, line.C, line.B) / d), -(Det(l.A, l.C, line.A, line.C) / d), 0)
	return result, true
}

func (l *Line) Length() float64 {
	x := l.End.X - l.Start.X
	y := l.End.Y - l.Start.Y
	return math.Sqrt(x*x + y*y)
}

func (l *Line) SquareLength() float64 {
	x := l.End.X - l.Start.X
	y := l.End.Y - l.Start.Y
	return x*x + y*y
}

func (l *Line) Reverse() *Line {
	return NewLineByVector(l.End, l.Start)
}

func (l *Line) GetPointAlong(t float64) *Vector {
	return l.End.Sub(l.Start).Norm().MulConstant(t).Add(l.Start)
}
