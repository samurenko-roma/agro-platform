package layoutsnapshot

import (
	"context"
	"errors"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type getLatestHandler struct{ proj Projection }

func NewGetLatest(proj Projection) queries.Handler { return &getLatestHandler{proj: proj} }

type GetLatestQuery struct{}

func (h *getLatestHandler) Ask(ctx context.Context, payload any) (any, error) {
	if _, ok := payload.(*GetLatestQuery); !ok {
		return nil, queries.ErrInvalidQueryType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}
	return h.proj.GetLatest(ctx, vo.ID(orgID))
}
