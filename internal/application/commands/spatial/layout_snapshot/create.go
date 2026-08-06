package layoutsnapshot

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	ls "github.com/samurenkoroma/agro-platform/internal/domain/spatial/entity/layout_snapshot"
	pus "github.com/samurenkoroma/agro-platform/internal/domain/spatial/entity/production_unit_snapshot"
	spatial "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

func (h *Handler) Create(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*CreateSnapshotCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}

	orgId, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, errors.New("user_id is required")
	}

	return h.uow.Execute(ctx, providers.NewSpatialProvider, func(provider repository.RepositoryProvider, exec uow.Execution) (any, error) {
		spatialProvider, ok := provider.(spatial.SpatialProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		farmID := vo.ID(orgId)

		units, err := spatialProvider.Units().ListByOwner(ctx, farmID)
		if err != nil {
			return nil, err
		}

		latest, err := spatialProvider.Snapshots().GetLatest(ctx, farmID)
		if err != nil {
			return nil, err
		}
		nextVersion := 1
		if latest != nil {
			nextVersion = latest.Root.Version + 1
		}

		snapshot, err := ls.New(farmID, nextVersion, vo.ID(userId), cmd.Description)
		if err != nil {
			return nil, err
		}

		for _, unit := range units {
			name := ""
			var capabilities []pu.Capability
			var metadata vo.Metadata
			if unit.Properties != nil {
				if n, ok := unit.Properties.Metadata["name"].(string); ok {
					name = n
				}
				capabilities = unit.Properties.Capabilities
				metadata = unit.Properties.Metadata
			}

			if err := snapshot.AddUnit(pus.ProductionUnitSnapshot{
				ID:             vo.NewID(),
				SnapshotID:     snapshot.Root.ID,
				OriginalUnitID: unit.ID,
				Type:           unit.Type,
				Name:           name,
				ParentID:       unit.ParentID,
				Capabilities:   capabilities,
				Metadata:       metadata,
				CreatedAt:      unit.CreatedAt,
			}); err != nil {
				return nil, err
			}
		}

		if err := spatialProvider.Snapshots().Save(ctx, snapshot); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(snapshot)

		return response.Id(snapshot.Root.ID), nil
	})
}
