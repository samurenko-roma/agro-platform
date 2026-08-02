package growingcycle

import (
	"context"
	"errors"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type handler struct {
	cycles Projection
}

func NewGetOne(cycles Projection) queries.Handler {
	return &handler{
		cycles: cycles,
	}
}

type GetOneQuery struct {
	Id string `json:"id" validate:"required"`
}

func (h *handler) Ask(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*GetOneQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}
	return h.cycles.Summary(ctx, vo.ID(orgID), vo.ID(q.Id))
}
