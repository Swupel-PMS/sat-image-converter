package geometry

import "math"

func Inside(v, min, max float64) bool {
	return (min <= v+Eps) && (v <= max+Eps)
}

func Det(a, b, c, d float64) float64 {
	return (a * d) - (b * c)
}

func GetBisector(l1 *Line, l2 *Line) *Line {
	q1 := math.Sqrt(l1.A*l1.A + l1.B*l1.B)
	q2 := math.Sqrt(l2.A*l2.A + l2.B*l2.B)
	A := l1.A/q1 - l2.A/q2
	B := l1.B/q1 - l2.B/q2
	C := l1.C/q1 - l2.C/q2
	return NewLineByValues(A, B, C)
}

func GetTanAngle(l1 *Line, l2 *Line) float64 {
	return (l1.A*l2.B - l2.A*l1.B) / (l1.A*l2.A + l1.B*l2.B)
}
