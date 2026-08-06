package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/account/auth"
	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/jwt"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

// RouterConfig конфигурация роутера
type RouterConfig struct {
	CommandRouter commands.Router
	QueryRouter   queries.Router
	Uow           uow.UnitOfWork
	JWTService    *jwt.Service
	Logger        *slog.Logger
}

// NewRouter создает новый HTTP роутер
func NewRouter(cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	// ========== AUTH ЭНДПОИНТЫ (без CQRS) ==========
	authHandler := auth.NewAuthHandler(cfg.Uow, cfg.JWTService)

	mux.HandleFunc("POST /auth/register", func(w http.ResponseWriter, r *http.Request) {
		var req *auth.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteValidationError(w, "invalid request body")
			return
		}

		// Валидация
		if req.Email == "" {
			response.WriteValidationError(w, "email is required")
			return
		}
		if req.Username == "" {
			response.WriteValidationError(w, "username is required")
			return
		}
		if req.Password == "" || len(req.Password) < 6 {
			response.WriteValidationError(w, "password must be at least 6 characters")
			return
		}

		res, err := authHandler.Register(r.Context(), req)
		if err != nil {
			response.WriteValidationError(w, err.Error())
		}
		response.Success(res).WriteJSON(w, http.StatusOK)
	})
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	//mux.HandleFunc("GET /auth/me", authHandler.Me)

	authMiddleware := NewAuthMiddleware(cfg.JWTService)
	// Защищенные эндпоинты (требуют аутентификации)
	mux.Handle("POST /auth/logout", authMiddleware.Authenticate(
		http.HandlerFunc(authHandler.Logout),
	))
	//mux.Handle("GET /auth/me", authMiddleware.Authenticate(
	//	http.HandlerFunc(authHandler.Me),
	//))

	// ========== CQRS ЭНДПОИНТЫ ==========
	// Команды и запросы идут через единый endpoint с аутентификацией

	protectedMux := http.NewServeMux()
	//protectedMux.Handle("/commands/", CommandEndpoint(cfg.CommandRouter))
	//protectedMux.Handle("/queries/", QueryEndpoint(cfg.QueryRouter))
	protectedMux.Handle("POST /commands/{name}", CommandByNameEndpoint(cfg.CommandRouter))
	protectedMux.Handle("POST /queries/{name}", QueryByNameEndpoint(cfg.QueryRouter))

	// Применяем middleware для защиты
	protectedHandler := authMiddleware.Authenticate(protectedMux)

	// Монтируем защищенные эндпоинты
	mux.Handle("/api/", http.StripPrefix("/api", protectedHandler))

	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("docs/swagger"))))
	// Опционально: эндпоинт для health check (без аутентификации)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(map[string]string{"status": "ok"}).WriteJSON(w, http.StatusOK)
	})

	// Применяем глобальные middleware (логирование, CORS, recovery)
	return withGlobalMiddleware(mux, cfg.Logger)
}

// withGlobalMiddleware применяет глобальные middleware ко всем запросам
func withGlobalMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	// Цепочка middleware (порядок важен!)
	handler := loggingMiddleware(next, log)
	handler = corsMiddleware(handler)
	handler = recoveryMiddleware(handler)
	return handler
}

// corsMiddleware добавляет CORS заголовки
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Organization-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// recoveryMiddleware восстанавливается после паники
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				response.WriteInternalError(w, fmt.Sprintf("internal server error %s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
