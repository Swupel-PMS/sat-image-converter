# GeoData to Image API

This is an api for converting geo data to an image. 
It supports different map types and dimensions, clipping and other image manipulations.
Check out the OpenAPI doc for the full feature set.

## Run

- Install Go [https://go.dev/dl/](https://go.dev/dl/)
- Run
````shell
go run main.go
````

## Environment variables



## Documentation

Run th api and open ``localhost:10001/swagger/index.html``

## TODO

### Must

- poly split - refactor and debug with test data

### Improvements

- change compression algo for delivering data
- refactor static map to use faster decompression