// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://swupel-pms.vercel.app/",
        "contact": {
            "name": "Francisco Cardoso",
            "url": "https://swupel-pms.vercel.app/",
            "email": "swupelpms@gmail.com"
        },
        "license": {
            "name": "Copyright 2022 swupelpms",
            "url": "mailto:swupelpms@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/info": {
            "get": {
                "description": "Informations and Feature set of the converter",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V1"
                ],
                "summary": "Information",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/converter.Information"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/v1/sat": {
            "post": {
                "description": "Generates an image from geo data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "image/png",
                    " application/zip"
                ],
                "tags": [
                    "V1"
                ],
                "summary": "Generate Image",
                "parameters": [
                    {
                        "description": " ",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.Payload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "409": {
                        "description": "Conflict"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "api.Payload": {
            "type": "object",
            "properties": {
                "clipping": {
                    "description": "Clips the image",
                    "type": "boolean"
                },
                "coordinates": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.GeoPoint"
                    }
                },
                "cut_for_max_zoom": {
                    "description": "Cuts object into smaller ones to archive max zoom if needed. Every feature is also applied on the sub images",
                    "type": "boolean"
                },
                "height": {
                    "description": "Initial image height",
                    "type": "integer"
                },
                "map_type": {
                    "description": "Map type for image generation",
                    "type": "string"
                },
                "optimized_size": {
                    "description": "Reduces the image dimension if possible.",
                    "type": "boolean"
                },
                "width": {
                    "description": "Initial image width",
                    "type": "integer"
                }
            }
        },
        "converter.Information": {
            "type": "object",
            "properties": {
                "clipping": {
                    "description": "Clipping feature enabled",
                    "type": "boolean"
                },
                "default_height": {
                    "description": "Default height",
                    "type": "integer"
                },
                "default_map_type": {
                    "description": "Default map type",
                    "type": "string"
                },
                "default_width": {
                    "description": "Default width",
                    "type": "integer"
                },
                "map_types": {
                    "description": "Available map types",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "description": "Module name",
                    "type": "string"
                }
            }
        },
        "model.GeoPoint": {
            "type": "object",
            "properties": {
                "lat": {
                    "description": "Latitude",
                    "type": "number"
                },
                "long": {
                    "description": "Longitude",
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:10001",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Sat Image API",
	Description:      "API for converting geo data to image",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}