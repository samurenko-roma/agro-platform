package production

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"

	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	growingcycle "github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/growing_cycle"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type growingCycleRepository struct {
	db uow.DB
}

func NewGrowingCycleRepository(db uow.DB) repository.GrowingCycleRepository {
	return &growingCycleRepository{db: db}
}

func (r *growingCycleRepository) Save(ctx context.Context, cycle *growingcycle.GrowingCycle) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO production_growing_cycles (
			id, farm_id, crop_id, variety_id, protocol_id,
			name, code, method, status, stage,
			created_at, updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,
			$11,$12
		)
		ON CONFLICT (id)
		DO UPDATE SET
			farm_id     = EXCLUDED.farm_id,
			crop_id     = EXCLUDED.crop_id,
			variety_id  = EXCLUDED.variety_id,
			protocol_id = EXCLUDED.protocol_id,
			name        = EXCLUDED.name,
			code        = EXCLUDED.code,
			method      = EXCLUDED.method,
			status      = EXCLUDED.status,
			stage       = EXCLUDED.stage,
			updated_at  = EXCLUDED.updated_at
	`,
		cycle.ID, cycle.FarmID, cycle.CropID, cycle.VarietyID, cycle.ProtocolID,
		cycle.Name, cycle.Code, cycle.Method, cycle.Status, cycle.Stage,
		cycle.CreatedAt, cycle.UpdatedAt,
	)

	return err
}

func (r *growingCycleRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*growingcycle.GrowingCycle, error) {
	root := &growingcycle.GrowingCycle{}

	err := r.db.QueryRow(ctx, `
		SELECT
			id, farm_id, crop_id, variety_id, protocol_id,
			name, code, method, status, stage,
			created_at, updated_at
		FROM production_growing_cycles
		WHERE id = $1 AND farm_id = $2
	`, id, farmID).Scan(scanGrowingCycle(root)...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (r *growingCycleRepository) GetByCode(ctx context.Context, code string, farmID vo.ID) (*growingcycle.GrowingCycle, error) {
	root := &growingcycle.GrowingCycle{}

	err := r.db.QueryRow(ctx, `
		SELECT
			id, farm_id, crop_id, variety_id, protocol_id,
			name, code, method, status, stage,
			created_at, updated_at
		FROM production_growing_cycles
		WHERE code = $1 AND farm_id = $2
	`, code, farmID).Scan(scanGrowingCycle(root)...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (r *growingCycleRepository) List(ctx context.Context, filter repository.ListFilter) ([]*growingcycle.GrowingCycle, error) {
	query := `
		SELECT
			id, farm_id, crop_id, variety_id, protocol_id,
			name, code, method, status, stage,
			created_at, updated_at
		FROM production_growing_cycles
		WHERE 1=1`

	args := make([]any, 0)
	argPos := 1

	if filter.FarmID != nil {
		query += ` AND farm_id = $` + strconv.Itoa(argPos)
		args = append(args, *filter.FarmID)
		argPos++
	}

	if filter.CropID != nil {
		query += ` AND crop_id = $` + strconv.Itoa(argPos)
		args = append(args, *filter.CropID)
		argPos++
	}

	if filter.Status != nil {
		query += ` AND status = $` + strconv.Itoa(argPos)
		args = append(args, *filter.Status)
		argPos++
	}

	query += ` ORDER BY created_at DESC`

	if filter.Limit > 0 {
		query += ` LIMIT $` + strconv.Itoa(argPos)
		args = append(args, filter.Limit)
		argPos++
	}

	if filter.Offset > 0 {
		query += ` OFFSET $` + strconv.Itoa(argPos)
		args = append(args, filter.Offset)
		argPos++
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*growingcycle.GrowingCycle, 0)

	for rows.Next() {
		root := &growingcycle.GrowingCycle{}
		if err := rows.Scan(scanGrowingCycle(root)...); err != nil {
			return nil, err
		}
		result = append(result, root)
	}

	return result, rows.Err()
}

func (r *growingCycleRepository) Delete(ctx context.Context, id vo.ID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM production_growing_cycles WHERE id = $1`, id)
	return err
}
