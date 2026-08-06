package productionunit

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	spatial "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

// Update production unit
// @Summary Обновить схему узла
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateCommand true "ID узла и схема"
// @Success 200 {object} response.SuccessResponse{data=response.IdResponse}
// @Router /api/commands/spatial.update_production_unit [post]
func docUpdateProductionUnit() {}
func (h *Handler) Update(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*UpdateCommand)
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

		if cmd.Schema != nil {
			unit.UpdateSchema(cmd.Schema)
		}
		if err := spatialProvider.Units().Save(ctx, unit); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(unit)
		return response.Id(unit.ID), nil
	})
}
