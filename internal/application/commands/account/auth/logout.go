package auth

import (
	"net/http"

	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Logout
// @Summary Выход (текущая реализация — заглушка, blacklist токенов не ведётся)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LogoutRequest false "Refresh-токен"
// @Success 200
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// В простой реализации просто возвращаем успех
	// Для полноценной реализации нужно добавить blacklist токенов в Redis
	response.Success(map[string]string{
		"message": "logout successful",
	}).WriteJSON(w, http.StatusOK)
}
