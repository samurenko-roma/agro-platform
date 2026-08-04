// Package swaggerdocs содержит "теневые" (никогда не вызываемые) функции
// с swag-аннотациями — источник для генерации docs/swagger/openapi.json.
//
// У каждой CQRS-команды/запроса — свой задокументированный путь вида
// /api/commands/{name} или /api/queries/{name}. Эти пути РЕАЛЬНЫ (см.
// CommandByNameEndpoint/QueryByNameEndpoint в internal/interfaces/http) —
// тело запроса в них это сразу "data", без конверта. Классический вызов
// через POST /api/commands с конвертом {"command"/"query", "data"}
// продолжает работать параллельно — это просто два способа обратиться к
// одному и тому же диспатчеру.
//
// ВАЖНО: @BasePath намеренно пустой (не "/api"), потому что /auth/* живут
// на корне (см. router.go), а не под /api — префикс "/api" прописан явно
// в @Router каждой команды/запроса.
package swaggerdocs

// @title Agro Platform API
// @version 1.0
// @description CQRS-style API. У каждой команды/запроса — свой путь вида
// @description /api/commands/{name} или /api/queries/{name}, тело =
// @description data. Тот же результат — POST /api/commands или
// @description /api/queries с конвертом {"command"/"query": "<name>",
// @description "data": {...}}.
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT из POST /auth/login. Формат: "Bearer <token>"
func GeneralInfo() {}
