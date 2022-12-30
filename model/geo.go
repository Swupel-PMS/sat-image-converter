package model

import (
	"github.com/golang/geo/s2"
	"sat-api/geometry"
)

type Point struct {
	X float64
	Y float64
}

type GeoPoint struct {
	Lat  float64 `json:"lat"`  // Latitude
	Long float64 `json:"long"` // Longitude
}

func (g GeoPoint) Equal(data GeoPoint) bool {
	return g.Lat == data.Lat && g.Long == data.Long
}

func (g GeoPoint) ToLatLng() s2.LatLng {
	return s2.LatLngFromDegrees(g.Lat, g.Long)
}

func (g GeoPoint) ToVector(multiplier float64) *geometry.Vector {
	p := s2.PointFromLatLng(g.ToLatLng())
	return geometry.NewVectorByGeo(p.Mul(multiplier))
}

type GeoData []GeoPoint

func NewGeoDataFromLatLng(data []s2.LatLng) GeoData {
	geoPoints := make([]GeoPoint, 0)
	for _, d := range data {
		geoPoints = append(geoPoints, GeoPoint{
			Lat:  d.Lat.Degrees(),
			Long: d.Lng.Degrees(),
		})
	}
	return geoPoints
}

func (g *GeoData) ToLatLng() []s2.LatLng {
	ret := make([]s2.LatLng, 0, len(*g))
	for _, point := range *g {
		ret = append(ret, point.ToLatLng())
	}
	return ret
}

func (g *GeoData) ToVector(multiplier float64) []*geometry.Vector {
	ret := make([]*geometry.Vector, 0, len(*g))
	for _, point := range *g {
		ret = append(ret, point.ToVector(multiplier))
	}
	return ret
}

func (g *GeoData) Sanitize() {
	if (*g)[0].Equal((*g)[len(*g)-1]) {
		*g = (*g)[:len(*g)-1]
	}
}
