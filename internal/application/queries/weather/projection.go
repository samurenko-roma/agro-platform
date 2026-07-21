package weather

import (
	"context"
	"time"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
)

type WeatherDataDTO struct {
	Temperature          *float64 `json:"temperature,omitempty"`
	TemperatureFeelsLike *float64 `json:"temperatureFeelsLike,omitempty"`
	Humidity             *float64 `json:"humidity,omitempty"`
	DewPoint             *float64 `json:"dewPoint,omitempty"`
	Precipitation        *float64 `json:"precipitation,omitempty"`
	PrecipitationProb    *float64 `json:"precipitationProbability,omitempty"`
	Rain                 *float64 `json:"rain,omitempty"`
	Snowfall             *float64 `json:"snowfall,omitempty"`
	WindSpeed            *float64 `json:"windSpeed,omitempty"`
	WindDirection        *float64 `json:"windDirection,omitempty"`
	WindGusts            *float64 `json:"windGusts,omitempty"`
	PressureSea          *float64 `json:"pressure,omitempty"`
	CloudCover           *float64 `json:"cloudCover,omitempty"`
	SolarRadiation       *float64 `json:"solarRadiation,omitempty"`
	UVIndex              *float64 `json:"uvIndex,omitempty"`
	SunshineDuration     *float64 `json:"sunshineDuration,omitempty"`
	SoilTemperature0cm   *float64 `json:"soilTemperature0cm,omitempty"`
	SoilTemperature6cm   *float64 `json:"soilTemperature6cm,omitempty"`
	SoilMoisture0_1cm    *float64 `json:"soilMoisture0to1cm,omitempty"`
	SoilMoisture1_3cm    *float64 `json:"soilMoisture1to3cm,omitempty"`
	Evapotranspiration   *float64 `json:"evapotranspiration,omitempty"`
	VapourPressureDef    *float64 `json:"vpd,omitempty"`
	WeatherCode          *int     `json:"weatherCode,omitempty"`
	Visibility           *float64 `json:"visibility,omitempty"`
	IsDay                *bool    `json:"isDay,omitempty"`
}

type WeatherRecordDTO struct {
	ID         vo.ID          `json:"id"`
	LocationID vo.ID          `json:"locationId"`
	Kind       string         `json:"kind"`
	Source     string         `json:"source"`
	Timestamp  time.Time      `json:"timestamp"`
	Data       WeatherDataDTO `json:"data"`
}

type LocationDTO struct {
	ID               vo.ID   `json:"id"`
	Name             string  `json:"name"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Timezone         string  `json:"timezone"`
	IsDefault        bool    `json:"isDefault"`
	ProductionUnitID *vo.ID  `json:"productionUnitId,omitempty"`
}

type Projection interface {
	GetCurrentByLocation(ctx context.Context, locationID vo.ID) (*WeatherRecordDTO, error)
	GetCurrentByFarm(ctx context.Context, farmID vo.ID) (*WeatherRecordDTO, error)
	ListForecast(ctx context.Context, locationID vo.ID, from time.Time) ([]*WeatherRecordDTO, error)
	ListHistorical(ctx context.Context, locationID vo.ID, from, to time.Time) ([]*WeatherRecordDTO, error)
	ListLocations(ctx context.Context, farmID vo.ID) ([]*LocationDTO, error)
}

func MapData(d weatherrecord.WeatherData) WeatherDataDTO {
	return WeatherDataDTO{
		Temperature:          d.Temperature,
		TemperatureFeelsLike: d.TemperatureFeelsLike,
		Humidity:             d.Humidity,
		DewPoint:             d.DewPoint,
		Precipitation:        d.Precipitation,
		PrecipitationProb:    d.PrecipitationProb,
		Rain:                 d.Rain,
		Snowfall:             d.Snowfall,
		WindSpeed:            d.WindSpeed,
		WindDirection:        d.WindDirection,
		WindGusts:            d.WindGusts,
		PressureSea:          d.PressureSea,
		CloudCover:           d.CloudCover,
		SolarRadiation:       d.SolarRadiation,
		UVIndex:              d.UVIndex,
		SunshineDuration:     d.SunshineDuration,
		SoilTemperature0cm:   d.SoilTemperature0cm,
		SoilTemperature6cm:   d.SoilTemperature6cm,
		SoilMoisture0_1cm:    d.SoilMoisture0_1cm,
		SoilMoisture1_3cm:    d.SoilMoisture1_3cm,
		Evapotranspiration:   d.Evapotranspiration,
		VapourPressureDef:    d.VapourPressureDef,
		WeatherCode:          d.WeatherCode,
		Visibility:           d.Visibility,
		IsDay:                d.IsDay,
	}
}
