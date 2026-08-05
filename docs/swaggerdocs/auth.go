package swaggerdocs

import (
	"github.com/samurenkoroma/agro-platform/internal/application/commands/account/auth"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

var (
	_ = auth.RegisterRequest{}

	_ = response.CommandResponse[any]{}
)

// Эти три — единственные РЕАЛЬНЫЕ пути в этой спеке (не /api/commands или
// /api/queries под капотом). Regex-перехватчик в docs/swagger/index.html
// их не трогает (не матчит /api/commands/ или /api/queries/), поэтому
// "Try it out" работает штатно, без переписывания.

// docRegister godoc
// @Summary Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterRequest true "Данные регистрации"
// @Success 201 {object} auth.RegisterResponse
// @Failure 400 {object} response.CommandResponse[any]
// @Router /auth/register [post]
func docRegister() {}

// docLogin godoc
// @Summary Вход, выдаёт access/refresh токены
// @Tags auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Email и пароль"
// @Success 200 {object} response.CommandResponse[any]
// @Failure 401 {object} response.CommandResponse[any]
// @Router /auth/login [post]
func docLogin() {}

// docLogout godoc
// @Summary Выход (текущая реализация — заглушка, blacklist токенов не ведётся)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body auth.LogoutRequest false "Refresh-токен"
// @Success 200 {object} response.CommandResponse[any]
// @Router /auth/logout [post]
func docLogout() {}
