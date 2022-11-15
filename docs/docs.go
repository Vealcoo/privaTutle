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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/media/image": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "UploadImage",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "file",
                        "description": "上傳圖片",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "有效時間",
                        "name": "expirationTime",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "瀏覽密碼",
                        "name": "password",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/media/video": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "UploadVideo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "file",
                        "description": "上傳影片",
                        "name": "video",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "有效時間",
                        "name": "expirationTime",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "瀏覽密碼",
                        "name": "password",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/media/{short}": {
            "get": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Media"
                ],
                "summary": "GetMedia",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "short",
                        "name": "short",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "password",
                        "name": "password",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/short": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Short"
                ],
                "summary": "Short",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/router.ShortInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/short/{short}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Short"
                ],
                "summary": "GetShort",
                "parameters": [
                    {
                        "type": "string",
                        "description": "short",
                        "name": "short",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/router.LoginInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register",
                "parameters": [
                    {
                        "description": "body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/router.RegisterInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/short/{page}/{limit}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "ShortList",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/user/short/{shortId}": {
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "DeleteShort",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "shortId",
                        "name": "shortId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "router.LoginInfo": {
            "type": "object",
            "required": [
                "userId",
                "userPassword"
            ],
            "properties": {
                "userId": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                },
                "userPassword": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                }
            }
        },
        "router.RegisterInfo": {
            "type": "object",
            "required": [
                "email",
                "sex",
                "userId",
                "userName",
                "userPassword"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "sex": {
                    "type": "string",
                    "enum": [
                        "female",
                        "male",
                        "none"
                    ]
                },
                "userId": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                },
                "userName": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 2
                },
                "userPassword": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 6
                }
            }
        },
        "router.ShortInfo": {
            "type": "object",
            "required": [
                "leadUrl"
            ],
            "properties": {
                "leadUrl": {
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
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Gin swagger",
	Description:      "Gin swagger",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
