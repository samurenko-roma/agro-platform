package operationevent

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	operationevent "github.com/samurenkoroma/agro-platform/internal/domain/operations/aggregate/operation_event"
	tl "github.com/samurenkoroma/agro-platform/internal/domain/operations/aggregate/timeline"
	opsrepo "github.com/samurenkoroma/agro-platform/internal/domain/operations/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type RecordOperationCommand struct {
	Type             string         `json:"type" validate:"required"`
	ProductionUnitID *string        `json:"productionUnitId"`
	GrowingCycleID   *string        `json:"growingCycleId"`
	Payload          map[string]any `json:"payload"`
}

func (h *Handler) Record(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*RecordOperationCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}
	userID, _ := ctx.Value("user_id").(string)

	return h.uow.Execute(ctx, providers.NewOperationsProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		ops, ok := p.(opsrepo.OperationsProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		farmID := vo.ID(orgID)
		e := operationevent.New(farmID, operationevent.OperationType(cmd.Type))

		if cmd.Payload != nil {
			for k, v := range cmd.Payload {
				e.Payload[k] = v
			}
		}
		if cmd.ProductionUnitID != nil {
			id := vo.ID(*cmd.ProductionUnitID)
			e.ProductionUnitID = &id
		}
		if cmd.GrowingCycleID != nil {
			id := vo.ID(*cmd.GrowingCycleID)
			e.GrowingCycleID = &id
		}
		if userID != "" {
			id := vo.ID(userID)
			e.PerformedBy = &id
		}

		if err := ops.Operations().Save(ctx, e); err != nil {
			return nil, err
		}

		// append to timeline
		var cycleID *vo.ID
		if e.GrowingCycleID != nil {
			cycleID = e.GrowingCycleID
		}
		timeline, err := ops.Timelines().GetByOwner(ctx, farmID, cycleID)
		if err != nil || timeline == nil {
			timeline = tl.New(farmID)
			timeline.GrowingCycleID = cycleID
		}
		item := tl.Item{
			ID:          vo.NewID(),
			Source:      tl.OperationSource,
			ReferenceID: e.ID,
			Title:       string(e.Type),
			Timestamp:   e.Timestamp,
			Metadata:    vo.NewMetadata(),
		}
		timeline.AddItem(item)
		if err := ops.Timelines().Save(ctx, timeline); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(timeline)

		return response.Id(e.ID), nil
	})
}
