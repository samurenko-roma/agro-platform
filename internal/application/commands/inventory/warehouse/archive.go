package warehouse

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	invrepo "github.com/samurenkoroma/agro-platform/internal/domain/inventory/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type ArchiveWarehouseCommand struct {
	WarehouseID string `json:"warehouseId" validate:"required"`
}

func (h *Handler) Archive(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*ArchiveWarehouseCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	if _, ok := ctx.Value("organization_id").(string); !ok {
		return nil, errors.New("organization_id is required")
	}
	return h.uow.Execute(ctx, providers.NewInventoryProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		inv, ok := p.(invrepo.InventoryProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}
		w, err := inv.Warehouses().GetByID(ctx, vo.ID(cmd.WarehouseID))
		if err != nil {
			return nil, err
		}
		w.Archive()
		if err := inv.Warehouses().Save(ctx, w); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(w)
		return response.Id(w.ID), nil
	})
}
