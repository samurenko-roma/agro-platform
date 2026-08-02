package production

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/planting"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type plantingRepository struct {
	db uow.DB
}

func NewPlantingRepository(db uow.DB) repository.PlantingRepository {
	return &plantingRepository{db: db}
}

func (r *plantingRepository) Save(ctx context.Context, root *planting.Planting) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO production_plantings (
			id, farm_id, cycle_id, planted_at, quantity, created_at, updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7
		)
		ON CONFLICT (id)
		DO UPDATE SET
			planted_at = EXCLUDED.planted_at,
			quantity   = EXCLUDED.quantity,
			updated_at = EXCLUDED.updated_at
	`,
		root.ID, root.FarmID, root.CycleID, root.PlantedAt, root.Quantity,
		root.CreatedAt, root.UpdatedAt,
	)

	return err
}

func (r *plantingRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*planting.Planting, error) {
	root := &planting.Planting{}

	err := r.db.QueryRow(ctx, `
		SELECT id, farm_id, cycle_id, planted_at, quantity, created_at, updated_at
		FROM production_plantings
		WHERE id = $1 AND farm_id = $2
	`, id, farmID).Scan(scanPlanting(root)...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (r *plantingRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*planting.Planting, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, farm_id, cycle_id, planted_at, quantity, created_at, updated_at
		FROM production_plantings
		WHERE cycle_id = $1
		ORDER BY planted_at
	`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*planting.Planting, 0)
	for rows.Next() {
		root := &planting.Planting{}
		if err := rows.Scan(scanPlanting(root)...); err != nil {
			return nil, err
		}
		result = append(result, root)
	}
	return result, rows.Err()
}

func (r *plantingRepository) Delete(ctx context.Context, id vo.ID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM production_plantings WHERE id = $1`, id)
	return err
}
