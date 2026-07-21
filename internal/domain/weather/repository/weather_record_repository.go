package repository

import (
	"context"
	"time"

	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
)

type RecordFilter struct {
	LocationID *vo.ID
	Kind       *weatherrecord.Kind
	From       *time.Time
	To         *time.Time
	Limit      int
}

type WeatherRecordRepository interface {
	Save(ctx context.Context, r *weatherrecord.WeatherRecord) error
	SaveBatch(ctx context.Context, records []*weatherrecord.WeatherRecord) error
	GetLatestCurrent(ctx context.Context, locationID vo.ID) (*weatherrecord.WeatherRecord, error)
	ListForecast(ctx context.Context, locationID vo.ID, from time.Time) ([]*weatherrecord.WeatherRecord, error)
	ListHistorical(ctx context.Context, filter RecordFilter) ([]*weatherrecord.WeatherRecord, error)
}
