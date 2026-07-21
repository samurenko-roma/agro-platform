package weather

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	weatherlocation "github.com/samurenkoroma/agro-platform/internal/domain/weather/aggregate/weather_location"
	"github.com/samurenkoroma/agro-platform/internal/domain/weather/repository"
)

type locationRepository struct{ db uow.DB }

func NewLocationRepository(db uow.DB) repository.WeatherLocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Save(ctx context.Context, l *weatherlocation.WeatherLocation) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO weather_locations (
			id, farm_id, production_unit_id, name,
			latitude, longitude, timezone, is_default,
			created_at, updated_at, archived_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (id) DO UPDATE SET
			name              = EXCLUDED.name,
			production_unit_id = EXCLUDED.production_unit_id,
			latitude          = EXCLUDED.latitude,
			longitude         = EXCLUDED.longitude,
			timezone          = EXCLUDED.timezone,
			is_default        = EXCLUDED.is_default,
			updated_at        = EXCLUDED.updated_at,
			archived_at       = EXCLUDED.archived_at
	`,
		l.ID, l.FarmID, l.ProductionUnitID, l.Name,
		l.Latitude, l.Longitude, l.Timezone, l.IsDefault,
		l.CreatedAt, l.UpdatedAt, l.ArchivedAt,
	)
	return err
}

func (r *locationRepository) GetByID(ctx context.Context, id vo.ID) (*weatherlocation.WeatherLocation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, farm_id, production_unit_id, name,
		       latitude, longitude, timezone, is_default,
		       created_at, updated_at, archived_at
		FROM weather_locations WHERE id = $1`, id)
	return scanLocation(row)
}

func (r *locationRepository) GetDefaultByFarm(ctx context.Context, farmID vo.ID) (*weatherlocation.WeatherLocation, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, farm_id, production_unit_id, name,
		       latitude, longitude, timezone, is_default,
		       created_at, updated_at, archived_at
		FROM weather_locations
		WHERE farm_id = $1 AND is_default = true AND archived_at IS NULL
		LIMIT 1`, farmID)
	l, err := scanLocation(row)
	if err != nil {
		return nil, nil // не найдено — не ошибка
	}
	return l, nil
}

func (r *locationRepository) ListByFarm(ctx context.Context, farmID vo.ID) ([]*weatherlocation.WeatherLocation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, farm_id, production_unit_id, name,
		       latitude, longitude, timezone, is_default,
		       created_at, updated_at, archived_at
		FROM weather_locations
		WHERE farm_id = $1 AND archived_at IS NULL
		ORDER BY is_default DESC, name ASC`, farmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*weatherlocation.WeatherLocation, 0)
	for rows.Next() {
		l, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, l)
	}
	return result, nil
}

func scanLocation(s interface{ Scan(...any) error }) (*weatherlocation.WeatherLocation, error) {
	var l weatherlocation.WeatherLocation
	err := s.Scan(
		&l.ID, &l.FarmID, &l.ProductionUnitID, &l.Name,
		&l.Latitude, &l.Longitude, &l.Timezone, &l.IsDefault,
		&l.CreatedAt, &l.UpdatedAt, &l.ArchivedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, weatherlocation.ErrLocationNotFound
		}
		return nil, err
	}
	return &l, nil
}
