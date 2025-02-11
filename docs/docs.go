// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "@lettons",
            "url": "https://t.me/lettons",
            "email": "kk7453603@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth": {
            "post": {
                "description": "Аутентифицирует пользователя и возвращает JWT-токен. При первой аутентификации пользователь создается автоматически.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Аутентификация и получение JWT-токена",
                "parameters": [
                    {
                        "description": "Данные для аутентификации",
                        "name": "authRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/buy/{item}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Покупает мерч для авторизованного пользователя. Имя товара передается в параметре пути.",
                "produces": [
                    "application/json"
                ],
                "summary": "Купить мерч за монетки",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Наименование мерча",
                        "name": "item",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Количество мерча",
                        "name": "count",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Покупка выполнена успешно",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/info": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает баланс монет, список купленных товаров и историю переводов (полученные и отправленные).",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение информации о монетах, инвентаре и истории транзакций",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.InfoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/sendCoin": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Переводит указанное количество монет от авторизованного пользователя к другому.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Отправить монеты другому пользователю",
                "parameters": [
                    {
                        "description": "Данные для перевода монет",
                        "name": "sendCoinRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SendCoinRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Перевод выполнен успешно",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AuthRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "description": "Пароль для аутентификации.",
                    "type": "string"
                },
                "username": {
                    "description": "Имя пользователя для аутентификации.",
                    "type": "string"
                }
            }
        },
        "models.AuthResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "description": "JWT-токен для доступа к защищённым ресурсам.",
                    "type": "string"
                }
            }
        },
        "models.CoinHistory": {
            "type": "object",
            "properties": {
                "received": {
                    "description": "Транзакции по полученным монетам.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.ReceivedTransaction"
                    }
                },
                "sent": {
                    "description": "Транзакции по отправленным монетам.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.SentTransaction"
                    }
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "description": "Сообщение об ошибке.",
                    "type": "string"
                }
            }
        },
        "models.InfoResponse": {
            "type": "object",
            "properties": {
                "coinHistory": {
                    "description": "История транзакций по монетам.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/models.CoinHistory"
                        }
                    ]
                },
                "coins": {
                    "description": "Количество доступных монет.",
                    "type": "integer"
                },
                "inventory": {
                    "description": "Список купленных товаров.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Item"
                    }
                }
            }
        },
        "models.Item": {
            "type": "object",
            "properties": {
                "quantity": {
                    "description": "Количество предметов.",
                    "type": "integer"
                },
                "type": {
                    "description": "Тип предмета.",
                    "type": "string"
                }
            }
        },
        "models.ReceivedTransaction": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "Количество полученных монет.",
                    "type": "integer"
                },
                "fromUser": {
                    "description": "Имя пользователя, отправившего монеты.",
                    "type": "string"
                }
            }
        },
        "models.SendCoinRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "Количество монет для перевода.",
                    "type": "integer"
                },
                "toUser": {
                    "description": "Имя пользователя, которому нужно отправить монеты.",
                    "type": "string"
                }
            }
        },
        "models.SentTransaction": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "Количество отправленных монет.",
                    "type": "integer"
                },
                "toUser": {
                    "description": "Имя пользователя, которому отправлены монеты.",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "0.0.0.0",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Merch Store API",
	Description:      "This is a merch store server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
