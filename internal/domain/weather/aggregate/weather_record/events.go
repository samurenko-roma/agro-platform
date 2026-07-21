package weatherrecord

import (
	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/event"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

const EventWeatherRecordCreated = "weather.record.created"

type WeatherRecordCreated struct {
	ev.BaseEvent
	Kind Kind
}

func NewWeatherRecordCreated(id vo.ID, kind Kind) WeatherRecordCreated {
	return WeatherRecordCreated{
		BaseEvent: ev.NewBaseEvent(id, EventWeatherRecordCreated),
		Kind:      kind,
	}
}
