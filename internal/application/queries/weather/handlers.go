package weather

import (
	"context"
	"errors"
	"time"

	"github.com/samurenkoroma/agro-platform/internal/application/queries"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type QueryHandler struct{ proj Projection }

func NewQueryHandler(proj Projection) *QueryHandler { return &QueryHandler{proj: proj} }

// --- get_current ---

type GetCurrentQuery struct {
	LocationID *string `json:"locationId"`
}

func (h *QueryHandler) GetCurrent(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*GetCurrentQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok || orgID == "" {
		return nil, errors.New("organization_id required")
	}

	if q.LocationID != nil {
		return h.proj.GetCurrentByLocation(ctx, vo.ID(*q.LocationID))
	}
	return h.proj.GetCurrentByFarm(ctx, vo.ID(orgID))
}

// --- list_forecast ---

type ListForecastQuery struct {
	LocationID string     `json:"locationId" validate:"required"`
	From       *time.Time `json:"from"`
}

func (h *QueryHandler) ListForecast(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*ListForecastQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	from := time.Now()
	if q.From != nil {
		from = *q.From
	}
	return h.proj.ListForecast(ctx, vo.ID(q.LocationID), from)
}

// --- list_historical ---

type ListHistoricalQuery struct {
	LocationID string    `json:"locationId" validate:"required"`
	From       time.Time `json:"from"       validate:"required"`
	To         time.Time `json:"to"         validate:"required"`
}

func (h *QueryHandler) ListHistorical(ctx context.Context, payload any) (any, error) {
	q, ok := payload.(*ListHistoricalQuery)
	if !ok {
		return nil, queries.ErrInvalidQueryType
	}
	return h.proj.ListHistorical(ctx, vo.ID(q.LocationID), q.From, q.To)
}

// --- list_locations ---

type ListLocationsQuery struct{}

func (h *QueryHandler) ListLocations(ctx context.Context, payload any) (any, error) {
	if _, ok := payload.(*ListLocationsQuery); !ok {
		return nil, queries.ErrInvalidQueryType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok || orgID == "" {
		return nil, errors.New("organization_id required")
	}
	return h.proj.ListLocations(ctx, vo.ID(orgID))
}
