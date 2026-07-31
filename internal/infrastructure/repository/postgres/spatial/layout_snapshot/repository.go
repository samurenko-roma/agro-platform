package layoutsnapshot

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
	ls "github.com/samurenkoroma/agro-platform/internal/domain/spatial/entity/layout_snapshot"
	pus "github.com/samurenkoroma/agro-platform/internal/domain/spatial/entity/production_unit_snapshot"
	"github.com/samurenkoroma/agro-platform/internal/domain/spatial/repository"
)

type snapshotRepository struct{ db uow.DB }

func NewRepository(db uow.DB) repository.LayoutSnapshotRepository {
	return &snapshotRepository{db: db}
}

func (r *snapshotRepository) Save(ctx context.Context, snapshot *ls.Aggregate) error {
	unitsJSON, err := json.Marshal(snapshot.Root.Units)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO spatial_layout_snapshots (
			id, farm_id, version, description, created_by, units, created_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET
			units = EXCLUDED.units
	`,
		snapshot.Root.ID, snapshot.Root.FarmID, snapshot.Root.Version,
		snapshot.Root.Description, snapshot.Root.CreatedBy, unitsJSON, snapshot.Root.CreatedAt,
	)
	return err
}

func (r *snapshotRepository) Get(ctx context.Context, id vo.ID) (*ls.Aggregate, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, farm_id, version, description, created_by, units, created_at
		FROM spatial_layout_snapshots WHERE id = $1`, id)
	snap, err := scanSnapshot(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return snap, err
}

func (r *snapshotRepository) GetLatest(ctx context.Context, farmID vo.ID) (*ls.Aggregate, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, farm_id, version, description, created_by, units, created_at
		FROM spatial_layout_snapshots
		WHERE farm_id = $1
		ORDER BY version DESC
		LIMIT 1`, farmID)
	snap, err := scanSnapshot(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return snap, err
}

func scanSnapshot(row interface{ Scan(...any) error }) (*ls.Aggregate, error) {
	var root ls.LayoutSnapshot
	var unitsJSON []byte

	if err := row.Scan(
		&root.ID, &root.FarmID, &root.Version, &root.Description,
		&root.CreatedBy, &unitsJSON, &root.CreatedAt,
	); err != nil {
		return nil, err
	}

	root.Units = make([]pus.ProductionUnitSnapshot, 0)
	if len(unitsJSON) > 0 {
		if err := json.Unmarshal(unitsJSON, &root.Units); err != nil {
			return nil, err
		}
	}

	return &ls.Aggregate{Root: root}, nil
}
