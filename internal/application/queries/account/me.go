package account

import (
	"context"
	"fmt"

	"github.com/samurenkoroma/agro-platform/internal/application/commands/account/auth"
	"github.com/samurenkoroma/agro-platform/internal/application/queries/account/dto"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/account/aggregate/user"
	domain "github.com/samurenkoroma/agro-platform/internal/domain/account/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type MeQuery struct {
}

// MeResponse ответ с информацией о текущем пользователе
type MeResponse struct {
	User         auth.User                   `json:"user"`
	Organization []*dto.UserOrganizationInfo `json:"organizations"`
	CurrentOrgId string                      `json:"currentOrgId"`
}

// Ask
// @Summary Текущий пользователь
// @Description Текущий пользователь + список организаций, в которых он состоит
// @Tags account
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=MeResponse}
// @Router /api/queries/account.me [post]
func (h *UserHandler) Ask(ctx context.Context, cmd any) (any, error) {
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

		orgs, err := orgRepo.ListByUser(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		var organizations []*dto.UserOrganizationInfo
		for _, o := range orgs {
			member, err2 := membershipRepo.FindByUserAndOrganization(ctx, user.ID, o.ID)
			if err2 != nil {
				return nil, err2
			}
			organizations = append(organizations, &dto.UserOrganizationInfo{
				OrganizationID:   o.ID,
				OrganizationName: o.Name,
				Role:             member.GetRoleName(),
			})
		}
		return MeResponse{
			User: auth.User{
				Id:    user.ID,
				Name:  user.Username,
				Email: user.Email,
				Role:  user.Role.String(),
			},
			Organization: organizations,
			CurrentOrgId: user.GetCurrentOrganizationID(),
		}, nil

	})
}
