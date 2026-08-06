package organization

import (
	"context"
	"fmt"

	"github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/queries/account/dto"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/organization"
	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/account/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type CreateOrganizationCmd struct {
	Name string `json:"name"`
}

// Create
// @Summary Создать организацию
// @Description Создать организацию (ферму/хозяйство), пользователь становится владельцем
// @Tags account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrganizationCmd true "Название организации"
// @Success 200 {object} response.SuccessResponse{data=dto.UserOrganizationInfo}
// @Failure 400 {object} response.ErrResponse "VALIDATION_ERROR"
// @Router /api/commands/account.create_organization [post]
func (h *OrganizationHandler) Create(ctx context.Context, cmd any) (any, error) {
	c, ok := cmd.(*CreateOrganizationCmd)
	if !ok {
		return nil, commands.ErrInvalidCommandType
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
		_, err := userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		newOrg, err := organization.NewSimpleOrganization(c.Name)
		if err != nil {
			return nil, err
		}

		if err := orgRepo.Save(ctx, newOrg); err != nil {
			return nil, err
		}

		membership, err := organization.NewMembership(userID, newOrg.ID, organization.OrgRoleOwner)
		if err != nil {
			return nil, err
		}

		if err := membershipRepo.Save(ctx, membership); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(newOrg)
		return dto.UserOrganizationInfo{
			OrganizationID:   newOrg.ID,
			OrganizationName: newOrg.Name,
			Role:             string(membership.Role),
			RoleName:         membership.GetRoleName(),
		}, nil
	})
}
