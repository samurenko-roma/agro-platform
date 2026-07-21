package repository

import (
	"context"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
)

type WeatherLocationRepository interface {
	Save(ctx context.Context, l *weatherlocation.WeatherLocation) error
	GetByID(ctx context.Context, id vo.ID) (*weatherlocation.WeatherLocation, error)
	GetDefaultByFarm(ctx context.Context, farmID vo.ID) (*weatherlocation.WeatherLocation, error)
	ListByFarm(ctx context.Context, farmID vo.ID) ([]*weatherlocation.WeatherLocation, error)
}
