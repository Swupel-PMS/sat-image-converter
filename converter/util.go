package converter

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"math"
	"sat-api/model"
)

func calculateOptimizedSize(points []model.Point, tolerance float64) (int, int, int, int) {
	xMin := math.MaxFloat64
	xMax := 0.0
	yMin := math.MaxFloat64
	yMax := 0.0
	for _, position := range points {
		xMin = math.Min(xMin, position.X)
		yMin = math.Min(yMin, position.Y)
		xMax = math.Max(xMax, position.X)
		yMax = math.Max(yMax, position.Y)
	}
	return int(math.Round(xMin - tolerance)),
		int(math.Round(yMin - tolerance)),
		int(math.Round(xMax + tolerance)),
		int(math.Round(yMax + tolerance))
}

func Crop(i image.Image, rectangle image.Rectangle) image.Image {
	newImg, ok := i.(*image.RGBA)
	if !ok {
		return newImg
	}
	return newImg.SubImage(rectangle)
}

func ImageToPNGReader(img image.Image) (io.Reader, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
