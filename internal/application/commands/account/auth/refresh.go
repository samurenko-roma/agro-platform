package auth

import (
	"encoding/json"
	"net/http"

	"github.com/samurenkoroma/agro-platform/internal/infrastructure/jwt"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

// RefreshRequest запрос на обновление токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshResponse ответ с новой парой токенов
type RefreshResponse struct {
	TokenPair *jwt.TokenPair `json:"tokenPair"`
}

// Refresh godoc
// @Summary Перевыпуск токена
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "refresh token"
// @Success 200 {object} response.SuccessResponse{data=RefreshResponse}
// @Failure 400 {object} response.ErrResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteValidationError(w, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.WriteValidationError(w, "refresh_token is required")
		return
	}

	tokenPair, err := h.jwtService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.WriteUnauthorized(w, "invalid or expired refresh token")
		return
	}

	response.Success(RefreshResponse{
		TokenPair: tokenPair,
	}).WriteJSON(w, http.StatusOK)
}
