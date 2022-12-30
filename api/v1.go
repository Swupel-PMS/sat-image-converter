package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sat-api/converter"
	"sat-api/model"
)

type Payload struct {
	Coordinates []model.GeoPoint `json:"coordinates"`
	converter.Configurations
}

type V1 struct {
	Converter converter.Converter
}

// Convert godoc
//
//	@Summary		Generate Image
//	@Description	Generates an image from geo data
//	@Tags			V1
//	@Accept			json
//	@Produce		png, application/zip
//	@Param			data body Payload true " "
//	@Success		200
//	@Failure		409
//	@Failure		500
//	@Router			/v1/sat [post]
func (v V1) Convert() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// parse payload
		var payload Payload
		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			writer.WriteHeader(http.StatusConflict)
			log.Println(err)
			return
		}
		if len(payload.Coordinates) < 3 {
			writer.WriteHeader(http.StatusConflict)
			return
		}
		result, err := v.Converter.Convert(request.Context(), payload.Coordinates, payload.Configurations)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = result.ParseResponse(request.Context(), writer)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

// Info godoc
//
//	@Summary		Information
//	@Description	Informations and Feature set of the converter
//	@Tags			V1
//	@Produce		json
//	@Success		200 {object} converter.Information
//	@Failure		500
//	@Router			/v1/info [get]
func (v V1) Info() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(writer).Encode(v.Converter.Information(request.Context()))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (v V1) Base() string {
	return "/v1"
}
