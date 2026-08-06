package productionunit

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type listChildrenHandler struct {
	units Projection
}

func NewListChildren(units Projection) queries.Handler {
	return &listChildrenHandler{units: units}
}

type ListChildrenQuery struct {
	ParentID string `json:"parentId" validate:"required"`
}

// Ask Список детских узлов
// @Summary Список детских узлов
// @Tags spatial
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]productionunit.DTO}
// @Router /api/queries/spatial.get_production_unit_children [post]
func (h *listChildrenHandler) Ask(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*ListChildrenQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.units.Children(ctx, vo.ID(q.ParentID))
}
