package productionunit

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	spatial2 "github.com/samurenkoroma/agro-platform/internal/application/services/spatial"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	pu "github.com/samurenkoroma/agro-platform/internal/domain/spatial/aggregate/production_unit"
	spatial "github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

// Configure production unit
// @Summary Сгенерировать схему узла
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ConfigureCommand true "ID родителя и схема раскладки"
// @Success 200 {object} response.SuccessResponse{data=response.IdResponse}
// @Failure 400 {object} response.ErrResponse "VALIDATION_ERROR"
// @Router /api/commands/spatial.configure_production_unit [post]
func (h *Handler) Configure(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*ConfigureCommand)
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

		parent, err := spatialProvider.Units().GetByID(ctx, cmd.Id, vo.ID(orgId))
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, ErrProductionUnitNotFound
		}
		if parent.ArchivedAt != nil {
			return nil, pu.ErrParentArchived
		}

		if err := spatial2.NewUnitLayoutGenerator(spatialProvider.Units(), exec).Generate(ctx, parent, cmd.Schema); err != nil {
			return nil, err
		}
		return response.Id(parent.ID), nil
	})
}
