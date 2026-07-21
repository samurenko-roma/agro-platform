package weatherrecord

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type WeatherRecord struct {
	ev.BaseAggregate
	ID         vo.ID
	LocationID vo.ID
	FarmID     vo.ID
	Kind       Kind       // CURRENT | FORECAST | HISTORICAL
	Source     Source     // OPEN_METEO | SENSOR | CUSTOM
	Timestamp  time.Time  // время наблюдения / прогноза
	ForecastAt *time.Time // когда был сделан прогноз (для FORECAST)
	Data       WeatherData
	CreatedAt  time.Time
}

func NewCurrent(locationID, farmID vo.ID, source Source, data WeatherData) *WeatherRecord {
	return newRecord(locationID, farmID, Current, source, time.Now(), nil, data)
}

func NewForecast(locationID, farmID vo.ID, source Source, forecastFor time.Time, data WeatherData) *WeatherRecord {
	now := time.Now()
	return newRecord(locationID, farmID, Forecast, source, forecastFor, &now, data)
}

func NewHistorical(locationID, farmID vo.ID, source Source, observedAt time.Time, data WeatherData) *WeatherRecord {
	return newRecord(locationID, farmID, Historical, source, observedAt, nil, data)
}

func newRecord(locationID, farmID vo.ID, kind Kind, source Source, ts time.Time, forecastAt *time.Time, data WeatherData) *WeatherRecord {
	r := &WeatherRecord{
		ID:         vo.NewID(),
		LocationID: locationID,
		FarmID:     farmID,
		Kind:       kind,
		Source:     source,
		Timestamp:  ts,
		ForecastAt: forecastAt,
		Data:       data,
		CreatedAt:  time.Now(),
	}
	r.AddEvent(NewWeatherRecordCreated(r.ID, kind))
	return r
}
