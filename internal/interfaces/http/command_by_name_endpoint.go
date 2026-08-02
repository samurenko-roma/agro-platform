package http

import (
	"io"
	"net/http"

	"github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/pkg/logger"
)

// CommandByNameEndpoint — тот же диспатч, что и CommandEndpoint, но имя
// команды берётся из пути (/api/commands/{name}), а тело запроса — сразу
// "data", без обёртки {"command": ..., "data": ...}. Существует ради
// удобного тестирования через Swagger/Scalar UI (у каждой команды свой
// путь и своя схема тела) и ради удобства ручных curl-вызовов. Обычный
// POST /api/commands с конвертом продолжает работать как прежде — этот
// route не заменяет его, а дополняет.
func CommandByNameEndpoint(router commands.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		name := r.PathValue("name")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.WriteValidationError(w, "invalid request body")
			return
		}

		log = log.With("command", name)

		handlerCmd, err := router.ResolveCommandPayload(name, body)
		if err != nil {
			log.Warn("command: failed to resolve", "error", err)
			response.WriteError(w, http.StatusBadRequest, response.CodeBadRequest,
				"failed to decode command: "+err.Error())
			return
		}

		result, err := router.Dispatch(r.Context(), name, handlerCmd)
		if err != nil {
			log.Warn("command: dispatch failed", "error", err)
			resp := response.FromError(err)
			resp.WriteJSON(w, getStatusCodeForError(resp.Error.Code))
			return
		}

		log.Info("command: executed")
		if result != nil {
			response.WriteSuccess(w, result)
		} else {
			response.WriteSuccess(w, map[string]string{"status": "ok"})
		}
	}
}
