package converter

import (
	"context"
	sm "github.com/flopp/go-staticmaps"
	"image"
	"image/color"
	"io"
	"sat-api/geometry"
	image2 "sat-api/image"
	"sat-api/model"
	"sync"
)

type OpenStreet struct {
	Name                    string
	Signature               string
	OptimizedSizeToleration float64
	AreaWeight              float64
	DefaultHeight           int
	DefaultWidth            int
	MaxZoom                 int
	MaximalArea             float64
	GeoDataMultiplier       float64
	Cache                   sm.TileCache
	ParallelAtNumberOf      int
}

func (o OpenStreet) Convert(ctx context.Context, data model.GeoData, configuration Configurations) (Result, error) {
	o.handleConfiguration(&configuration)
	data.Sanitize()
	poly := geometry.NewPolygon(data.ToVector(o.GeoDataMultiplier))
	polyArea := poly.CountSquare()
	areas := make([]*sm.Area, 0)
	// if the area is bigger than the max at max zoom
	if configuration.CutForMaxZoom && polyArea >= o.MaximalArea {
		divider := polyArea / o.MaximalArea
		//dividerRounded := int(math.Round(divider))
		// TODO: debug
		poly1, poly2, _, ok := poly.Split(divider)
		poly1LatLng := poly1.ToLatLng(o.GeoDataMultiplier)
		poly2LatLng := poly2.ToLatLng(o.GeoDataMultiplier)

		if ok {
			area := sm.NewArea(data.ToLatLng(), color.Black, color.Transparent, o.AreaWeight)
			areas = append(areas, area)
			if len(poly1LatLng) >= 3 {
				area2 := sm.NewArea(poly1LatLng, color.White, color.Transparent, o.AreaWeight)
				areas = append(areas, area2)
			}
			if len(poly2LatLng) >= 3 {
				area3 := sm.NewArea(poly2LatLng, color.White, color.Transparent, o.AreaWeight)
				areas = append(areas, area3)
			}
		}
	} else {
		// if no cutting of the area is needed just add the overall area
		weight := o.AreaWeight
		if configuration.Clipping {
			weight = 0.0
		}
		area := sm.NewArea(data.ToLatLng(), color.Black, color.Transparent, weight)
		areas = append(areas, area)
	}

	return Result{
		ContentType: "image/png",
		Images:      o.convertToImagesResult(ctx, configuration, areas),
	}, nil
}

func (o OpenStreet) Information(ctx context.Context) Information {
	providers := sm.GetTileProviders()
	mapTypes := make([]string, 0, len(providers))
	for _, provider := range providers {
		mapTypes = append(mapTypes, provider.Name)
	}
	return Information{
		Name:           o.Name,
		Clipping:       true,
		MapTypes:       mapTypes,
		DefaultMapType: o.defaultTile().Name,
		DefaultHeight:  o.DefaultHeight,
		DefaultWidth:   o.DefaultWidth,
	}
}

func (o OpenStreet) convertToImagesResult(ctx context.Context, configuration Configurations, areas []*sm.Area) []Image {
	images := make([]Image, 0)
	if len(areas) >= o.ParallelAtNumberOf {
		wg := &sync.WaitGroup{}
		mutex := &sync.Mutex{}
		for _, area := range areas {
			select {
			case <-ctx.Done():
				return images
			default:
				wg.Add(1)
				go func(a *sm.Area) {
					defer wg.Done()
					img := o.convertToImageResult(configuration, a)
					mutex.Lock()
					images = append(images, img)
					mutex.Unlock()
				}(area)
			}
			wg.Wait()
		}
		return images
	}
	for _, area := range areas {
		select {
		case <-ctx.Done():
			return images
		default:
			img := o.convertToImageResult(configuration, area)
			images = append(images, img)
		}
	}
	return images
}

func (o OpenStreet) convertToImageResult(configuration Configurations, area *sm.Area) Image {
	var raw io.Reader
	img, err := o.generateImg(configuration, area)
	if err == nil {
		raw, err = image2.ToPNGReader(img)
	}
	return Image{
		Data:  raw,
		Error: err,
		Geo:   model.NewGeoDataFromLatLng(area.Positions),
	}
}

func (o OpenStreet) generateImg(configuration Configurations, area *sm.Area) (image.Image, error) {
	ctx := sm.NewContext()
	ctx.SetTileProvider(o.getTileProvider(configuration.MapType))
	ctx.SetSize(configuration.Width, configuration.Height)
	ctx.SetCache(o.Cache)
	ctx.AddObject(area)
	img, tr, err := ctx.RenderWithTransformer()
	if err != nil {
		return nil, err
	}
	data := model.NewGeoDataFromLatLng(area.Positions)
	if configuration.Clipping {
		img = o.clip(img, data, tr, configuration.OptimizedSize)
	} else if configuration.OptimizedSize {
		img = o.optimizeSize(img, data, tr)
	}
	return img, nil
}

func (o OpenStreet) clip(i image.Image, data model.GeoData, transformer *sm.Transformer, crop bool) image.Image {
	absolutPositions := o.absolutPositions(data, transformer)
	xMin, yMin, xMax, yMax := image2.CalculateOptimizedSize(absolutPositions, 0)
	points := make([]*geometry.Vector, 0)
	for _, pos := range absolutPositions {
		points = append(points, geometry.NewVector(pos.X, pos.Y, 0))
	}
	absolutPoly := geometry.NewPolygon(points)
	imagePoly := image2.NewImagePolygon(absolutPoly, image2.Bound{
		X: image2.Point{
			X: xMin,
			Y: xMax,
		},
		Y: image2.Point{
			X: yMin,
			Y: yMax,
		},
	})
	dst := imagePoly.Clip(i, crop)
	return dst
}

func (o OpenStreet) optimizeSize(i image.Image, data model.GeoData, transformer *sm.Transformer) image.Image {
	xMin, yMin, xMax, yMax := image2.CalculateOptimizedSize(o.absolutPositions(data, transformer), o.OptimizedSizeToleration)
	return image2.Crop(i, image.Rect(xMin, yMin, xMax, yMax))
}

func (o OpenStreet) absolutPositions(data model.GeoData, transformer *sm.Transformer) []model.Point {
	absolutPositions := make([]model.Point, 0)
	for _, position := range data.ToLatLng() {
		x, y := transformer.LatLngToXY(position)
		absolutPositions = append(absolutPositions, model.Point{X: x, Y: y})
	}
	return absolutPositions
}

func (o OpenStreet) handleConfiguration(configurations *Configurations) {
	if configurations.Width <= 0 {
		configurations.Width = o.DefaultWidth
	}
	if configurations.Height <= 0 {
		configurations.Height = o.DefaultHeight
	}
}

func (o OpenStreet) getTileProvider(name string) *sm.TileProvider {
	provider, ok := sm.GetTileProviders()[name]
	if ok {
		return provider
	}
	return o.defaultTile()
}

func (o OpenStreet) defaultTile() *sm.TileProvider {
	t := new(sm.TileProvider)
	t.Name = "arcgis-worldimagery"
	t.Attribution = o.Signature
	t.TileSize = 256
	t.URLPattern = "https://server.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile/%[2]d/%[4]d/%[3]d"
	t.Shards = []string{}
	return t
}
