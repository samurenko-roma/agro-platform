package eventhandlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	"github.com/samurenkoroma/agro-platform/internal/domain/shared/event"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	spatialrepo "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	spatialProviders "github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/bus"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

func RegisterAllocationHandlers(eventBus bus.EventBus, unitOfWork uow.UnitOfWork) {
	eventBus.Register(allocation.EventAllocated, onAllocationAllocated(unitOfWork))
	eventBus.Register(allocation.EventReleased, onAllocationReleased(unitOfWork))
}

func onAllocationAllocated(unitOfWork uow.UnitOfWork) bus.EventHandler {
	return func(ctx context.Context, e event.DomainEvent) error {
		evt, ok := e.(allocation.AllocationAllocated)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", e)
		}
		orgId, ok := ctx.Value("organization_id").(string)
		if !ok {
			return errors.New("organization_id is required")
		}
		_, err := unitOfWork.Execute(ctx, spatialProviders.NewSpatialProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
			spatial, ok := p.(spatialrepo.SpatialProvider)
			if !ok {
				return nil, repository.ErrInvalidProviderType
			}
			unit, err := spatial.Units().GetByID(ctx, evt.ProductionUnitID, vo.ID(orgId))
			if err != nil {
				return nil, fmt.Errorf("production unit %s: %w", evt.ProductionUnitID, err)
			}
			if unit == nil {
				return nil, fmt.Errorf("production unit %s not found", evt.ProductionUnitID)
			}

			unit.Occupy()

			if err := spatial.Units().Save(ctx, unit); err != nil {
				return nil, err
			}

			exec.RegisterAggregate(unit)

			return nil, nil
		})
		return err
	}
}

func onAllocationReleased(unitOfWork uow.UnitOfWork) bus.EventHandler {
	return func(ctx context.Context, e event.DomainEvent) error {
		evt, ok := e.(allocation.AllocationReleased)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", e)
		}
		orgId, ok := ctx.Value("organization_id").(string)
		if !ok {
			return errors.New("organization_id is required")
		}

		_, err := unitOfWork.Execute(ctx, spatialProviders.NewSpatialProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
			spatial, ok := p.(spatialrepo.SpatialProvider)
			if !ok {
				return nil, repository.ErrInvalidProviderType
			}

			unit, err := spatial.Units().GetByID(ctx, evt.ProductionUnitID, vo.ID(orgId))
			if err != nil {
				return nil, fmt.Errorf("production unit %s: %w", evt.ProductionUnitID, err)
			}
			if unit == nil {
				return nil, fmt.Errorf("production unit %s not found", evt.ProductionUnitID)
			}

			unit.Release()

			if err := spatial.Units().Save(ctx, unit); err != nil {
				return nil, err
			}

			exec.RegisterAggregate(unit)
			return nil, nil
		})
		return err
	}
}
