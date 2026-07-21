package weatherprovider

import (
	"context"
	"time"

	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
)

// WeatherProvider — порт для любого источника погодных данных.
// Open-Meteo, датчик, Windy, Weather API — все реализуют этот интерфейс.
type WeatherProvider interface {
	// Name возвращает идентификатор провайдера (для логов и хранения Source).
	Name() weatherrecord.Source

	// FetchCurrent возвращает текущие погодные данные для точки.
	FetchCurrent(ctx context.Context, location *weatherlocation.WeatherLocation) (*weatherrecord.WeatherData, error)

	// FetchForecast возвращает почасовой прогноз на days дней вперёд.
	FetchForecast(ctx context.Context, location *weatherlocation.WeatherLocation, days int) ([]ForecastPoint, error)

	// FetchHistorical возвращает исторические данные за период.
	FetchHistorical(ctx context.Context, location *weatherlocation.WeatherLocation, from, to time.Time) ([]HistoricalPoint, error)
}

type ForecastPoint struct {
	Time time.Time
	Data weatherrecord.WeatherData
}

type HistoricalPoint struct {
	Time time.Time
	Data weatherrecord.WeatherData
}
