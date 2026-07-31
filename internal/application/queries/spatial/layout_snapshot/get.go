package layoutsnapshot

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type getHandler struct{ proj Projection }

func NewGet(proj Projection) queries.Handler { return &getHandler{proj: proj} }

type GetQuery struct {
	ID string `json:"id" validate:"required"`
}

func (h *getHandler) Ask(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*GetQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.proj.Get(ctx, vo.ID(q.ID))
}
