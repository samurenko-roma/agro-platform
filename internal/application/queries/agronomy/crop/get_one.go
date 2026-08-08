package crop

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
)

type getCropHandler struct{ crops Projection }
type OneQuery struct {
	ID string `json:"id" validate:"required"`
}

func NewGetOneHandler(crops Projection) queries.Handler {
	return &getCropHandler{crops: crops}
}

// Ask godoc
// @Summary Получить культуру по ID
// @Description Реальный вызов: POST /api/queries, {"query": "agronomy.get_crop", "data": <тело ниже>}
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body OneQuery true "ID культуры"
// @Success 200 {object} response.SuccessResponse{data=Detail}
// @Router /api/queries/agronomy.get_crop [post]
func (h *getCropHandler) Ask(ctx context.Context, query any) (any, error) {
	q, ok := query.(*OneQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.crops.Get(ctx, q.ID)
}
