package geometry

import (
	"github.com/golang/geo/s2"
	"math"
)

type Polygon struct {
	Poly Vectors
}

func NewEmptyPolygon() *Polygon {
	return &Polygon{Poly: make([]*Vector, 0)}
}

func NewPolygon(vectors Vectors) *Polygon {
	return &Polygon{Poly: vectors}
}

func (p *Polygon) countSquareSigned() float64 {
	pointsCount := p.Length()
	if pointsCount < 3 {
		return 0
	}
	result := 0.0
	result += p.Vector(0).X * (p.Vector(pointsCount-1).Y - p.Vector(1).Y)
	for i := 1; i < pointsCount-1; i++ {
		result += p.Vector(i).X * (p.Vector(i-1).Y - p.Vector(i+1).Y)
	}
	result += p.Vector(pointsCount-1).X * (p.Vector(pointsCount-2).Y - p.Vector(0).Y)
	return result / 2.0
}

func (p *Polygon) CountSquare() float64 {
	return math.Abs(p.countSquareSigned())
}

func (p *Polygon) IsClockwise() bool {
	sum := 0.0
	t := p.Length() - 1
	for i := 0; i < t; i++ {
		sum += (p.Vector(i+1).X - p.Vector(i).X) * (p.Vector(i+1).Y + p.Vector(i).Y)
	}
	sum += (p.Vector(0).X - p.Vector(t).X) * (p.Vector(0).Y + p.Vector(t).Y)
	return sum <= 0
}

func (p *Polygon) Clear() {
	p.Poly = make([]*Vector, 0)
}

func (p *Polygon) Empty() bool {
	return p.Length() == 0
}

func (p *Polygon) Split(square float64) (*Polygon, *Polygon, *Line, bool) {
	polygonSize := p.Length()

	polygon := p.Poly
	if !p.IsClockwise() {
		p.Poly.Reserve()
	}
	poly1 := NewEmptyPolygon()
	poly2 := NewEmptyPolygon()
	cutLine := NewEmptyLine()

	if p.CountSquare()-square <= Eps {
		return nil, nil, nil, false
	}

	minCutLineExits := false
	minSqLength := math.MaxFloat64

	for i := 0; i < polygonSize; i++ {
		for j := i + 1; j < polygonSize; j++ {
			p1, p2 := p.CreateSubPoly(i, j)
			l1 := NewLineByVector(polygon[i], polygon[i+1])
			next := j + 1
			if next >= polygonSize {
				next = 0
			}
			l2 := NewLineByVector(polygon[j], polygon[next])
			cut, ok := GetCut(l1, l2, square, p1, p2)
			if ok {
				sqLength := cut.SquareLength()
				if sqLength < minSqLength && p.IsSegmentInsidePoly(cut, i, j) {
					minSqLength = sqLength
					poly1 = p1
					poly2 = p2
					cutLine = cut
					minCutLineExits = true
				}
			}

		}
	}
	if minCutLineExits {
		poly1.Add(cutLine.Start)
		poly1.Add(cutLine.End)
		poly2.Add(cutLine.End)
		poly2.Add(cutLine.Start)
		return poly1, poly2, cutLine, true
	}
	poly1 = NewPolygon(polygon)
	return poly1, poly2, cutLine, false
}

func (p *Polygon) FindDistance(point *Vector) float64 {
	distance := math.MaxFloat64
	for i := 0; i < len(p.Poly)-1; i++ {
		line := NewLineByVector(p.Poly[i], p.Poly[i+1])
		pV := line.GetSegmentNearestPoint(point)
		l := pV.Sub(point).Length()
		if l < distance {
			distance = l
		}
	}
	line := NewLineByVector(p.Poly[len(p.Poly)-1], p.Poly[0])
	pV := line.GetSegmentNearestPoint(point)
	l := pV.Sub(point).Length()
	if l < distance {
		distance = l
	}
	return distance
}

func (p *Polygon) FindNearestPoint(point *Vector) *Vector {
	result := &Vector{}
	distance := math.MaxFloat64
	for i := 0; i < p.Length()-1; i++ {
		line := NewLineByVector(p.Poly.Vector(i), p.Poly.Vector(i+1))
		pV := line.GetSegmentNearestPoint(point)
		l := pV.Sub(point).Length()
		if l < distance {
			distance = l
			result = pV
		}
	}
	line := NewLineByVector(p.Vector(p.Length()-1), p.Poly.Vector(0))
	pV := line.GetSegmentNearestPoint(point)
	l := pV.Sub(point).Length()
	if l < distance {
		distance = l
		result = pV
	}
	return result
}

