package layoutsnapshot

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	lsquery "github.com/samurenkoroma/agro-platform/internal/application/queries/spatial/layout_snapshot"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type projection struct{ db uow.DB }

func New(db uow.DB) lsquery.Projection { return &projection{db: db} }

func (p *projection) Get(ctx context.Context, id vo.ID) (*lsquery.SnapshotDTO, error) {
	row := p.db.QueryRow(ctx, `
		SELECT id, farm_id, version, description, created_by, units, created_at
		FROM spatial_layout_snapshots WHERE id = $1`, id)
	dto, err := scanSnapshot(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return dto, err
}

func (p *projection) GetLatest(ctx context.Context, farmID vo.ID) (*lsquery.SnapshotDTO, error) {
	row := p.db.QueryRow(ctx, `
		SELECT id, farm_id, version, description, created_by, units, created_at
		FROM spatial_layout_snapshots
		WHERE farm_id = $1
		ORDER BY version DESC
		LIMIT 1`, farmID)
	dto, err := scanSnapshot(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return dto, err
}

func scanSnapshot(row interface{ Scan(...any) error }) (*lsquery.SnapshotDTO, error) {
	var dto lsquery.SnapshotDTO
	var unitsJSON []byte
	if err := row.Scan(
		&dto.ID, &dto.FarmID, &dto.Version, &dto.Description,
		&dto.CreatedBy, &unitsJSON, &dto.CreatedAt,
	); err != nil {
		return nil, err
	}
	dto.Units = make([]lsquery.UnitDTO, 0)
	if len(unitsJSON) > 0 {
		if err := json.Unmarshal(unitsJSON, &dto.Units); err != nil {
			return nil, err
		}
	}
	return &dto, nil
}
