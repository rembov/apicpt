// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/api/auth/login": {
            "post": {
                "description": "Выполняет вход в систему с указанием email и пароля",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Аутентификация"
                ],
                "summary": "Вход в систему",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email пользователя",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Пароль пользователя",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Токены доступа и обновления",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Неверный логин или пароль",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/refresh-token": {
            "post": {
                "description": "Обновляет токен доступа с помощью токена обновления",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Аутентификация"
                ],
                "summary": "Обновление токена",
                "parameters": [
                    {
                        "description": "Токен обновления",
                        "name": "refreshToken",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Новый токен доступа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Токен обновления недействителен",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/register": {
            "post": {
                "description": "Регистрирует нового пользователя с указанными email, паролем и ролью",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Аутентификация"
                ],
                "summary": "Регистрация нового пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email пользователя",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Пароль пользователя",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Роль пользователя (Author или Reader)",
                        "name": "role",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Токены доступа и обновления",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат email",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Email уже существует",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/posts": {
            "get": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "Возвращает список постов",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Получение постов"
                ],
                "summary": "Получение списка постов",
                "responses": {
                    "200": {
                        "description": "Список постов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "Создает новый пост с заголовком и содержимым",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Управление постами"
                ],
                "summary": "Создание поста",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Уникальный ключ",
                        "name": "idempotencyKey",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Заголовок поста",
                        "name": "title",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Содержимое поста",
                        "name": "content",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пост успешно создан",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Пользователь не является автором",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "Уникальный ключ уже использовался",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/posts/{postId}": {
            "post": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "Редактирует пост по указанному идентификатору",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Управление постами"
                ],
                "summary": "Редактирование поста",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор поста",
                        "name": "postId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Новый заголовок поста",
                        "name": "title",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Новое содержимое поста",
                        "name": "content",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пост успешно обновлен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Доступ запрещен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пост не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/posts/{postId}/images": {
            "post": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "Добавляет изображение к посту",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Управление постами"
                ],
                "summary": "Добавление картинки к посту",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор поста",
                        "name": "postId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Изображение",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Картинка добавлена к посту",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Доступ запрещен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пост не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/posts/{postId}/images/{imageId}": {
            "delete": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "Удаляет изображение из поста по идентификатору",
                "tags": [
                    "Управление постами"
                ],
                "summary": "Удаление картинки из поста",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор поста",
                        "name": "postId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Идентификатор картинки",
                        "name": "imageId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Картинка успешно удалена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "403": {
                        "description": "Доступ запрещён",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пост или картинка не найдены",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/posts/{postId}/status": {
            "patch": {
                "security": [
                    {
                        "bearerAuth": []
                    }
                ],
                "description": "Публикует пост по указанному идентификатору",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Управление постами"
                ],
                "summary": "Публикация поста",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор поста",
                        "name": "postId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Статус поста (Published)",
                        "name": "status",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пост успешно опубликован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверное значение статуса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Пост не найден",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
