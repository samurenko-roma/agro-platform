package variety

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
)

type listVarietyHandler struct {
	varieties Projection
}

func New(varieties Projection) queries.Handler {
	return &listVarietyHandler{
		varieties: varieties,
	}
}

type ListQuery struct {
	CropID string `json:"cropId,omitempty"` // tomato, eggplant, cucumber
}

// Ask godoc
// @Summary Список сортов культуры
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ListQuery true "ID культуры"
// @Success 200 {object} response.SuccessResponse{data=[]ListItem}
// @Router /api/queries/agronomy.list_varieties [post]
func (h *listVarietyHandler) Ask(ctx context.Context, query any) (any, error) {
	q, ok := query.(*ListQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}

	return h.varieties.List(ctx, ListFilter{
		CropID: q.CropID,
	})

}
