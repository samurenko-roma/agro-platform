package warehouse

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	warehousedomain "github.com/samurenkoroma/agro-platform/internal/domain/inventory/aggregate/warehouse"
	invrepo "github.com/samurenkoroma/agro-platform/internal/domain/inventory/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type CreateWarehouseCommand struct {
	Name string  `json:"name" validate:"required"`
	Code *string `json:"code"`
}

func (h *Handler) Create(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*CreateWarehouseCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}
	return h.uow.Execute(ctx, providers.NewInventoryProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		inv, ok := p.(invrepo.InventoryProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}
		w := warehousedomain.New(vo.ID(orgID), cmd.Name)
		w.Code = cmd.Code
		if err := inv.Warehouses().Save(ctx, w); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(w)
		return response.Id(w.ID), nil
	})
}
