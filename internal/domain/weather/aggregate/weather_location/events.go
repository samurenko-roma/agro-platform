package weatherlocation

import (
	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/event"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

const (
	EventWeatherLocationCreated  = "weather.location.created"
	EventWeatherLocationArchived = "weather.location.archived"
)

type WeatherLocationCreated struct{ ev.BaseEvent }

func NewWeatherLocationCreated(id vo.ID) WeatherLocationCreated {
	return WeatherLocationCreated{ev.NewBaseEvent(id, EventWeatherLocationCreated)}
}

type WeatherLocationArchived struct{ ev.BaseEvent }

func NewWeatherLocationArchived(id vo.ID) WeatherLocationArchived {
	return WeatherLocationArchived{ev.NewBaseEvent(id, EventWeatherLocationArchived)}
}
