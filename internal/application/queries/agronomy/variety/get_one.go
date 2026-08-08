package variety

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
)

type getVarietyHandler struct {
	varieties Projection
}

type OneQuery struct {
	Id string `json:"id" validate:"required"`
}

func NewGetOneHandler(varieties Projection) queries.Handler {
	return &getVarietyHandler{
		varieties: varieties,
	}
}

// Ask
// @Summary Получить сорт по ID
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body OneQuery true "ID сорта"
// @Success 200 {object} response.SuccessResponse{data=Detail}
// @Router /api/queries/agronomy.get_variety [post]
func (h *getVarietyHandler) Ask(ctx context.Context, query any) (any, error) {
	q, ok := query.(*OneQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.varieties.Get(ctx, q.Id)
}
