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

func (h *listChildrenHandler) Ask(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*ListChildrenQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.units.Children(ctx, vo.ID(q.ParentID))
}
