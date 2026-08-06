package allocation

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	production "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

func (h *Handler) AllocateProductionUnit(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*AllocateProductionUnitCommand)
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

		farmID := vo.ID(orgId)

		cycle, err := productionProvider.GrowingCycles().GetByID(ctx, cmd.CycleID, farmID)
		if err != nil {
			return nil, err
		}
		if cycle == nil {
			return nil, ErrGrowingCycleNotFound
		}

		item := allocation.New(farmID, cmd.CycleID, cmd.ProductionUnitID, cmd.Area, cmd.StartedAt)

		if err := productionProvider.Allocation().Save(ctx, item); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(item)
		return response.Id(item.ID), nil
	})
}
