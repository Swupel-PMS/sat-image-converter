package main

import (
	"context"
	sm "github.com/flopp/go-staticmaps"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sat-api/api"
	"sat-api/converter"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	//debug.SetMaxStack(math.MaxInt)
	/*
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
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
