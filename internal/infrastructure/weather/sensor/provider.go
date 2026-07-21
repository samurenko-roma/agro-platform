package sensorprovider

import (
	"context"
	"fmt"
	"time"

	sensordomain "github.com/samurenkoroma/agro-platform/internal/domain/environment/aggregate/sensor"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
	weatherprovider "github.com/samurenkoroma/agro-platform/internal/domain/weather/provider"
)

// SensorReader — минимальный интерфейс для получения датчиков фермы.
// Внедряется снаружи, чтобы не создавать циклической зависимости.
type SensorReader interface {
	ListByProductionUnit(ctx context.Context, unitID vo.ID) ([]*sensordomain.Sensor, error)
	ListByFarm(ctx context.Context, farmID vo.ID) ([]*sensordomain.Sensor, error)
}

// Provider агрегирует показания датчиков фермы в WeatherData.
// Не поддерживает прогноз и исторические данные — только текущее.
type Provider struct {
	sensors SensorReader
}

func NewProvider(sensors SensorReader) *Provider {
	return &Provider{sensors: sensors}
}

func (p *Provider) Name() weatherrecord.Source {
	return weatherrecord.SourceSensor
}

func (p *Provider) FetchCurrent(ctx context.Context, location *weatherlocation.WeatherLocation) (*weatherrecord.WeatherData, error) {
	var sensors []*sensordomain.Sensor
	var err error

	if location.ProductionUnitID != nil {
		sensors, err = p.sensors.ListByProductionUnit(ctx, *location.ProductionUnitID)
	} else {
		sensors, err = p.sensors.ListByFarm(ctx, location.FarmID)
	}
	if err != nil {
		return nil, fmt.Errorf("sensor provider: %w", err)
	}

	data := &weatherrecord.WeatherData{}
	for _, s := range sensors {
		if s.Status != sensordomain.Online || s.Value.Current == nil {
			continue
		}
		v := *s.Value.Current
		switch s.Type {
		case sensordomain.Temperature:
			data.Temperature = &v
		case sensordomain.Humidity:
			data.Humidity = &v
		case sensordomain.CO2:
			// CO2 не в WeatherData напрямую — можно расширить
		case sensordomain.SoilMoisture:
			data.SoilMoisture0_1cm = &v
		}
	}
	return data, nil
}

// FetchForecast датчики не поддерживают.
func (p *Provider) FetchForecast(_ context.Context, _ *weatherlocation.WeatherLocation, _ int) ([]weatherprovider.ForecastPoint, error) {
	return nil, fmt.Errorf("sensor provider does not support forecast")
}

// FetchHistorical датчики не поддерживают через этот интерфейс
// (исторические данные датчиков хранятся в Telemetry).
func (p *Provider) FetchHistorical(_ context.Context, _ *weatherlocation.WeatherLocation, _, _ time.Time) ([]weatherprovider.HistoricalPoint, error) {
	return nil, fmt.Errorf("sensor provider does not support historical data, use telemetry")
}
