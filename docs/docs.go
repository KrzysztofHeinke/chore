// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/kv/{type}": {
            "get": {
                "description": "Get the value of the key",
                "summary": "Get Specific key",
                "parameters": [
                    {
                        "enum": [
                            "templates",
                            "auths",
                            "binds"
                        ],
                        "type": "string",
                        "description": "type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "key of the file",
                        "name": "key",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Set a key with value",
                "consumes": [
                    "text/plain"
                ],
                "summary": "New key or replace key",
                "parameters": [
                    {
                        "enum": [
                            "templates",
                            "auths",
                            "binds"
                        ],
                        "type": "string",
                        "description": "type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "key of the file",
                        "name": "key",
                        "in": "query"
                    },
                    {
                        "description": "any value to store",
                        "name": "payload",
                        "in": "body",
                        "schema": {
                            "type": "primitive"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Set a key with value",
                "consumes": [
                    "text/plain"
                ],
                "summary": "New key",
                "parameters": [
                    {
                        "enum": [
                            "templates",
                            "auths",
                            "binds"
                        ],
                        "type": "string",
                        "description": "type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "key of the file",
                        "name": "key",
                        "in": "query"
                    },
                    {
                        "description": "any value to store",
                        "name": "payload",
                        "in": "body",
                        "schema": {
                            "type": "primitive"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete key",
                "summary": "Delete key",
                "parameters": [
                    {
                        "enum": [
                            "templates",
                            "auths",
                            "binds"
                        ],
                        "type": "string",
                        "description": "type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "key of the file",
                        "name": "key",
                        "in": "query"
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/kv/{type}/": {
            "get": {
                "description": "Get list of keys",
                "summary": "Get List",
                "parameters": [
                    {
                        "enum": [
                            "templates",
                            "auths",
                            "binds"
                        ],
                        "type": "string",
                        "description": "type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "key of the file",
                        "name": "key",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "default": "true",
                        "description": "is it for listing",
                        "name": "list",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "folder/ or file",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "418": {
                        "description": "not a list",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/send": {
            "get": {
                "description": "Send request to api.",
                "consumes": [
                    "*/*"
                ],
                "summary": "Send request",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
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
	Version:     "",
	Host:        "",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "",
	Description: "storage and send API",
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
	swag.Register(swag.Name, &s{})
}
