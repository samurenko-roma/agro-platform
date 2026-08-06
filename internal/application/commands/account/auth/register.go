package auth

import (
	"context"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	account "github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/account/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=32"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// RegisterResponse ответ на регистрацию
type RegisterResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Message  string `json:"message"`
}

// Register godoc
// @Summary Регистрация
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Данные регистрации"
// @Success 201 {object} response.SuccessResponse{data=RegisterResponse}
// @Failure 400 {object} response.ErrResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*RegisterRequest)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}

	return h.uow.Execute(ctx, providers.NewAccountProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {
		// Приводим провайдер к нужному типу
		authProvider, ok := provider.(domain.AccountProvider)
		if !ok {
			if !ok {
			}
			if !ok {
				return nil, repository.ErrInvalidProviderType
			}
			return nil, repository.ErrInvalidProviderType
		}
		userRepo := authProvider.Users()

		// Проверяем, не существует ли пользователь
		existing, _ := userRepo.FindByEmail(ctx, cmd.Email)
		if existing != nil {
			return nil, account.ErrUserAlreadyExists
		}

		existing, _ = userRepo.FindByUsername(ctx, cmd.Username)
		if existing != nil {
			return nil, account.ErrUserAlreadyExists
		}

		// Создаем пользователя
		user, err := account.NewUser(
			cmd.Email, cmd.Username, cmd.Password,
			cmd.FirstName, cmd.LastName, cmd.Phone,
		)
		if err != nil {
			return nil, err
		}

		// Сохраняем
		if err := userRepo.Save(ctx, user); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(user)
		return RegisterResponse{
			UserID:   user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     string(user.Role),
			Message:  "User registered successfully",
		}, nil
	})

}
