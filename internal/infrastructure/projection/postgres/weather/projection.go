package weather

import (
	"context"
	"encoding/json"
	"time"

	weatherquery "github.com/samurenkoroma/agro-platform/internal/application/queries/weather"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherrecord "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_record"
)

type Projection struct{ db uow.DB }

func New(db uow.DB) weatherquery.Projection { return &Projection{db: db} }

func (p *Projection) GetCurrentByLocation(ctx context.Context, locationID vo.ID) (*weatherquery.WeatherRecordDTO, error) {
	row := p.db.QueryRow(ctx, `
		SELECT id, location_id, kind, source, timestamp, data
		FROM weather_records
		WHERE location_id = $1 AND kind = 'CURRENT'
		ORDER BY timestamp DESC LIMIT 1`, locationID)
	return scanRecordDTO(row)
}

func (p *Projection) GetCurrentByFarm(ctx context.Context, farmID vo.ID) (*weatherquery.WeatherRecordDTO, error) {
	row := p.db.QueryRow(ctx, `
		SELECT r.id, r.location_id, r.kind, r.source, r.timestamp, r.data
		FROM weather_records r
		JOIN weather_locations l ON l.id = r.location_id
		WHERE l.farm_id = $1 AND l.is_default = true AND r.kind = 'CURRENT'
		ORDER BY r.timestamp DESC LIMIT 1`, farmID)
	return scanRecordDTO(row)
}

func (p *Projection) ListForecast(ctx context.Context, locationID vo.ID, from time.Time) ([]*weatherquery.WeatherRecordDTO, error) {
	rows, err := p.db.Query(ctx, `
		SELECT id, location_id, kind, source, timestamp, data
		FROM weather_records
		WHERE location_id = $1 AND kind = 'FORECAST' AND timestamp >= $2
		ORDER BY timestamp ASC`, locationID, from)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRecordDTOs(rows)
}

func (p *Projection) ListHistorical(ctx context.Context, locationID vo.ID, from, to time.Time) ([]*weatherquery.WeatherRecordDTO, error) {
	rows, err := p.db.Query(ctx, `
		SELECT id, location_id, kind, source, timestamp, data
		FROM weather_records
		WHERE location_id = $1 AND kind = 'HISTORICAL'
		  AND timestamp >= $2 AND timestamp <= $3
		ORDER BY timestamp DESC`, locationID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRecordDTOs(rows)
}

func (p *Projection) ListLocations(ctx context.Context, farmID vo.ID) ([]*weatherquery.LocationDTO, error) {
	rows, err := p.db.Query(ctx, `
		SELECT id, name, latitude, longitude, timezone, is_default, production_unit_id
		FROM weather_locations
		WHERE farm_id = $1 AND archived_at IS NULL
		ORDER BY is_default DESC, name ASC`, farmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*weatherquery.LocationDTO, 0)
	for rows.Next() {
		var d weatherquery.LocationDTO
		if err := rows.Scan(&d.ID, &d.Name, &d.Latitude, &d.Longitude, &d.Timezone, &d.IsDefault, &d.ProductionUnitID); err != nil {
			return nil, err
		}
		result = append(result, &d)
	}
	return result, nil
}

func scanRecordDTO(s interface{ Scan(...any) error }) (*weatherquery.WeatherRecordDTO, error) {
	var dto weatherquery.WeatherRecordDTO
	var dataJSON []byte
	if err := s.Scan(&dto.ID, &dto.LocationID, &dto.Kind, &dto.Source, &dto.Timestamp, &dataJSON); err != nil {
		return nil, err
	}
	var data weatherrecord.WeatherData
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, err
	}
	dto.Data = weatherquery.MapData(data)
	return &dto, nil
}

func scanRecordDTOs(rows interface {
	Next() bool
	Scan(...any) error
	Close()
}) ([]*weatherquery.WeatherRecordDTO, error) {
	result := make([]*weatherquery.WeatherRecordDTO, 0)
	for rows.Next() {
		dto, err := scanRecordDTO(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, dto)
	}
	return result, nil
}
