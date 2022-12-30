package geometry

import (
	"github.com/golang/geo/r3"
	"math"
)

var Eps = 1e-6

type Vector struct {
	r3.Vector
}

// NewVector - creates new vector from absolut points
func NewVector(x, y, z float64) *Vector {
	return &Vector{r3.Vector{
		X: x,
		Y: y,
		Z: z,
	}}
}

// NewVectorByGeo - creates a new vector by providing the underlying object
func NewVectorByGeo(vector r3.Vector) *Vector {
	return &Vector{vector}
}

func (v *Vector) Dot(vector *Vector) float64 {
	return v.Vector.Dot(vector.Vector)
}

func (v *Vector) SquareLength() float64 {
	return v.Dot(v)
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.SquareLength())
}

func (v *Vector) Sub(vector *Vector) *Vector {
	return NewVectorByGeo(v.Vector.Sub(vector.Vector))
}

func (v *Vector) Add(vector *Vector) *Vector {
	return NewVectorByGeo(v.Vector.Add(vector.Vector))
}

func (v *Vector) Mul(vector *Vector) *Vector {
	return NewVector(v.X*vector.X, v.Y*vector.Y, v.Z*vector.Z)
}

func (v *Vector) Div(vector *Vector) *Vector {
	return NewVector(v.X/vector.X, v.Y/vector.Y, v.Z/vector.Z)
}

func (v *Vector) MulConstant(value float64) *Vector {
	return NewVector(v.X*value, v.Y*value, v.Z*value)
}

func (v *Vector) DivConstant(value float64) *Vector {
	return NewVector(v.X/value, v.Y/value, v.Z/value)
}

func (v *Vector) Norm() *Vector {
	// TODO: maybe change to squareLength
	l := v.Length()
	if l == 0 {
		return &Vector{}
	}
	return v.DivConstant(l)
}

func (v *Vector) Equal(vector *Vector) bool {
	return v.X == vector.X && v.Y == vector.Y && v.Z == vector.Z
}
