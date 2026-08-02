package production

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	harvestbatch "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/harvest_batch"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type harvestRepository struct {
	db uow.DB
}

func NewHarvestRepository(db uow.DB) repository.HarvestBatchRepository {
	return &harvestRepository{db: db}
}

func (r *harvestRepository) Save(ctx context.Context, root *harvestbatch.HarvestBatch) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO production_harvest_batch (
			id, farm_id, cycle_id, harvested_at, quantity, created_at, updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7
		)
		ON CONFLICT (id)
		DO UPDATE SET
			harvested_at = EXCLUDED.harvested_at,
			quantity     = EXCLUDED.quantity,
			updated_at   = EXCLUDED.updated_at
	`,
		root.ID, root.FarmID, root.CycleID, root.HarvestedAt, root.Quantity,
		root.CreatedAt, root.UpdatedAt,
	)

	return err
}

func (r *harvestRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*harvestbatch.HarvestBatch, error) {
	root := &harvestbatch.HarvestBatch{}

	err := r.db.QueryRow(ctx, `
		SELECT id, farm_id, cycle_id, harvested_at, quantity, created_at, updated_at
		FROM production_harvest_batch
		WHERE id = $1 AND farm_id = $2
	`, id, farmID).Scan(scanHarvest(root)...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (r *harvestRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*harvestbatch.HarvestBatch, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, farm_id, cycle_id, harvested_at, quantity, created_at, updated_at
		FROM production_harvest_batch
		WHERE cycle_id = $1
		ORDER BY harvested_at
	`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*harvestbatch.HarvestBatch, 0)
	for rows.Next() {
		root := &harvestbatch.HarvestBatch{}
		if err := rows.Scan(scanHarvest(root)...); err != nil {
			return nil, err
		}
		result = append(result, root)
	}
	return result, rows.Err()
}

func (r *harvestRepository) Delete(ctx context.Context, id vo.ID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM production_harvest_batch WHERE id = $1`, id)
	return err
}
