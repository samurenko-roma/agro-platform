package weather

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrepo "github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type SyncForecastCommand struct {
	LocationID *string `json:"locationId"`
	Days       int     `json:"days" validate:"min=1,max=16"`
}

func (h *Handler) SyncForecast(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*SyncForecastCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok || orgID == "" {
		return nil, errors.New("organization_id required")
	}
	if cmd.Days == 0 {
		cmd.Days = 7
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

		records, err := h.service.SyncForecast(ctx, location, cmd.Days)
		if err != nil {
			return nil, err
		}
		return map[string]any{"synced": len(records)}, nil
	})
}
