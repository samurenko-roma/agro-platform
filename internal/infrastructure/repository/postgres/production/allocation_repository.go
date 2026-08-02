package production

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/samurenkoroma/agro-platform/internal/application/uow"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/aggregate/allocation"
	"github.com/samurenkoroma/agro-platform/internal/domain/production/repository"
	vo "github.com/samurenkoroma/agro-platform/internal/domain/shared/valueobject"
)

type allocationRepository struct {
	db uow.DB
}

func NewAllocationRepository(db uow.DB) repository.AllocationRepository {
	return &allocationRepository{db: db}
}

func (r *allocationRepository) ListActiveByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, farm_id, cycle_id, production_unit_id, area, started_at, ended_at, created_at, updated_at
		FROM production_allocations
		WHERE production_unit_id = $1 AND ended_at IS NULL
	`, productionUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*allocation.Allocation, 0)
	for rows.Next() {
		root := &allocation.Allocation{}
		if err := rows.Scan(scanAllocation(root)...); err != nil {
			return nil, err
		}
		result = append(result, root)
	}
	return result, rows.Err()
}

func (r *allocationRepository) Save(ctx context.Context, root *allocation.Allocation) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO production_allocations (
			id, farm_id, cycle_id, production_unit_id, area,
			started_at, ended_at, created_at, updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9
		)
		ON CONFLICT (id)
		DO UPDATE SET
			cycle_id = EXCLUDED.cycle_id,
			production_unit_id = EXCLUDED.production_unit_id,
			area = EXCLUDED.area,
			started_at = EXCLUDED.started_at,
			ended_at = EXCLUDED.ended_at,
			updated_at = EXCLUDED.updated_at
	`,
		root.ID, root.FarmID, root.CycleID, root.ProductionUnitID, root.Area,
		root.StartedAt, root.EndedAt, root.CreatedAt, root.UpdatedAt,
	)

	return err
}

func (r *allocationRepository) GetByID(ctx context.Context, id vo.ID, farmID vo.ID) (*allocation.Allocation, error) {
	root := &allocation.Allocation{}

	err := r.db.QueryRow(ctx, `
		SELECT id, farm_id, cycle_id, production_unit_id, area, started_at, ended_at, created_at, updated_at
		FROM production_allocations
		WHERE id = $1 AND farm_id = $2
	`, id, farmID).Scan(scanAllocation(root)...)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return root, nil
}

func (r *allocationRepository) ListByCycleID(ctx context.Context, cycleID vo.ID) ([]*allocation.Allocation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, farm_id, cycle_id, production_unit_id, area, started_at, ended_at, created_at, updated_at
		FROM production_allocations
		WHERE cycle_id = $1
		ORDER BY started_at
	`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*allocation.Allocation, 0)
	for rows.Next() {
		root := &allocation.Allocation{}
		if err := rows.Scan(scanAllocation(root)...); err != nil {
			return nil, err
		}
		result = append(result, root)
	}
	return result, rows.Err()
}

func (r *allocationRepository) ListByProductionUnitID(ctx context.Context, productionUnitID vo.ID) ([]*allocation.Allocation, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, farm_id, cycle_id, production_unit_id, area, started_at, ended_at, created_at, updated_at
		FROM production_allocations
		WHERE production_unit_id = $1
		ORDER BY started_at
	`, productionUnitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*allocation.Allocation, 0)
	for rows.Next() {
		root := &allocation.Allocation{}
		if err := rows.Scan(scanAllocation(root)...); err != nil {
			return nil, err
		}
		result = append(result, root)
	}
	return result, rows.Err()
}

func (r *allocationRepository) Delete(ctx context.Context, id vo.ID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM production_allocations WHERE id = $1`, id)
	return err
}
