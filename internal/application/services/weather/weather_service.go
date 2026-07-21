package weatherservice

import (
	"context"
	"fmt"
	"time"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
	weatherprovider "github.com/samurenkoroma/agro-platform/internal/domain/weather/provider"
	"github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
	"github.com/samurenkoroma/agro-platform/pkg/logger"
)

// WeatherService — application service для получения погодных данных.
// Не зависит от конкретных провайдеров — работает через интерфейс WeatherProvider.
type WeatherService struct {
	primary   weatherprovider.WeatherProvider   // основной источник (open-meteo)
	fallbacks []weatherprovider.WeatherProvider // запасные (датчики, другие API)
	records   repository.WeatherRecordRepository
	locations repository.WeatherLocationRepository
}

func New(
	primary weatherprovider.WeatherProvider,
	fallbacks []weatherprovider.WeatherProvider,
	records repository.WeatherRecordRepository,
	locations repository.WeatherLocationRepository,
) *WeatherService {
	return &WeatherService{
		primary:   primary,
		fallbacks: fallbacks,
		records:   records,
		locations: locations,
	}
}

// SyncCurrent получает текущую погоду для локации и сохраняет в БД.
func (s *WeatherService) SyncCurrent(ctx context.Context, location *weatherlocation.WeatherLocation) (*weatherrecord.WeatherRecord, error) {
	log := logger.FromContext(ctx)

	data, source, err := s.fetchWithFallback(ctx, func(p weatherprovider.WeatherProvider) (*weatherrecord.WeatherData, error) {
		return p.FetchCurrent(ctx, location)
	})
	if err != nil {
		return nil, fmt.Errorf("SyncCurrent: all providers failed: %w", err)
	}

	record := weatherrecord.NewCurrent(location.ID, location.FarmID, source, *data)
	if err := s.records.Save(ctx, record); err != nil {
		return nil, fmt.Errorf("SyncCurrent: save record: %w", err)
	}

	log.Info("weather: synced current", "location", location.Name, "source", source)
	return record, nil
}

// SyncForecast получает прогноз на days дней и сохраняет в БД.
func (s *WeatherService) SyncForecast(ctx context.Context, location *weatherlocation.WeatherLocation, days int) ([]*weatherrecord.WeatherRecord, error) {
	log := logger.FromContext(ctx)

	points, err := s.primary.FetchForecast(ctx, location, days)
	if err != nil {
		return nil, fmt.Errorf("SyncForecast: %w", err)
	}

	records := make([]*weatherrecord.WeatherRecord, 0, len(points))
	for _, pt := range points {
		r := weatherrecord.NewForecast(location.ID, location.FarmID, s.primary.Name(), pt.Time, pt.Data)
		records = append(records, r)
	}

	if err := s.records.SaveBatch(ctx, records); err != nil {
		return nil, fmt.Errorf("SyncForecast: save batch: %w", err)
	}

	log.Info("weather: synced forecast", "location", location.Name, "points", len(records))
	return records, nil
}

// SyncHistorical получает исторические данные за период и сохраняет в БД.
func (s *WeatherService) SyncHistorical(ctx context.Context, location *weatherlocation.WeatherLocation, from, to time.Time) ([]*weatherrecord.WeatherRecord, error) {
	log := logger.FromContext(ctx)

	points, err := s.primary.FetchHistorical(ctx, location, from, to)
	if err != nil {
		return nil, fmt.Errorf("SyncHistorical: %w", err)
	}

	records := make([]*weatherrecord.WeatherRecord, 0, len(points))
	for _, pt := range points {
		r := weatherrecord.NewHistorical(location.ID, location.FarmID, s.primary.Name(), pt.Time, pt.Data)
		records = append(records, r)
	}

	if err := s.records.SaveBatch(ctx, records); err != nil {
		return nil, fmt.Errorf("SyncHistorical: save batch: %w", err)
	}

	log.Info("weather: synced historical", "location", location.Name, "from", from, "to", to, "records", len(records))
	return records, nil
}

// GetCurrentForFarm возвращает последнюю текущую запись для дефолтной локации фермы.
func (s *WeatherService) GetCurrentForFarm(ctx context.Context, farmID vo.ID) (*weatherrecord.WeatherRecord, error) {
	location, err := s.locations.GetDefaultByFarm(ctx, farmID)
	if err != nil || location == nil {
		return nil, fmt.Errorf("GetCurrentForFarm: no default location for farm %s", farmID)
	}
	return s.records.GetLatestCurrent(ctx, location.ID)
}

// fetchWithFallback пробует primary провайдер, при ошибке — fallbacks по очереди.
func (s *WeatherService) fetchWithFallback(
	ctx context.Context,
	fetch func(weatherprovider.WeatherProvider) (*weatherrecord.WeatherData, error),
) (*weatherrecord.WeatherData, weatherrecord.Source, error) {
	log := logger.FromContext(ctx)

	data, err := fetch(s.primary)
	if err == nil {
		return data, s.primary.Name(), nil
	}
	log.Warn("weather: primary provider failed, trying fallbacks", "error", err)

	for _, fb := range s.fallbacks {
		data, err = fetch(fb)
		if err == nil {
			return data, fb.Name(), nil
		}
		log.Warn("weather: fallback provider failed", "provider", fb.Name(), "error", err)
	}
	return nil, "", fmt.Errorf("all providers failed, last error: %w", err)
}
