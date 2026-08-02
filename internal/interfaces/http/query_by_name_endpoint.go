package http

import (
	"io"
	"net/http"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/pkg/logger"
)

// QueryByNameEndpoint — аналог CommandByNameEndpoint для запросов.
func QueryByNameEndpoint(router queries.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.FromContext(r.Context())
		name := r.PathValue("name")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.WriteValidationError(w, "invalid request body")
			return
		}

		log = log.With("query", name)

		result, err := router.Dispatch(r.Context(), name, body)
		if err != nil {
			log.Warn("query: dispatch failed", "error", err)
			resp := response.FromError(err)
			resp.WriteJSON(w, getStatusCodeForError(resp.Error.Code))
			return
		}

		log.Info("query: executed")
		response.WriteSuccess(w, result)
	}
}
