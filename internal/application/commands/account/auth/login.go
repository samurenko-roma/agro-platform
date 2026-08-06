package auth

import (
	"context"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	account "github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/account/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/jwt"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
type LoginResult struct {
	TokenPair    *jwt.TokenPair `json:"tokenPair"`
	User         User           `json:"user"`
	CurrentOrgId string         `json:"currentOrgId,omitempty"`
}

// Login
// @Summary Вход
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Email и пароль"
// @Success 200 {object} response.SuccessResponse{data=LoginResult}
// @Failure 400 {object} response.ErrResponse
// @Router /auth/login [post]
func docLogin() {}
func (h *AuthHandler) Login(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*LoginRequest)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	return h.uow.Execute(ctx, providers.NewAccountProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {
		// Приводим провайдер к нужному типу
		authProvider, ok := provider.(domain.AccountProvider)
		if !ok {
			if !ok {
			}
			return nil, repository.ErrInvalidProviderType
		}

		userRepo := authProvider.Users()
		membershipRepo := authProvider.Memberships()

		// Ищем пользователя
		user, err := userRepo.FindByEmail(ctx, cmd.Email)
		if err != nil {
			return nil, account.ErrInvalidCredentials
		}

		// Проверяем пароль
		if !user.CheckPassword(cmd.Password) {
			return nil, account.ErrInvalidCredentials
		}

		// Проверяем статус
		if !user.IsActive() {
			return nil, account.ErrUserInactive
		}

		currentOrgID := user.GetCurrentOrganizationID()
		var orgRole string
		if currentOrgID != "" {
			// Получаем все членства пользователя
			membership, err := membershipRepo.FindByUserAndOrganization(ctx, user.ID, currentOrgID)
			if err != nil {
				return nil, err
			}
			orgRole = membership.GetRoleName()
		}

		// Генерируем токены с текущей организацией

		tokenPair, err := h.jwtService.GenerateTokenPair(
			user.ID,
			user.Username,
			user.Email,
			string(user.Role),
			currentOrgID,
			orgRole,
		)
		if err != nil {
			return nil, err
		}

		// Обновляем время последнего входа
		user.UpdateLastLogin()
		userRepo.Update(ctx, user)

		exec.RegisterAggregate(user)
		return LoginResult{
			TokenPair: tokenPair,
			User: User{
				Id:    user.ID,
				Name:  user.Username,
				Email: user.Email,
				Role:  user.Role.String(),
			},
			CurrentOrgId: currentOrgID,
		}, nil

	})
}
