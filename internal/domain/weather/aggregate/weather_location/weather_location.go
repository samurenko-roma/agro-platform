package weatherlocation

import (
	"time"

	ev "github.com/samurenkoroma/agro-platform/internal/domain/shared/aggregate"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

// WeatherLocation описывает точку, для которой собираются погодные данные.
// Привязана к ферме, опционально к ProductionUnit (поле, теплица).
// Хранит координаты и имя для отображения.
type WeatherLocation struct {
	ev.BaseAggregate
	ID               vo.ID
	FarmID           vo.ID
	ProductionUnitID *vo.ID
	Name             string
	Latitude         float64
	Longitude        float64
	Timezone         string
	IsDefault        bool // основная точка фермы
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ArchivedAt       *time.Time
}

func New(farmID vo.ID, name string, lat, lon float64) *WeatherLocation {
	now := time.Now()
	w := &WeatherLocation{
		ID:        vo.NewID(),
		FarmID:    farmID,
		Name:      name,
		Latitude:  lat,
		Longitude: lon,
		Timezone:  "UTC",
		CreatedAt: now,
		UpdatedAt: now,
	}
	w.AddEvent(NewWeatherLocationCreated(w.ID))
	return w
}

func (w *WeatherLocation) SetDefault() {
	w.IsDefault = true
	w.UpdatedAt = time.Now()
}

func (w *WeatherLocation) Archive() {
	now := time.Now()
	w.ArchivedAt = &now
	w.UpdatedAt = now
	w.AddEvent(NewWeatherLocationArchived(w.ID))
}
