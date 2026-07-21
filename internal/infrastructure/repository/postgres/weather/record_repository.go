package weather

import (
	"context"
	"encoding/json"
	"time"

	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
	"github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
)

type recordRepository struct{ db uow.DB }

func NewRecordRepository(db uow.DB) repository.WeatherRecordRepository {
	return &recordRepository{db: db}
}

func (r *recordRepository) Save(ctx context.Context, rec *weatherrecord.WeatherRecord) error {
	data, err := json.Marshal(rec.Data)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO weather_records (
			id, location_id, farm_id, kind, source,
			timestamp, forecast_at, data, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO NOTHING`,
		rec.ID, rec.LocationID, rec.FarmID, rec.Kind, rec.Source,
		rec.Timestamp, rec.ForecastAt, data, rec.CreatedAt,
	)
	return err
}

func (r *recordRepository) SaveBatch(ctx context.Context, records []*weatherrecord.WeatherRecord) error {
	for _, rec := range records {
		if err := r.Save(ctx, rec); err != nil {
			return err
		}
	}
	return nil
}

func (r *recordRepository) GetLatestCurrent(ctx context.Context, locationID vo.ID) (*weatherrecord.WeatherRecord, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, location_id, farm_id, kind, source,
		       timestamp, forecast_at, data, created_at
		FROM weather_records
		WHERE location_id = $1 AND kind = 'CURRENT'
		ORDER BY timestamp DESC
		LIMIT 1`, locationID)
	return scanRecord(row)
}

func (r *recordRepository) ListForecast(ctx context.Context, locationID vo.ID, from time.Time) ([]*weatherrecord.WeatherRecord, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, location_id, farm_id, kind, source,
		       timestamp, forecast_at, data, created_at
		FROM weather_records
		WHERE location_id = $1
		  AND kind = 'FORECAST'
		  AND timestamp >= $2
		ORDER BY timestamp ASC`, locationID, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRecords(rows)
}

func (r *recordRepository) ListHistorical(ctx context.Context, filter repository.RecordFilter) ([]*weatherrecord.WeatherRecord, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, location_id, farm_id, kind, source,
		       timestamp, forecast_at, data, created_at
		FROM weather_records
		WHERE location_id = $1
		  AND kind = 'HISTORICAL'
		  AND ($2::timestamptz IS NULL OR timestamp >= $2)
		  AND ($3::timestamptz IS NULL OR timestamp <= $3)
		ORDER BY timestamp DESC`,
		filter.LocationID, filter.From, filter.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRecords(rows)
}

func scanRecord(s interface{ Scan(...any) error }) (*weatherrecord.WeatherRecord, error) {
	var rec weatherrecord.WeatherRecord
	var dataJSON []byte
	if err := s.Scan(
		&rec.ID, &rec.LocationID, &rec.FarmID, &rec.Kind, &rec.Source,
		&rec.Timestamp, &rec.ForecastAt, &dataJSON, &rec.CreatedAt,
	); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(dataJSON, &rec.Data); err != nil {
		return nil, err
	}
	return &rec, nil
}

func scanRecords(rows interface {
	Next() bool
	Scan(...any) error
	Close()
}) ([]*weatherrecord.WeatherRecord, error) {
	result := make([]*weatherrecord.WeatherRecord, 0)
	for rows.Next() {
		rec, err := scanRecord(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, rec)
	}
	return result, nil
}
