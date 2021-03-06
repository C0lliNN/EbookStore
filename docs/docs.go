// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2022-05-22 09:09:56.655435769 -0300 -03 m=+0.060537754
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://github.com/C0lliNN",
        "contact": {
            "name": "Raphael Collin",
            "email": "raphael_professional@yahoo.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/books": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Fetch Books",
                "parameters": [
                    {
                        "type": "string",
                        "name": "authorName",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "description",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "perPage",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "title",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/catalog.PaginatedBooksResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Create a new Book",
                "parameters": [
                    {
                        "maxLength": 100,
                        "type": "string",
                        "name": "authorName",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "description",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "name": "price",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "releaseDate",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "maxLength": 100,
                        "type": "string",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Book Poster",
                        "name": "poster",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Book Content in PDF",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/catalog.BookResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/books/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Fetch Book by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/catalog.BookResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Delete a Book",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Success"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Catalog"
                ],
                "summary": "Update the provided Book",
                "parameters": [
                    {
                        "description": "Book Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/catalog.UpdateBook"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Book ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Success"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/healthcheck": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "REST API Healtcheck",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login using email and password",
                "parameters": [
                    {
                        "description": "Login Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/auth.CredentialsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Shop"
                ],
                "summary": "Fetch Orders",
                "parameters": [
                    {
                        "type": "integer",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "perPage",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "status",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/shop.PaginatedOrdersResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Shop"
                ],
                "summary": "Create a new Order",
                "parameters": [
                    {
                        "description": "Order Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/shop.CreateOrder"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/shop.OrderResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/orders/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Shop"
                ],
                "summary": "Fetch Order by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "orderId ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/shop.OrderResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/orders/{id}/download": {
            "get": {
                "produces": [
                    "application/pdf"
                ],
                "tags": [
                    "Shop"
                ],
                "summary": "Download the book for the given Order",
                "parameters": [
                    {
                        "description": "Order Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/shop.CreateOrder"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success"
                    },
                    "402": {
                        "description": "Payment Required",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/password-reset": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Reset the password for the given email",
                "parameters": [
                    {
                        "description": "Register Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.PasswordResetRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "success"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Register Payload",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/auth.CredentialsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/stripe/webhook": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Shop"
                ],
                "summary": "Handle stripe webhooks",
                "responses": {
                    "200": {
                        "description": "Success"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.CredentialsResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "auth.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                }
            }
        },
        "auth.PasswordResetRequest": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        },
        "auth.RegisterRequest": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName",
                "password",
                "passwordConfirmation"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string",
                    "maxLength": 150
                },
                "lastName": {
                    "type": "string",
                    "maxLength": 150
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                },
                "passwordConfirmation": {
                    "type": "string"
                }
            }
        },
        "catalog.BookResponse": {
            "type": "object",
            "properties": {
                "authorName": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "posterImageLink": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "releaseDate": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "catalog.PaginatedBooksResponse": {
            "type": "object",
            "properties": {
                "currentPage": {
                    "type": "integer"
                },
                "perPage": {
                    "type": "integer"
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/catalog.BookResponse"
                    }
                },
                "totalItems": {
                    "type": "integer"
                },
                "totalPages": {
                    "type": "integer"
                }
            }
        },
        "catalog.UpdateBook": {
            "type": "object",
            "properties": {
                "authorName": {
                    "type": "string",
                    "maxLength": 100
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string",
                    "maxLength": 100
                }
            }
        },
        "server.ErrorResponse": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "shop.CreateOrder": {
            "type": "object",
            "required": [
                "bookId"
            ],
            "properties": {
                "bookId": {
                    "type": "string",
                    "maxLength": 36
                }
            }
        },
        "shop.OrderResponse": {
            "type": "object",
            "properties": {
                "bookId": {
                    "type": "string"
                },
                "clientSecret": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "paymentIntentId": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "shop.PaginatedOrdersResponse": {
            "type": "object",
            "properties": {
                "currentPage": {
                    "type": "integer"
                },
                "perPage": {
                    "type": "integer"
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/shop.OrderResponse"
                    }
                },
                "totalItems": {
                    "type": "integer"
                },
                "totalPages": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "http://localhost:8080",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "E-book Store",
	Description: "Endpoints available in the E-book store REST API.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
