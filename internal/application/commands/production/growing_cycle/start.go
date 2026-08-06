package growingcycle

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
	production "github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

// Start регистрирует новый цикл выращивания СРАЗУ с физическим размещением
// (аллокациями) и, опционально, посевами — "быстрый старт" в одно действие.
// Для регистрации цикла без немедленного размещения используйте Create.
func (h *Handler) Start(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*StartGrowingCycleCMD)
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

		cycle := growingcycle.New(
			farmID, cmd.CropID,
			cmd.VarietyID, cmd.ProtocolID,
			cmd.Name, cmd.Code, cmd.Method)

		cycle.ChangeState(cmd.Stage)
		cycle.Start()

		if err := productionProvider.GrowingCycles().Save(ctx, cycle); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(cycle)

		for _, a := range cmd.Allocations {
			alloc := allocation.New(farmID, cycle.ID, a.ProductionUnitID, a.Area, &a.StartedAt)
			if err := productionProvider.Allocation().Save(ctx, alloc); err != nil {
				return nil, err
			}
			exec.RegisterAggregate(alloc)
		}

		for _, pl := range cmd.Plantings {
			p := planting.New(farmID, cycle.ID, pl.PlantedAt, pl.Quantity)
			if err := productionProvider.Planting().Save(ctx, p); err != nil {
				return nil, err
			}
			exec.RegisterAggregate(p)
		}

		return response.Id(cycle.ID), nil
	})
}
