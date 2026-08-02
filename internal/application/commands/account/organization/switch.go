package organization

import (
	"context"
	"errors"
	"fmt"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/queries/account/dto"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/account/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/jwt"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type SwitchOrganizationCmd struct {
	OrganizationID string `json:"organization_id"`
}

type SwitchOrganizationResult struct {
	TokenPair  *jwt.TokenPair           `json:"tokenPair"`
	CurrentOrg dto.UserOrganizationInfo `json:"currentOrg"`
}

func (h *OrganizationHandler) Switch(ctx context.Context, cmd any) (any, error) {
	c, ok := cmd.(*SwitchOrganizationCmd)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}

	// Получаем текущего пользователя из контекста
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, user.ErrUnauthorized
	}

	return h.uow.Execute(ctx, providers.NewAccountProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {

		authProvider, ok := provider.(domain.AccountProvider)
		if !ok {
			return nil, fmt.Errorf("expected FarmProvider, got %T", provider)
		}

		userRepo := authProvider.Users()
		membershipRepo := authProvider.Memberships()
		orgRepo := authProvider.Organizations()

		// Получаем пользователя
		user, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// Проверяем членство в организации
		membership, err := membershipRepo.FindByUserAndOrganization(ctx, userID, c.OrganizationID)
		if err != nil {
			return nil, errors.New("you don't have access to this organization")
		}

		if !membership.IsActive {
			return nil, errors.New("membership is not active")
		}

		// Получаем информацию об организации
		org, err := orgRepo.FindByID(ctx, c.OrganizationID)
		if err != nil {
			return nil, err
		}

		// Обновляем текущую организацию в профиле
		user.SetCurrentOrganization(org.ID)
		if err := userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(user)
		// Генерируем новые токены с новой организацией
		tokenPair, err := h.jwtService.GenerateTokenPair(
			user.ID,
			user.Username,
			user.Email,
			string(user.Role),
			c.OrganizationID,
			string(membership.Role),
		)
		if err != nil {
			return nil, err
		}

		return SwitchOrganizationResult{
			TokenPair: tokenPair,
			CurrentOrg: dto.UserOrganizationInfo{
				OrganizationID:   org.ID,
				OrganizationName: org.Name,
				Role:             string(membership.Role),
				RoleName:         membership.GetRoleName(),
			},
		}, nil
	})

}
