package planting

import (
	"context"
	"errors"
	"time"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	production "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

func (h *Handler) Change(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*ChangePlantingCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}

	orgId, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}

	return h.uow.Execute(ctx, providers.NewProductionProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {
		productionProvider, ok := provider.(production.ProductionProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		root, err := productionProvider.Planting().GetByID(ctx, cmd.ID, vo.ID(orgId))
		if err != nil {
			return nil, err
		}
		if root == nil {
			return nil, ErrPlantingNotFound
		}

		root.PlantedAt = cmd.PlantedAt
		root.Quantity = cmd.Quantity
		root.UpdatedAt = time.Now()

		if err := productionProvider.Planting().Save(ctx, root); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(root)
		return response.Id(root.ID), nil
	})
}
