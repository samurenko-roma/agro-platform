// Package swaggerdocs содержит "теневые" (никогда не вызываемые) функции
// с swag-аннотациями — источник для генерации docs/swagger/openapi.json.
//
// У каждой CQRS-команды/запроса — свой задокументированный путь вида
// /api/commands/{name} или /api/queries/{name}. Эти пути РЕАЛЬНЫ (см.
// CommandByNameEndpoint/QueryByNameEndpoint в internal/interfaces/http) —
// тело запроса в них это сразу "data", без конверта. Классический вызов
// через POST /api/commands с конвертом {"command"/"query", "data"}
// продолжает работать параллельно, это просто два способа обратиться к
// одному и тому же диспатчеру.
//
// ВАЖНО: @BasePath намеренно пустой (не "/api"), потому что /auth/* живут
// на корне (см. router.go), а не под /api — префикс "/api" прописан явно
// в @Router каждой команды/запроса.
package swaggerdocs

// @title Agro Platform API
// @version 1.0
// @description CQRS-style API. Все команды — POST /api/commands с телом
// @description {"command": "<name>", "data": {...}}; все запросы — POST
// @description /api/queries с телом {"query": "<name>", "data": {...}}.
// @description Каждая операция ниже задокументирована как отдельный путь
// @description для читаемости; "Try it out" отправляет уже переписанный,
// @description реальный запрос.
// @host api.lab.note
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT из POST /auth/login. Формат: "Bearer <token>"
func GeneralInfo() {}
