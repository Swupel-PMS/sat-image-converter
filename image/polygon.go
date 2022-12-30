package image

import (
	"image"
	"image/color"
	"image/draw"
	"sat-api/geometry"
)

type Point struct {
	X int
	Y int
}

type Bound struct {
	X Point
	Y Point
}

type Polygon struct {
	*geometry.Polygon
	Bound Bound
}

func NewImagePolygon(polygon *geometry.Polygon, bound Bound) *Polygon {
	return &Polygon{
		Polygon: polygon,
		Bound:   bound,
	}
}

func (c *Polygon) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Polygon) Bounds() image.Rectangle {
	return image.Rect(c.Bound.X.X, c.Bound.Y.X, c.Bound.X.Y, c.Bound.Y.Y)
}

func (c *Polygon) At(x, y int) color.Color {
	if c.IsPointInside(geometry.NewVector(float64(x), float64(y), 0)) {
		return color.Alpha{A: 255}
	}
	return color.Alpha{A: 0}
}

func (c *Polygon) Clip(from image.Image, crop bool) image.Image {
	dst := image.NewRGBA(from.Bounds())
	draw.DrawMask(dst, dst.Bounds(), from, image.Point{}, c, image.Point{}, draw.Over)
	// TODO: refactor draw mask, to crop automatically
	if crop {
		return Crop(dst, image.Rect(c.Bound.X.X, c.Bound.Y.X, c.Bound.X.Y, c.Bound.Y.Y))
	}
	return dst
}
