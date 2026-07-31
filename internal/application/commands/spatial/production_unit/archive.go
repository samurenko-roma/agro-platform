package productionunit

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/response"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	spatial "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

func (h *Handler) Archive(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*ArchiveCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}

	orgId, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}

	return h.uow.Execute(ctx, providers.NewSpatialProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {
		spatialProvider, ok := provider.(spatial.SpatialProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		unit, err := spatialProvider.Units().GetByID(ctx, cmd.Id, vo.ID(orgId))
		if err != nil {
			return nil, err
		}
		if unit == nil {
			return nil, ErrProductionUnitNotFound
		}

		children, err := spatialProvider.Units().GetChildren(ctx, unit.ID)
		if err != nil {
			return nil, err
		}
		for _, c := range children {
			if c.ArchivedAt == nil {
				return nil, ErrHasActiveChildren
			}
		}

		if err := unit.Archive(); err != nil {
			return nil, err
		}

		if err := spatialProvider.Units().Save(ctx, unit); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(unit)
		return response.Id(unit.ID), nil
	})
}
