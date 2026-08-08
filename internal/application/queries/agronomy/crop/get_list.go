package crop

import (
	"context"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
)

type listCropHandler struct {
	crops Projection
}

type ListQuery struct {
	Search   *string  `json:"search"`
	Category []string `json:"categories"`
	Archived *bool    `json:"archived"`
}

func NewGetListHandler(crops Projection) queries.Handler {
	return &listCropHandler{crops: crops}
}

// Ask godoc
// @Summary Список культур
// @Description Список культур (с фильтром по категории)
// @Tags agronomy
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ListQuery false "Фильтр по категориям/поиску (опционально)"
// @Success 200 {object} response.SuccessResponse{data=[]ListItem}
// @Router /api/queries/agronomy.list_crops [post]
func (h *listCropHandler) Ask(ctx context.Context, query any) (any, error) {
	q, ok := query.(*ListQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.crops.List(ctx,
		ListFilter{
			Search:   q.Search,
			Category: q.Category,
			Archived: q.Archived,
		},
	)
}
