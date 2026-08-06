package item

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	inventoryitem "github.com/samurenkoroma/agro-platform/internal/domain/inventory/aggregate/inventory_item"
	invrepo "github.com/samurenkoroma/agro-platform/internal/domain/inventory/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type CreateItemCommand struct {
	Name        string  `json:"name"        validate:"required"`
	Type        string  `json:"type"        validate:"required"`
	Unit        string  `json:"unit"        validate:"required"`
	SKU         *string `json:"sku"`
	WarehouseID *string `json:"warehouseId"`
}

func (h *Handler) Create(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*CreateItemCommand)
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
		item := inventoryitem.New(vo.ID(orgID), cmd.Name, inventoryitem.Type(cmd.Type), inventoryitem.Unit(cmd.Unit))
		item.SKU = cmd.SKU
		if cmd.WarehouseID != nil {
			id := vo.ID(*cmd.WarehouseID)
			item.WarehouseID = &id
		}
		if err := inv.Items().Save(ctx, item); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(item)
		return response.Id(item.ID), nil
	})
}
