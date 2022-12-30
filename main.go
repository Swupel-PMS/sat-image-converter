package main

import (
	"context"
	"fmt"
	sm "github.com/flopp/go-staticmaps"
	"image"
	"image/png"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"sat-api/api"
	"sat-api/converter"
	"sat-api/geometry"
	image2 "sat-api/image"
	"sat-api/model"
	"strings"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	/*
		ClipDir("./toconvert", [][]model.Point{{
			{27, 95},
			{100, 550},
			{300, 600}, {668, 155}}, {
			{2, 65},
			{300, 650},
			{668, 65}}, {
			{2, 65},
			{300, 650},
			{668, 65}}, {
			{2, 65},
			{300, 650},
			{668, 65}}, {
			{2, 65},
			{300, 650},
			{668, 65}}})

	*/
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()
	err := api.Start(ctx, api.V1{Converter: converter.OpenStreet{
		Name:                    "OpenStreet",
		Signature:               "",
		OptimizedSizeToleration: 10,
		AreaWeight:              1.5,
		DefaultHeight:           512,
		DefaultWidth:            512,
		MaxZoom:                 20,
		MaximalArea:             25340.20196778653, // e-10
		GeoDataMultiplier:       10000000,
		Cache:                   sm.NewTileCacheFromUserCache(0777),
		ParallelAtNumberOf:      3, // TODO: analyze this
	}})
	if err != nil {
		log.Println(err)
	}
}

func ClipDir(path string, points [][]model.Point) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	pathToConvert := make([]string, 0)
	if !info.IsDir() {
		return fmt.Errorf("no dir")
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range dir {
		if !entry.IsDir() {
			pathToConvert = append(pathToConvert, entry.Name())
		}
	}
	wg := &sync.WaitGroup{}
	for i := 0; i < len(pathToConvert); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			err := ClipFile(filepath.Join(path, pathToConvert[index]), points[index])
			if err != nil {
				log.Println(err)
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func ClipFile(path string, points []model.Point) error {
	imgFile, err := os.Open(path)
	defer imgFile.Close()
	if err != nil {
		return err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	xMin, yMin, xMax, yMax := image2.CalculateOptimizedSize(points, 0)
	vectors := make([]*geometry.Vector, 0)
	for _, point := range points {
		vectors = append(vectors, geometry.NewVector(point.X, point.Y, 0))
	}
	absolutPoly := geometry.NewPolygon(vectors)
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
	imgClipped := imagePoly.Clip(img, true)
	name := fmt.Sprintf("%s", filepath.Join(filepath.Dir(path),
		fmt.Sprintf("%s%s", fmt.Sprintf("%s-clipped", strings.TrimSuffix(filepath.Base(path),
			filepath.Ext(path))), ".png")))
	f, err := os.Create(name)
	defer f.Close()
	if err != nil {
		return err
	}
	return png.Encode(f, imgClipped)
}
