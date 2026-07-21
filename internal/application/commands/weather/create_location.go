package weather

import (
	"context"
	"errors"

	command "github.com/samurenkoroma/agro-platform/internal/application/commands"
	"github.com/samurenkoroma/agro-platform/internal/application/commands/response"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrepo "github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
	"github.com/samurenkoroma/agro-platform/internal/infrastructure/repository/providers"
	"github.com/samurenkoroma/agro-platform/internal/shared/repository"
)

type CreateLocationCommand struct {
	Name             string  `json:"name"      validate:"required"`
	Latitude         float64 `json:"latitude"  validate:"required"`
	Longitude        float64 `json:"longitude" validate:"required"`
	Timezone         string  `json:"timezone"`
	ProductionUnitID *string `json:"productionUnitId"`
	IsDefault        bool    `json:"isDefault"`
}

func (h *Handler) CreateLocation(ctx context.Context, payload any) (any, error) {
	cmd, ok := payload.(*CreateLocationCommand)
	if !ok {
		return nil, command.ErrInvalidCommandType
	}
	orgID, ok := ctx.Value("organization_id").(string)
	if !ok || orgID == "" {
		return nil, errors.New("organization_id required")
	}

	return h.uow.Execute(ctx, providers.NewWeatherProvider, func(p repository.RepositoryProvider, exec uow.Execution) (any, error) {
		wp, ok := p.(weatherrepo.WeatherProvider)
		if !ok {
			return nil, repository.ErrInvalidProviderType
		}

		farmID := vo.ID(orgID)
		loc := weatherlocation.New(farmID, cmd.Name, cmd.Latitude, cmd.Longitude)

		if cmd.Timezone != "" {
			loc.Timezone = cmd.Timezone
		}
		if cmd.ProductionUnitID != nil {
			id := vo.ID(*cmd.ProductionUnitID)
			loc.ProductionUnitID = &id
		}
		if cmd.IsDefault {
			loc.SetDefault()
		}

		if err := wp.Locations().Save(ctx, loc); err != nil {
			return nil, err
		}
		exec.RegisterAggregate(loc)
		return response.Id(loc.ID), nil
	})
}
