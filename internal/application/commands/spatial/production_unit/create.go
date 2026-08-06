package productionunit

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	spatial "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	topology "github.com/samurenkoroma/agro-platform/internal/domain/spatial/service"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

var topologyRules topology.TopologyRules = topology.DefaultTopology{}

// Create production unit
// @Summary Создать узел
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCommand true "Параметры узла"
// @Success 200 {object} response.SuccessResponse{data=response.IdResponse}
// @Failure 400 {object} response.ErrResponse "VALIDATION_ERROR"
// @Router /api/commands/spatial.create_production_unit [post]
func (h *Handler) Create(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*CreateCommand)
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

		var parentCode string

		if cmd.ParentID != nil {
			parent, err := spatialProvider.Units().GetByID(ctx, *cmd.ParentID, vo.ID(orgId))
			if err != nil {
				return nil, err
			}
			if parent == nil {
				return nil, pu.ErrParentNotFound
			}
			if parent.ArchivedAt != nil {
				return nil, pu.ErrParentArchived
			}
			if !topologyRules.CanAttach(parent.Type, cmd.Type) {
				return nil, pu.ErrInvalidHierarchy
			}
			parentCode = parent.Code
		}

		seq, err := spatialProvider.Units().GetNextSequence(ctx, vo.ID(orgId), cmd.ParentID, cmd.Type)
		if err != nil {
			return nil, err
		}

		code := pu.BuildCode(parentCode, cmd.Type, seq)

		unit := pu.New(vo.ID(orgId), cmd.ParentID, cmd.Type, code, &cmd.Name, seq)
		unit.Properties.AddCapabilities(cmd.Capabilities)

		if cmd.Dimensions != nil {
			unit.AddDimensions(cmd.Dimensions)
		}
		if cmd.Geometry != nil {
			unit.SetGeometry(*cmd.Geometry)
		}
		if cmd.Area != nil {
			m2, err := cmd.Area.Unit.ToSquareMeters(cmd.Area.Value)
			if err != nil {
				return nil, err
			}
			unit.SetArea(m2)
		}

		if err := spatialProvider.Units().Save(ctx, unit); err != nil {
			return nil, err
		}

		exec.RegisterAggregate(unit)
		return response.Id(unit.ID), nil
	})
}
