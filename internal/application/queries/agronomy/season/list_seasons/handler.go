package listseasons

import (
	"context"
	"errors"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	"github.com/samurenkoroma/agro-platform/internal/application/queries/agronomy/season"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	projection "github.com/samurenkoroma/agro-platform/internal/infrastructure/projection/postgres/agronomy/season"
)

type handler struct {
	seasons season.Projection
}

func New(db uow.DB) queries.Handler {
	return &handler{seasons: projection.New(db)}
}

type Query struct {
}

// Ask godoc
// @Summary Список сезонов организации
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse
// @Router /api/queries/agronomy.list_seasons [post]
func (h *handler) Ask(ctx context.Context, query any) (any, error) {
	_, ok := query.(*Query)
	if !ok {
		return nil, errors.New("invalid query type")
	}
	orgId, ok := ctx.Value("organization_id").(string)
	if !ok {
		return nil, errors.New("organization_id is required")
	}

	return h.seasons.List(ctx, season.SeasonFilter{
		OwnerId: orgId,
	})
}
