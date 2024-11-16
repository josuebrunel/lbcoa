// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "demo.com",
        "contact": {
            "name": "API Support",
            "url": "http://demo.com/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Returns a JSON response with FizzBuzz result",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "fizzbuzz"
                ],
                "summary": "FizzBuzz",
                "operationId": "fizzbuzz",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Int1",
                        "name": "int1",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Int2",
                        "name": "int2",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Str1",
                        "name": "str1",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Str2",
                        "name": "str2",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/fizzbuzz_pkg_apiresponse.ApiResponse-string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/fizzbuzz_pkg_apiresponse.ApiResponse-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/fizzbuzz_pkg_apiresponse.ApiResponse-string"
                        }
                    }
                }
            }
        },
        "/stat": {
            "get": {
                "description": "Returns a JSON response of the query string with the most hits",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "fizzbuzz"
                ],
                "summary": "Stat",
                "operationId": "fizzbuzz-stat",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/fizzbuzz_pkg_apiresponse.ApiResponse-string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/fizzbuzz_pkg_apiresponse.ApiResponse-string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "fizzbuzz_pkg_apiresponse.ApiResponse-string": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "FizzBuzz Swagger UI",
	Description:      "FizzBuzz",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
