package productionunit

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	_ "github.com/samurenkoroma/agro-platform/internal/interfaces/http/response"
)

type handler struct {
	units Projection
}

func NewGetOne(units Projection) queries.Handler {
	return &handler{
		units: units,
	}
}

type GetOneQuery struct {
	Id string `json:"id" validate:"required"`
}

// Ask docGetProductionUnit godoc
// @Summary Получить узел вместе с поддеревом
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GetOneQuery true "ID узла"
// @Success 200 {object} response.SuccessResponse{data=DTO}
// @Router /api/queries/spatial.get_production_unit [post]
func (h *handler) Ask(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*GetOneQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}

	return h.units.Get(ctx, vo.ID(q.Id))
}