func (p *Polygon) CountCenter() *Vector {
	return p.centroid()
}

func (p *Polygon) SplitNearestEdge(point *Vector) {
	result := &Vector{}
	ri := -1
	distance := math.MaxFloat64
	for i := 0; i < p.Length()-1; i++ {
		line := NewLineByVector(p.Vector(i), p.Vector(i+1))
		pV := line.GetSegmentNearestPoint(point)
		l := pV.Sub(point).Length()
		if l < distance {
			distance = l
			ri = i
			result = pV
		}
	}
	line := NewLineByVector(p.Vector(p.Length()-1), p.Vector(0))
	pV := line.GetSegmentNearestPoint(point)
	l := pV.Sub(point).Length()
	if l < distance {
		distance = l
		ri = p.Length() - 1
		result = pV
	}
	if ri != -1 {
		p.Poly.Insert(result, ri+1)
	}
}

func (p *Polygon) IsPointInside(point *Vector) bool {
	return p.isPointInsidePoly(point)
}

func (p *Polygon) IsSegmentInsidePoly(l *Line, excludeLine1 int, excludeLine2 int) bool {
	pointsCount := p.Length()
	for i := 0; i < pointsCount; i++ {
		if i != excludeLine1 && i != excludeLine2 {
			p1 := p.Vector(i)
			next := i + 1
			if next >= pointsCount {
				next = 0
			}
			p2 := p.Vector(next)
			nL := NewLineByVector(p1, p2)
			pV, ok := nL.CrossSegmentSegment(nL)
			if ok {
				if p1.Sub(pV).SquareLength() > Eps {
					if p2.Sub(pV).SquareLength() > Eps {
						return false
					}
				}
			}
		}
	}
	return p.isPointInsidePoly(l.GetPointAlong(0.5))
}

func (p *Polygon) CreateSubPoly(line1 int, line2 int) (*Polygon, *Polygon) {
	poly1 := NewEmptyPolygon()
	poly2 := NewEmptyPolygon()
	pc1 := line2 - line1
	for i := 1; i < pc1; i++ {
		poly1.Add(p.Vector(i + line1))
	}
	polySize := p.Length()
	pc2 := polySize - pc1
	for i := 1; i < pc2; i++ {
		poly2.Add(p.Vector((i + line2) % polySize))
	}
	return poly1, poly2
}

func (p *Polygon) centroid() *Vector {
	n := p.Length()
	result := Vector{}
	for i := 0; i < n; i++ {
		result.Add(p.Vector(i))
	}
	return result.DivConstant(float64(n))
}

func (p *Polygon) isPointInsidePoly(point *Vector) bool {
	pointsCount := p.Length() - 1
	l := NewDirectedLine(point, NewVector(0.0, 1e100, 0.0))
	result := 0
	for i := 0; i < pointsCount; i++ {
		line := NewLineByVector(p.Vector(i), p.Vector(i+1))
		_, ok := l.CrossSegmentSegment(line)
		if ok {
			result++
		}
	}
	line := NewLineByVector(p.Vector(pointsCount), p.Vector(0))
	_, ok := l.CrossSegmentSegment(line)
	if ok {
		result++
	}
	return result%2 != 0
}

func (p *Polygon) Add(vector *Vector) {
	p.Poly.Add(vector)
}

func (p *Polygon) ToLatLng(divider float64) []s2.LatLng {
	latLngS := make([]s2.LatLng, 0)
	for _, vector := range p.Poly {
		latLngS = append(latLngS, s2.LatLngFromPoint(s2.Point{Vector: vector.DivConstant(divider).Vector}))
	}
	return latLngS
}

func (p *Polygon) ToPoints() []s2.Point {
	points := make([]s2.Point, 0)
	for _, vector := range p.Poly {
		points = append(points, s2.Point{Vector: vector.Vector})
	}
	return points
}

func (p *Polygon) Vector(i int) *Vector {
	return p.Poly.Vector(i)
}

func (p *Polygon) Length() int {
	return len(p.Poly)
}
