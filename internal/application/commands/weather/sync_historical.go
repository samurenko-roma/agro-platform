package weather

import (
	"context"
	"errors"
	"time"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrepo "github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type SyncHistoricalCommand struct {
	LocationID *string   `json:"locationId"`
	From       time.Time `json:"from" validate:"required"`
	To         time.Time `json:"to"   validate:"required"`
}

func (h *Handler) SyncHistorical(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*SyncHistoricalCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok || orgID == "" {
		return nil, errors.New("organization_id required")
	}
	if cmd.To.Before(cmd.From) {
		return nil, errors.New("to must be after from")
	}

	return h.uow.Execute(ctx, providers.NewWeatherProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		wp, ok := p.(weatherrepo.WeatherProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		farmID := vo.ID(orgID)

		var location *weatherlocation.WeatherLocation
		var err error
		if cmd.LocationID != nil {
			location, err = wp.Locations().GetByID(ctx, vo.ID(*cmd.LocationID))
		} else {
			location, err = wp.Locations().GetDefaultByFarm(ctx, farmID)
		}
		if err != nil || location == nil {
			return nil, errors.New("location not found")
		}

		records, err := h.service.SyncHistorical(ctx, location, cmd.From, cmd.To)
		if err != nil {
			return nil, err
		}
		return map[string]any{"synced": len(records)}, nil
	})
}
