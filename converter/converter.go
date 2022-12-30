package converter

import (
	"context"
	"io"
	"sat-api/model"
)

type Converter interface {
	Convert(ctx context.Context, geo model.GeoData, configuration Configurations) (Result, error)
	Information(ctx context.Context) Information
}

type Image struct {
	Data  io.Reader     `json:"-"`   // raw image data
	Error error         `json:"-"`   // error happened during image conversion
	Geo   model.GeoData `json:"geo"` // geo data for the image
}

type Configurations struct {
	Height        int    `json:"height"`           // Initial image height
	Width         int    `json:"width"`            // Initial image width
	OptimizedSize bool   `json:"optimized_size"`   // Reduces the image dimension if possible.
	CutForMaxZoom bool   `json:"cut_for_max_zoom"` // Cuts object into smaller ones to archive max zoom if needed. Every feature is also applied on the sub images
	Clipping      bool   `json:"clipping"`         // Clips the image
	MapType       string `json:"map_type"`         // Map type for image generation
}

type Information struct {
	Name           string   `json:"name"`             // Module name
	Clipping       bool     `json:"clipping"`         // Clipping feature enabled
	MapTypes       []string `json:"map_types"`        // Available map types
	DefaultMapType string   `json:"default_map_type"` // Default map type
	DefaultHeight  int      `json:"default_height"`   // Default height
	DefaultWidth   int      `json:"default_width"`    // Default width
}
