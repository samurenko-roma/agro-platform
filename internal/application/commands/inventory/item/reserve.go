package item

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/inventory/aggregate/movement"
	invrepo "github.com/samurenkoroma/agro-platform/internal/domain/inventory/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type ReserveCommand struct {
	ItemID        string  `json:"itemId"     validate:"required"`
	Amount        float64 `json:"amount"     validate:"required,gt=0"`
	ReferenceType *string `json:"referenceType"`
	ReferenceID   *string `json:"referenceId"`
}

func (h *Handler) Reserve(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*ReserveCommand)
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
		item, err := inv.Items().GetByID(ctx, vo.ID(cmd.ItemID))
		if err != nil {
			return nil, err
		}
		if err := item.Reserve(cmd.Amount); err != nil {
			return nil, err
		}
		if err := inv.Items().Save(ctx, item); err != nil {
			return nil, err
		}
		m := movement.New(vo.ID(orgID), item.ID, movement.Reserve, cmd.Amount)
		if cmd.ReferenceType != nil && cmd.ReferenceID != nil {
			m.Reference = &movement.Reference{
				Type: movement.ReferenceType(*cmd.ReferenceType),
				ID:   *cmd.ReferenceID,
			}
		}
		if err := inv.Movements().Save(ctx, m); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(item)
		exec.RegisterAggregate(m)
		return response.Id(item.ID), nil
	})
}
